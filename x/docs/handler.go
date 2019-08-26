package docs

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgShareDocument:
			return handleShareDocument(ctx, keeper, msg)
		case MsgDocumentReceipt:
			return handleDocumentReceipt(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ----------------------------------
// --- ShareDocument
// ----------------------------------

func handleShareDocument(ctx sdk.Context, keeper Keeper, msg MsgShareDocument) sdk.Result {
	keeper.ShareDocument(ctx, types.Document(msg))
	return sdk.Result{}
}

// ----------------------------------
// --- DocumentReceipt
// ----------------------------------

func handleDocumentReceipt(ctx sdk.Context, keeper Keeper, msg MsgDocumentReceipt) sdk.Result {
	keeper.ShareDocumentReceipt(ctx, types.DocumentReceipt(msg))
	return sdk.Result{}
}
