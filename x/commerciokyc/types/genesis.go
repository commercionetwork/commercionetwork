package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		LiquidityPoolAmount: sdk.Coins(nil),
		Invites:             []*Invite(nil),
		Memberships:         []*Membership(nil),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate
	coins := gs.LiquidityPoolAmount

	if coins.IsAnyNegative() {
		return errors.New("liquidity pool amount cannot contain negative values")
	}
	return nil
}
