package types

import (
	"fmt"
	//"cosmos-sync/conf/server"
	"cosmos-sync/logger"
	"cosmos-sync/store"
	//"cosmos-sync/util/constant"

  "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
  "github.com/cosmos/cosmos-sdk/client/context"
  "github.com/cosmos/cosmos-sdk/codec"
  "github.com/cosmos/cosmos-sdk/x/auth"
  "github.com/cosmos/cosmos-sdk/x/bank"
  "github.com/cosmos/cosmos-sdk/x/distribution"
  dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
  "github.com/cosmos/cosmos-sdk/x/gov"
  "github.com/cosmos/cosmos-sdk/x/slashing"
  "github.com/cosmos/cosmos-sdk/x/staking"
  stags "github.com/cosmos/cosmos-sdk/x/staking/tags"
  staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
  "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tm "github.com/tendermint/tendermint/types"
	"regexp"
	"strconv"
	"strings"
)

type (
	Int = types.Int
	MsgTransfer = bank.MsgSend

	MsgStakeCreate         = staking.MsgCreateValidator
	MsgStakeEdit           = staking.MsgEditValidator
	MsgStakeDelegate       = staking.MsgDelegate
	MsgStakeBeginUnbonding = staking.MsgUndelegate
	MsgBeginRedelegate     = staking.MsgBeginRedelegate
	StakeValidator         = staking.Validator
	Delegation             = staking.Delegation
	UnbondingDelegation    = staking.UnbondingDelegation

	MsgUnjail              = slashing.MsgUnjail

	//MsgSetWithdrawAddress          = distribution.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward     = distribution.MsgWithdrawDelegatorReward
	//MsgWithdrawDelegatorReward     = distribution.MsgWithdrawDelegatorRewardsAll
	MsgWithdrawValidatorCommission = distribution.MsgWithdrawValidatorCommission

	MsgDeposit                       = gov.MsgDeposit
	//MsgSubmitProposal                = gov.MsgSubmitProposal
	//MsgSubmitSoftwareUpgradeProposal = gov.MsgSubmitSoftwareUpgradeProposal
	MsgVote                          = gov.MsgVote
	Proposal                         = gov.Proposal
	SdkVote                          = gov.Vote

	ResponseDeliverTx = abci.ResponseDeliverTx

	StdTx      = auth.StdTx
	SdkCoins   = types.Coins
	KVPair     = types.KVPair
	AccAddress = types.AccAddress
	ValAddress = types.ValAddress
	Dec        = types.Dec
	Validator  = tm.Validator
	Tx         = tm.Tx
	Block      = tm.Block
	BlockMeta  = tm.BlockMeta
	HexBytes   = cmn.HexBytes
	TmKVPair   = cmn.KVPair

	ABCIQueryOptions = rpcclient.ABCIQueryOptions
	Client           = rpcclient.Client
	HTTP             = rpcclient.HTTP
	ResultStatus     = ctypes.ResultStatus
)

// msg struct for delegation withdraw for all of the delegator's delegations
type MsgWithdrawDelegatorRewardsAll struct {
	DelegatorAddr types.AccAddress `json:"delegator_addr"`
}

type MsgSubmitProposal struct {
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   gov.ProposalKind   `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       types.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit types.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
}

func NewMsgWithdrawDelegatorRewardsAll(delAddr types.AccAddress) MsgWithdrawDelegatorRewardsAll {
	return MsgWithdrawDelegatorRewardsAll{
		DelegatorAddr: delAddr,
	}
}

func (msg MsgWithdrawDelegatorRewardsAll) Route() string { return "distr" }
func (msg MsgWithdrawDelegatorRewardsAll) Type() string  { return "withdraw_delegator_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorRewardsAll) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.AccAddress(msg.DelegatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorRewardsAll) GetSignBytes() []byte {
	bz := dtypes.MsgCdc.MustMarshalJSON(msg)
	return types.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDelegatorRewardsAll) ValidateBasic() types.Error {
	if msg.DelegatorAddr.Empty() {
		return dtypes.ErrNilDelegatorAddr(dtypes.DefaultCodespace)
	}
	/*if msg.ValidatorAddress.Empty() {
		return dtypes.ErrNilValidatorAddr(dtypes.DefaultCodespace)
	}*/
	return nil
}

func NewMsgSubmitProposal(title, description string, proposalType gov.ProposalKind, proposer types.AccAddress, initialDeposit types.Coins) MsgSubmitProposal {
	return MsgSubmitProposal{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
	}
}

//nolint
func (msg MsgSubmitProposal) Route() string { return "gov" }
func (msg MsgSubmitProposal) Type() string  { return "submit_proposal" }

// Implements Msg.
func (msg MsgSubmitProposal) ValidateBasic() types.Error {
	if len(msg.Title) == 0 {
		return gov.ErrInvalidTitle(gov.DefaultCodespace, "No title present in proposal")
	}
	if len(msg.Title) > gov.MaxTitleLength {
		return gov.ErrInvalidTitle(gov.DefaultCodespace, fmt.Sprintf("Proposal title is longer than max length of %d", gov.MaxTitleLength))
	}
	if len(msg.Description) == 0 {
		return gov.ErrInvalidDescription(gov.DefaultCodespace, "No description present in proposal")
	}
	if len(msg.Description) > gov.MaxDescriptionLength {
		return gov.ErrInvalidDescription(gov.DefaultCodespace, fmt.Sprintf("Proposal description is longer than max length of %d", gov.MaxDescriptionLength))
	}
	var pt = msg.ProposalType
	if pt != gov.ProposalTypeText &&
		pt != gov.ProposalTypeParameterChange &&
		pt != gov.ProposalTypeSoftwareUpgrade {
		return gov.ErrInvalidProposalType(gov.DefaultCodespace, msg.ProposalType)
	}
	if msg.Proposer.Empty() {
		return types.ErrInvalidAddress(msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return types.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if msg.InitialDeposit.IsAnyNegative() {
		return types.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	return nil
}

func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf("MsgSubmitProposal{%s, %s, %s, %v}", msg.Title, msg.Description, msg.ProposalType, msg.InitialDeposit)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	bz := codec.New().MustMarshalJSON(msg)
	return types.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Proposer}
}

var (
	ValidatorsKey        = staking.ValidatorsKey
	GetValidatorKey      = staking.GetValidatorKey
	GetDelegationKey     = staking.GetDelegationKey
	GetDelegationsKey    = staking.GetDelegationsKey
	GetUBDKey            = staking.GetUBDKey
	GetUBDsKey           = staking.GetUBDsKey
	ValAddressFromBech32 = types.ValAddressFromBech32

	UnmarshalValidator      = staketypes.UnmarshalValidator
	MustUnmarshalValidator  = staketypes.MustUnmarshalValidator
	UnmarshalDelegation     = staketypes.UnmarshalDelegation
	MustUnmarshalDelegation = staketypes.MustUnmarshalDelegation
	MustUnmarshalUBD        = staketypes.MustUnmarshalUBD

	Bech32ifyValPub      = types.Bech32ifyValPub
	RegisterCodec        = types.RegisterCodec
	AccAddressFromBech32 = types.AccAddressFromBech32
	BondStatusToString   = types.BondStatusToString

	NewDecFromStr = types.NewDecFromStr

	AddressStoreKey   = auth.AddressStoreKey
	GetAccountDecoder = context.GetAccountDecoder

	KeyProposal      = gov.KeyProposal
	KeyVotesSubspace = gov.KeyVotesSubspace

	NewHTTP = rpcclient.NewHTTP

	//tags
	TagGovProposalID                   = "proposal-id"
	TagDistributionReward              = "withdraw-reward-total"
	TagStakeActionCompleteRedelegation = stags.ActionCompleteRedelegation
	TagStakeDelegator                  = stags.Delegator
	TagStakeSrcValidator               = stags.SrcValidator
	TagAction                          = types.TagAction

	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	/*if server.Network == constant.NetworkMainnet {
		types.SetNetworkType(types.Mainnet)
	}*/
	cdc = app.MakeCodec()
}

func GetCodec() *codec.Codec {
	return cdc
}

//
func ParseCoins(coinsStr string) (coins store.Coins) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin := ParseCoin(coinStr)
		coins = append(coins, coin)
	}
	return coins
}

func ParseCoin(coinStr string) (coin store.Coin) {
	var (
		reDnm  = `[A-Za-z\-]{2,15}`
		reAmt  = `[0-9]+[.]?[0-9]*`
		reSpc  = `[[:space:]]*`
		reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
	)

	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		logger.Error("invalid coin expression", logger.Any("coin", coinStr))
		return coin
	}
	denom, amount := matches[2], matches[1]

	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Error("Convert str to int failed", logger.Any("amount", amount))
		return coin
	}

	return store.Coin{
		Denom:  denom,
		Amount: amt,
	}
}

func BuildFee(fee auth.StdFee) store.Fee {
	return store.Fee{
		Amount: ParseCoins(fee.Amount.String()),
		Gas:    int64(fee.Gas),
	}
}
