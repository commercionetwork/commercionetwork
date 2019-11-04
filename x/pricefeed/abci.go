package pricefeed

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//EndBlocker ensures that prices will update at most once per block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	if err := k.ComputeAndUpdateCurrentPrices(ctx); err != nil {
		fmt.Printf("Error during updating of current prices: %s", err)
	}
}
