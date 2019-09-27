package pricefeed

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S START
}

//EndBlocker ensures that prices will update at most once per block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.SetCurrentPrices(ctx)
	_ = fmt.Sprintf("error is occurred: \n %s", err)
}
