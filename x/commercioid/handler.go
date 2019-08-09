package commercioid

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "commercioid" type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized commercioid message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSetIdentity allows to handle a MsgSetIdentity checking that the user that wants to set an identity is
// the real owner of that identity.
// If the user is not allowed to use that identity, returns an error.
func handleMsgSetIdentity(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetIdentity) sdk.Result {
	keeper.SaveIdentity(ctx, msg.Owner, msg.DidDocumentUri)
	return sdk.Result{}
}
