package types

import (
	"encoding/json"
	"cosmos-sync/util/constant"
)

type BeginRedelegate struct {
	DelegatorAddr    string `json:"delegator_addr"`
	ValidatorSrcAddr string `json:"validator_src_addr"`
	ValidatorDstAddr string `json:"validator_dst_addr"`
	SharesAmount     string `json:"shares_amount"`
}

func (s BeginRedelegate) Type() string {
	return constant.TxTypeBeginRedelegate
}

func (s BeginRedelegate) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func NewBeginRedelegate(msg MsgBeginRedelegate) BeginRedelegate {
	shares := msg.Amount.String()
	return BeginRedelegate{
		DelegatorAddr:    msg.DelegatorAddress.String(),
		ValidatorSrcAddr: msg.ValidatorSrcAddress.String(),
		ValidatorDstAddr: msg.ValidatorDstAddress.String(),
		SharesAmount:     shares,
	}
}
