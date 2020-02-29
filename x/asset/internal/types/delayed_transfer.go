package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/barkisnet/barkis/types"
)

type DelayedTransfer struct {
	From        sdk.AccAddress `json:"from" yaml:"from"`
	To          sdk.AccAddress `json:"to" yaml:"to"`
	Amount      sdk.Coins      `json:"amount" yaml:"amount"`
	MaturedTime int64          `json:"matured_time" yaml:"matured_time"`
	Sequence    int64          `json:"sequence" yaml:"sequence"`
}

func NewDelayedTransfer(from, to sdk.AccAddress, amount sdk.Coins, maturedTime int64, sequence int64) *DelayedTransfer {
	return &DelayedTransfer{
		From:        from,
		To:          to,
		Amount:      amount,
		MaturedTime: maturedTime,
		Sequence:    sequence,
	}
}

func (delayedTransfer DelayedTransfer) String() string {
	return fmt.Sprintf(`DelayedTransfer:
	from:        %s
	to:          %s
	amount:      %s
	matureTime:  %s
	sequence:    %d`, delayedTransfer.From.String(), delayedTransfer.To.String(), delayedTransfer.Amount.String(),
		time.Unix(delayedTransfer.MaturedTime, 0).Format("2006-01-02T15:04:05.000Z"), delayedTransfer.Sequence)
}

type DelayedTransferList []*DelayedTransfer

func (delayedTransferList DelayedTransferList) String() (out string) {
	for _, delayedTransfer := range delayedTransferList {
		out += delayedTransfer.String() + "\n"
	}
	return strings.TrimSpace(out)
}
