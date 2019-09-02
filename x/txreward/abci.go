package txreward

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	//Get the number of active validators
	activeValidators := k.StakeKeeper.GetLastValidators(ctx)
	valNumber := len(activeValidators)

	k.StakeKeeper

}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S END
}
