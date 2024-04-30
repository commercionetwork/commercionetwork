package government

// import (
// 	"fmt"

// 	"github.com/commercionetwork/commercionetwork/x/government/keeper"
// 	"github.com/commercionetwork/commercionetwork/x/government/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	errorsmod "cosmossdk.io/errors"
// )

// // NewHandler ...
// func NewHandler(k keeper.Keeper) sdk.Handler {
// 	msgServer := keeper.NewMsgServerImpl(k)
// 	_ = msgServer

// 	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
// 		ctx = ctx.WithEventManager(sdk.NewEventManager())

// 		switch msg := msg.(type) {
// 		default:
// 			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
// 			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
// 		}
// 	}
// }
