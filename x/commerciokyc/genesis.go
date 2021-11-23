package commerciokyc

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets commerciokyc information for genesis.
// TODO move all keeper invocation in keeper package
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {

	// Get the module account
	moduleAcc := keeper.GetMembershipModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// Get the initial pool coins
	// TODO RESOLVE POOL ISSUE
	if keeper.GetModuleBalance(ctx, moduleAcc.GetAddress()).IsZero() {
		if err := keeper.SetLiquidityPoolToAccount(ctx, data.LiquidityPoolAmount); err != nil {
			panic(err)
		}
		keeper.SetModuleAccount(ctx, moduleAcc)
	}

	// Import the signers
	for _, signer := range data.TrustedServiceProviders {
		tsp, err := sdk.AccAddressFromBech32(signer)
		if err != nil {
			panic(err)
		}
		keeper.AddTrustedServiceProvider(ctx, tsp)
	}

	// Import all the invites
	for _, invite := range data.Invites {
		keeper.SaveInvite(ctx, *invite)
	}

	// Import the memberships
	for _, membership := range data.Memberships {
		mOwner, _ := sdk.AccAddressFromBech32(membership.Owner)
		mTsp, _ := sdk.AccAddressFromBech32(membership.TspAddress)
		err := keeper.AssignMembership(ctx, mOwner, membership.MembershipType, mTsp, *membership.ExpiryAt)
		if err != nil {
			panic(err)
		}
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// create the Memberships set
	/*var liquidityPoolAmount []*sdk.Coin
	for _, coin := range k.GetPoolFunds(ctx) {
		liquidityPoolAmount = append(liquidityPoolAmount, &coin)
	}*/
	var trustedServiceProviders []string
	for _, tsp := range k.GetTrustedServiceProviders(ctx).Addresses {
		trustedServiceProviders = append(trustedServiceProviders, tsp)
	}

	return &types.GenesisState{
		LiquidityPoolAmount:     k.GetPoolFunds(ctx),
		Invites:                 k.GetInvites(ctx),
		TrustedServiceProviders: trustedServiceProviders,
		Memberships:             k.GetMemberships(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state types.GenesisState) error {
	return state.Validate()
}
