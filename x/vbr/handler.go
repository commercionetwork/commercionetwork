package vbr

// import (
// 	"fmt"

// 	"github.com/commercionetwork/commercionetwork/x/vbr/keeper"
// 	"github.com/commercionetwork/commercionetwork/x/vbr/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	errorsmod "cosmossdk.io/errors"
// )

// // NewHandler ...
// func NewHandler(k keeper.Keeper) sdk.Handler {
// 	msgServer := keeper.NewMsgServerImpl(k)

// 	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
// 		ctx = ctx.WithEventManager(sdk.NewEventManager())

// 		switch msg := msg.(type) {
// 		case *types.MsgIncrementBlockRewardsPool:
// 			res, err := msgServer.IncrementBlockRewardsPool(sdk.WrapSDKContext(ctx), msg)
// 			return sdk.WrapServiceResult(ctx, res, err)
// 		case *types.MsgSetParams:
// 			res, err := msgServer.SetParams(sdk.WrapSDKContext(ctx), msg)
// 			return sdk.WrapServiceResult(ctx, res, err)
// 		default:
// 			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
// 			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
// 		}
// 	}
// }
