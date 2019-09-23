//In this file, according to cosmos sdk modules should be declared Begin and End Blocker functions
package id

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S BEGIN
}

func EndBlocker(ctx sdk.Context, k Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S END
}
