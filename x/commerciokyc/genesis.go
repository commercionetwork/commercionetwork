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

	// Import all the invites
	for _, invite := range data.Invites {
		keeper.SaveInvite(ctx, *invite)
	}

	// Import the memberships
	for _, membership := range data.Memberships {
		err := keeper.AssignMembership(ctx, *membership)
		if err != nil {
			panic(err)
		}
	}

	// Import the signers
	for _, signer := range data.TrustedServiceProviders {
		keeper.AddTrustedServiceProvider(ctx, sdk.AccAddress(signer))
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
	for _, tsp := range k.GetTrustedServiceProviders(ctx) {
		trustedServiceProviders = append(trustedServiceProviders, tsp.String())
	}

	//var invites []*types.Invite
	/*for _, invite := range k.GetInvites(ctx) {
		invites = append(invites, &invite)
	}*/

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
