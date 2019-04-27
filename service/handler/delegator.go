package handler

import (
	"cosmos-sync/logger"
	"cosmos-sync/store"
	"cosmos-sync/store/document"
	"cosmos-sync/types"
	"cosmos-sync/util/constant"
	"cosmos-sync/util/helper"
	"sync"
	"fmt"
)

// init delegator for genesis validator
func InitDelegator() {
	validators := helper.GetValidators()
	for _, validator := range validators {
		valAddr := validator.OperatorAddress.String()
		valAccAddr := helper.ValAddrToAccAddr(valAddr)
		modifyDelegator(valAccAddr, valAddr)
	}
}

// save or update validator or delegator info
// by parse stake tx

// Different transaction types correspond to different operations
//TxTypeStakeCreateValidator
//	1:insert new validator (---> CompareAndUpdateValidators)
//	2:insert delegator
//
//TxTypeStakeEditValidator
//	1:update validator
//
//TxTypeStakeDelegate
//	1:update validator (---> CompareAndUpdateValidators)
//	2:insert delegator(or update delegator existed )
//
//TxTypeStakeBeginUnbonding
//	1:update validator (---> CompareAndUpdateValidators)
//	2:update delegator
//
//TxTypeBeginRedelegate
//	1:update validator(src,dest) (---> CompareAndUpdateValidators)
//	2:update delegator(src,dest)
func SaveOrUpdateDelegator(docTx document.CommonTx, mutex sync.Mutex) {

	logger.Debug("Start", logger.String("method", "saveDelegator"))

	switch docTx.Type {
	case constant.TxTypeStakeCreateValidator:
		modifyDelegator(docTx.From, docTx.To)
		break
	case constant.TxTypeStakeEditValidator:
		updateValidator(docTx.From)
		break
	case constant.TxTypeStakeDelegate, constant.TxTypeStakeBeginUnbonding:
		modifyDelegator(docTx.From, docTx.To)
		break
	case constant.TxTypeBeginRedelegate:
		delAddress := docTx.From
		msg := docTx.Msg.(types.BeginRedelegate)
		valSrcAddr := msg.ValidatorSrcAddr
		valDstAddr := msg.ValidatorDstAddr

		modifyDelegator(delAddress, valSrcAddr)
		modifyDelegator(delAddress, valDstAddr)
		break
	}

	logger.Debug("End", logger.String("method", "saveDelegator"))
}

func modifyDelegator(delAddress, valAddress string) {
	logger.Info("delegator info has changed", logger.String("delAddress", delAddress), logger.String("valAddress", valAddress))
	// get delegation
	delegation := BuildDelegation(delAddress, valAddress)

	// get unbondingDelegation
	ud := BuildUnbondingDelegation(delAddress, valAddress)

	delegator := document.Delegator{
		Address:       delAddress,
		ValidatorAddr: valAddress,

		Shares:         delegation.Shares,
		OriginalShares: delegation.OriginalShares,
		BondedHeight:   delegation.Height,

		UnbondingDelegation: document.UnbondingDelegation{
			CreationHeight: ud.CreationHeight,
			MinTime:        ud.MinTime,
			InitialBalance: ud.InitialBalance,
			Balance:        ud.Balance,
		},
	}

	if delegator.BondedHeight < 0 &&
		delegator.UnbondingDelegation.CreationHeight < 0 {
		store.Delete(delegator)
		logger.Info("delete delegator", logger.String("delAddress", delAddress), logger.String("valAddress", valAddress))
	} else {
		store.SaveOrUpdate(delegator)
		logger.Info("saveOrUpdate delegator", logger.String("delAddress", delAddress), logger.String("valAddress", valAddress))
	}
}

func BuildDelegation(delAddress, valAddress string) (res tempDelegation) {
	d := helper.GetDelegation(delAddress, valAddress)

	if d.DelegatorAddress == nil {
		// represents delegation is nil
		res.Height = -1
		return res
	}

	floatShares := helper.ParseFloat(d.Shares.String())
	res = tempDelegation{
		Shares:         floatShares,
		OriginalShares: d.Shares.String(),
		Height:         -1,
	}

	return res
}

func BuildUnbondingDelegation(delAddress, valAddress string) (res document.UnbondingDelegation) {
	ud := helper.GetUnbondingDelegation(delAddress, valAddress)

	// doesn't have unbonding delegation
	if ud.DelegatorAddress == nil {
		// represents unbonding delegation is nil
		res.CreationHeight = -1
		return res
	}

  initBalances := ""
  balances := ""
  creationHeight := ud.Entries[0].CreationHeight
  completionTime := ud.Entries[0].CompletionTime
	for _, entry := range ud.Entries {
		initBalances += fmt.Sprintf("%vatom,", entry.InitialBalance)
		balances += fmt.Sprintf("%vatom,", entry.Balance)
	}
  initBalance := types.ParseCoins(initBalances[:len(initBalances)-1])
  balance := types.ParseCoins(balances[:len(balances)-1])
	res = document.UnbondingDelegation{
		CreationHeight: creationHeight,
		MinTime:        completionTime.Unix(),
		InitialBalance: initBalance,
		Balance:        balance,
	}

	return res
}

// Delegation represents the bond with tokens held by an account.  It is
// owned by one delegator, and is associated with the voting power of one
// pubKey.
type tempDelegation struct {
	Shares         float64
	OriginalShares string
	Height         int64 // Last height bond updated
}
