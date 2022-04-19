package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// InitGenesis sets commerciokyc information for genesis.
// TODO move all keeper invocation in keeper package
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {

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
		// Need use sigle keeper methods in AssignMembership to assign membership avoid expired issue
		if !types.IsMembershipTypeValid(membership.MembershipType) {
			panic(sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", membership.MembershipType)))
		}
		if k.IsTrustedServiceProvider(ctx, mOwner) && membership.MembershipType != types.MembershipTypeBlack {
			panic(sdkErr.Wrap(sdkErr.ErrUnauthorized,
				fmt.Sprintf("account \"%s\" is a Trust Service Provider: remove from tsps list before", mOwner),
			))
		}
		// Delete membership if exists
		_ = k.DeleteMembership(ctx, mOwner)

		store := ctx.KVStore(k.storeKey)
		staddr := k.storageForAddr(mOwner)
		if store.Has(staddr) {
			panic(sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"cannot add membership \"%s\" for address %s: user already has a membership",
					membership.MembershipType,
					mOwner,
				),
			))
		}
		// Save membership
		membership := types.NewMembership(membership.MembershipType, mOwner, mTsp, membership.ExpiryAt.UTC())
		store.Set(staddr, k.cdc.MustMarshal(&membership))
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
	}
}
