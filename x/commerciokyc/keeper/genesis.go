package keeper

import (
	"fmt"
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets commerciokyc information for genesis.
// TODO move all keeper invocation in keeper package
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {

	// Setup params
	params := types.Params{
		CheckMembershipsEpochIdentifier: data.Params.CheckMembershipsEpochIdentifier,
	}

	if err := k.UpdateParams(ctx, params); err != nil {
		panic(err)
	}

	// Get the module account
	moduleAcc := k.GetMembershipModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// Get the initial pool coins
	// TODO RESOLVE POOL ISSUE
	if k.GetModuleBalance(ctx, moduleAcc.GetAddress()).IsZero() {
		if err := k.SetLiquidityPoolToAccount(ctx, data.LiquidityPoolAmount); err != nil {
			panic(err)
		}
		k.SetModuleAccount(ctx, moduleAcc)
	}

	// Import the signers
	for _, signer := range data.TrustedServiceProviders {
		tsp, err := sdk.AccAddressFromBech32(signer)
		if err != nil {
			panic(err)
		}
		k.AddTrustedServiceProvider(ctx, tsp)
	}

	// Import all the invites
	for _, invite := range data.Invites {
		k.SaveInvite(ctx, *invite)
	}

	// Import the memberships
	for _, membership := range data.Memberships {
		mOwner, _ := sdk.AccAddressFromBech32(membership.Owner)
		mTsp, _ := sdk.AccAddressFromBech32(membership.TspAddress)
		// TODO need remove membership before init
		if time.Now().After(*membership.ExpiryAt) {
			continue
		}
		err := k.AssignMembership(ctx, mOwner, membership.MembershipType, mTsp, *membership.ExpiryAt)
		if err != nil {
			panic(err)
		}
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var trustedServiceProviders []string
	for _, tsp := range k.GetTrustedServiceProviders(ctx).Addresses {
		trustedServiceProviders = append(trustedServiceProviders, tsp)
	}

	return &types.GenesisState{
		LiquidityPoolAmount:     k.GetPoolFunds(ctx),
		Invites:                 k.GetInvites(ctx),
		TrustedServiceProviders: trustedServiceProviders,
		Memberships:             k.GetMemberships(ctx),
		Params:                  k.GetParams(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state types.GenesisState) error {
	return state.Validate()
}
