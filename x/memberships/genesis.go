package memberships

import (
	"errors"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - accreditations genesis state
type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins         `json:"liquidity_pool_amount"`     // Liquidity pool from which to get the rewards
	Invites                 []types.Invite    `json:"invites"`                   // List of invites
	TrustedServiceProviders ctypes.Addresses  `json:"trusted_service_providers"` // List of trusted service providers
	Credentials             types.Credentials `json:"credentials"`               // List of verifiable credentials
	StableCreditsDenom      string            `json:"stable_credits_denom"`      // Stable credits denom used during membership buying
	Memberships             []Membership      `json:"memberships"`               // List of all the existing memberships
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(stableCreditsDenom string) GenesisState {
	return GenesisState{
		StableCreditsDenom: stableCreditsDenom,
	}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Set the liquidity pool
	if data.LiquidityPoolAmount != nil {
		keeper.SetPoolFunds(ctx, data.LiquidityPoolAmount)
	}

	// Import the signers
	for _, signer := range data.TrustedServiceProviders {
		keeper.AddTrustedServiceProvider(ctx, signer)
	}

	// Import all the invites
	for _, invite := range data.Invites {
		keeper.SaveInvite(ctx, invite)
	}

	// Import the credentials
	for _, credential := range data.Credentials {
		keeper.SaveCredential(ctx, credential)
	}

	// Import the memberships
	for _, membership := range data.Memberships {
		_, err := keeper.AssignMembership(ctx, membership.Owner, membership.MembershipType)
		if err != nil {
			panic(err)
		}
	}

	// Set the stable credits denom
	keeper.SetStableCreditsDenom(ctx, data.StableCreditsDenom)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		LiquidityPoolAmount:     keeper.GetPoolFunds(ctx),
		Invites:                 keeper.GetInvites(ctx),
		TrustedServiceProviders: keeper.GetTrustedServiceProviders(ctx),
		Credentials:             keeper.GetCredentials(ctx),
		Memberships:             keeper.GetMembershipsSet(ctx),
		StableCreditsDenom:      keeper.GetStableCreditsDenom(ctx),
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
