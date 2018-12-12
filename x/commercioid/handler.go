package commercioid

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
// At the moment, there is only one Msg/Handler.

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetName
func handleMsgSetIdentity(ctx sdk.Context, keeper Keeper, msg MsgSetIdentity) sdk.Result {

	// Checks if the the msg sender is the same as the current owner
	if keeper.HasOwner(ctx, msg.DID) && !msg.Owner.Equals(keeper.GetOwner(ctx, msg.DID)) {
		// If not, throw an error
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}

	// If so, set the DDO reference to the value specified in the msg.
	keeper.SetIdentity(ctx, msg.DID, msg.DDOReference)

	// Also, store the owner of the identity so one will be able to claim it or edit it
	keeper.SetOwner(ctx, msg.DID, msg.Owner)

	// return
	return sdk.Result{}
}
