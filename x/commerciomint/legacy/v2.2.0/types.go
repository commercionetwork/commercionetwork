package v2_2_0

import (
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = types.ModuleName
)

type Position struct {
	Owner        sdk.AccAddress `json:"owner"`
	Collateral   sdk.Int        `json:"collateral"`
	Credits      sdk.Coin       `json:"credits"`
	CreatedAt    time.Time      `json:"created_at"`
	ID           string         `json:"id"`
	ExchangeRate sdk.Dec        `json:"exchange_rate"`
}

// GenesisState - commerciomint genesis state
type GenesisState struct {
	Positions           []Position    `json:"positions"`
	LiquidityPoolAmount sdk.Coins     `json:"pool_amount"`
	CollateralRate      sdk.Dec       `json:"collateral_rate"`
	FreezePeriod        time.Duration `json:"freeze_period"`
}
