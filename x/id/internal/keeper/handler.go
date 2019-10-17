package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/id"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper id.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case id.MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", id.ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSetIdentity allows to handle a MsgSetIdentity checking that the user that wants to set an identity is
// the real owner of that identity.
// If the user is not allowed to use that identity, returns an error.
func handleMsgSetIdentity(ctx sdk.Context, keeper id.Keeper, msg id.MsgSetIdentity) sdk.Result {
	keeper.SaveIdentity(ctx, msg.Owner, msg.DidDocument)
	return sdk.Result{}
}
