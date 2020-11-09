package memberships

import (
	"errors"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

// GenesisState - accreditations genesis state
type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins         `json:"liquidity_pool_amount"`     // Liquidity pool from which to get the rewards
	Invites                 types.Invites     `json:"invites"`                   // List of invites
	TrustedServiceProviders ctypes.Addresses  `json:"trusted_service_providers"` // List of trusted service providers
	Memberships             types.Memberships `json:"memberships"`               // List of all the existing memberships
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(stableCreditsDenom string) GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, supplyKeeper supply.Keeper, data GenesisState) {
	moduleAcc := keeper.GetMembershipModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(data.LiquidityPoolAmount); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	// Import the signers
	for _, signer := range data.TrustedServiceProviders {
		keeper.AddTrustedServiceProvider(ctx, signer)
	}

	// Import all the invites
	for _, invite := range data.Invites {
		keeper.SaveInvite(ctx, invite)
	}

	// Import the memberships
	for _, membership := range data.Memberships {
		err := keeper.AssignMembership(ctx, membership.Owner, membership.MembershipType, membership.TspAddress, membership.ExpiryAt)
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	// create the Memberships set
	height := ctx.BlockHeight()

	return GenesisState{
		LiquidityPoolAmount:     keeper.GetPoolFunds(ctx),
		Invites:                 keeper.GetInvites(ctx),
		TrustedServiceProviders: keeper.GetTrustedServiceProviders(ctx),
		Memberships:             keeper.GetExportMemberships(ctx, height),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.LiquidityPoolAmount.IsAnyNegative() {
		return errors.New("liquidity pool amount cannot contain negative values")
	}

	return nil
}
