package commercioid

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "commercioid" type messages.
// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetIdentity:
			return handleMsgCreateIdentity(ctx, keeper, msg)
		case types.MsgCreateConnection:
			return handleMsgCreateConnection(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized commercioid message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ----------------------------------
// --- Create identity
// ----------------------------------

func handleMsgCreateIdentity(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetIdentity) sdk.Result {

	// Checks if the the msg sender is the same as the current owner
	if !keeper.CanBeUsedBy(ctx, msg.Owner, msg.DID) {
		// If not, throw an error
		return sdk.ErrUnauthorized("Incorrect Signer").Result()
	}

	// If so, set the DDO reference to the value specified in the msg.
	keeper.CreateIdentity(ctx, msg.Owner, msg.DID, msg.DDOReference)

	// return
	return sdk.Result{}
}

// ----------------------------------
// --- Create connection
// ----------------------------------

func handleMsgCreateConnection(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateConnection) sdk.Result {

	// Checks if the the msg sender is the same as the current owner
	if !keeper.CanBeUsedBy(ctx, msg.Signer, msg.FirstUser) || !keeper.CanBeUsedBy(ctx, msg.Signer, msg.SecondUser) {
		// If not, throw an error
		return sdk.ErrUnauthorized("The signer must own either the first or the second DID").Result()
	}

	// If so, set the DDO reference to the value specified in the msg.
	keeper.AddConnection(ctx, msg.FirstUser, msg.SecondUser)

	// return
	return sdk.Result{}
}
