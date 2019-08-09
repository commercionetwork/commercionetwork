package commerciodocs

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgShareDocument:
			return handleShareDocument(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized commerciodocs message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ----------------------------------
// --- ShareDocument
// ----------------------------------

func handleShareDocument(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgShareDocument) sdk.Result {
	// Share the document
	err := keeper.ShareDocument(ctx, msg.Document)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
