package keeper

import (
	"context"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context, Keeper, bankKeeper.Keeper) {
	keeper, ctx := SetupKeeper(t)

	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx), *keeper, keeper.bankKeeper
}
