package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {

	return &GenesisState{
		LiquidityPoolAmount: sdk.Coins(nil),
		Invites:             []*Invite(nil),
		Memberships:         []*Membership(nil),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	coins := gs.LiquidityPoolAmount

	if coins.IsAnyNegative() {
		return errors.New("liquidity pool amount cannot contain negative values")
	}

	for _, invite := range gs.Invites {
		if err := invite.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid invite: %s", err)
		}
	}

	for _, membership := range gs.Memberships {
		if err := membership.Validate(); err != nil {
			return fmt.Errorf("invalid invite: %s", err)
		}
	}

	return nil
}
