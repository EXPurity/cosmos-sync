// This package is used for Query balance of account

package helper

import (
	"cosmos-sync/logger"
	"cosmos-sync/store"
	"cosmos-sync/types"
	"cosmos-sync/util/constant"
)

// query account balance from sdk store
func QueryAccountBalance(address string) store.Coins {
	cdc := types.GetCodec()

	addr, err := types.AccAddressFromBech32(address)
	if err != nil {
		logger.Error("get addr from hex failed", logger.Any("err", err))
		return nil
	}

	res, err := Query(types.AddressStoreKey(addr), "acc",
		constant.StoreDefaultEndPath)

	if err != nil {
		logger.Error("Query balance from tendermint failed", logger.Any("err", err))
		return nil
	}

	// balance is empty
	if len(res) <= 0 {
		return nil
	}

	decoder := types.GetAccountDecoder(cdc)
	account, err := decoder(res)
	if err != nil {
		logger.Error("decode account failed", logger.Any("err", err))
		return nil
	}

	return types.ParseCoins(account.GetCoins().String())
}

func ValAddrToAccAddr(address string) (accAddr string) {
	valAddr, err := types.ValAddressFromBech32(address)
	if err != nil {
		logger.Error("ValAddressFromBech32 decode account failed", logger.String("address", address))
		return
	}

	return types.AccAddress(valAddr.Bytes()).String()
}
