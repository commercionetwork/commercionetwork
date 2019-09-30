package mint

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S BEGIN
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S END
}
