package pricefeed

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//EndBlocker ensures that prices will update at most once per block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.SetCurrentPrices(ctx)
	_ = fmt.Sprintf("error is occurred: \n %s", err)
}
