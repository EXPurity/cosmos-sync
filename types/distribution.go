package types

import (
	"encoding/json"
	"cosmos-sync/util/constant"
)

type WithdrawDelegatorRewardsAllMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
}

func NewWithdrawDelegatorRewardsAllMsg(msg MsgWithdrawDelegatorRewardsAll) WithdrawDelegatorRewardsAllMsg {
	return WithdrawDelegatorRewardsAllMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
	}
}

func (s WithdrawDelegatorRewardsAllMsg) Type() string {
	return constant.TxTypeWithdrawDelegatorRewardsAll
}

func (s WithdrawDelegatorRewardsAllMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type WithdrawDelegatorRewardMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawDelegatorRewardMsg(msg MsgWithdrawDelegatorReward) WithdrawDelegatorRewardMsg {
	return WithdrawDelegatorRewardMsg{
		DelegatorAddr: msg.DelegatorAddress.String(),
		ValidatorAddr: msg.ValidatorAddress.String(),
	}
}

func (s WithdrawDelegatorRewardMsg) Type() string {
	return constant.TxTypeWithdrawDelegatorReward
}

func (s WithdrawDelegatorRewardMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type WithdrawValidatorRewardsAllMsg struct {
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawValidatorRewardsAllMsg(msg MsgWithdrawValidatorCommission) WithdrawValidatorRewardsAllMsg {
	return WithdrawValidatorRewardsAllMsg{
		ValidatorAddr: msg.ValidatorAddress.String(),
	}
}

func (s WithdrawValidatorRewardsAllMsg) Type() string {
	return constant.TxTypeWithdrawValidatorRewardsAll
}

func (s WithdrawValidatorRewardsAllMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}
