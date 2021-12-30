package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// get the proposer public key for this block
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress {
	return k.distKeeper.GetPreviousProposerConsAddr(ctx)
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, addr sdk.ConsAddress) {
	k.distKeeper.SetPreviousProposerConsAddr(ctx, addr)
}
