package v2_2_0

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "commerciomint"
	BondDenom    = "ucommercio"
	CreditsDenom = "uccc"
)

type Position struct {
	Owner        sdk.AccAddress `json:"owner"`
	Collateral   math.Int        `json:"collateral"`
	Credits      sdk.Coin       `json:"credits"`
	CreatedAt    time.Time      `json:"created_at"`
	ID           string         `json:"id"`
	ExchangeRate math.LegacyDec `json:"exchange_rate"`
}

// GenesisState - commerciomint genesis state
type GenesisState struct {
	Positions           []Position    `json:"positions"`
	LiquidityPoolAmount sdk.Coins     `json:"pool_amount"`
	CollateralRate      math.LegacyDec `json:"collateral_rate"`
	FreezePeriod        time.Duration `json:"freeze_period"`
}
