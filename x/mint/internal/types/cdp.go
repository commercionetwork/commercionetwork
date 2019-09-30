package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//CDP stands for Collateralized Debt Position
type CDP struct {
	Owner           sdk.AccAddress `json:"owner"`
	DepositedAmount sdk.Coins      `json:"deposited_amount"`
	LiquidityAmount sdk.Coins      `json:"liquidity_amount"`
	Timestamp       *big.Int       `json:"timestamp"`
}

func (current CDP) Equals(cdp CDP) bool {
	return current.Owner.Equals(cdp.Owner) &&
		current.DepositedAmount.IsEqual(cdp.DepositedAmount) &&
		current.LiquidityAmount.IsEqual(cdp.LiquidityAmount) &&
		current.Timestamp == cdp.Timestamp
}

type CDPs []CDP

func (cdps CDPs) AppendIfMissing(cdp CDP) (CDPs, bool) {
	for _, ele := range cdps {
		if ele.Equals(cdp) {
			return nil, true
		}
	}
	return append(cdps, cdp), false
}
