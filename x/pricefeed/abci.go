package pricefeed

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/keeper"
)

//EndBlocker ensures that prices will update at most once per block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.ComputeAndUpdateCurrentPrices(ctx)
}
