package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CDP struct {
	Owner           sdk.AccAddress `json:"owner"`
	DepositedAmount sdk.Coins      `json:"deposited_amount"`
	LiquidityAmount sdk.Coins      `json:"liquidity_amount"`
	Timestamp       *big.Int       `json:"timestamp"`
}
