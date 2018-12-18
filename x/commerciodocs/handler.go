package commerciodocs

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgStoreDocument:
			return handleStoreDocument(ctx, keeper, msg)
		case MsgShareDocument:
			return handleShareDocument(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized commerciodocs message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ----------------------------------
// --- StoreDocument
// ----------------------------------

func handleStoreDocument(ctx sdk.Context, keeper Keeper, msg MsgStoreDocument) sdk.Result {

	// Checks if the the document already has an owner.
	// If it does, checks that msg sender is the same as the current owner
	if keeper.HasOwner(ctx, msg.Reference) && !keeper.IsOwner(ctx, msg.Owner, msg.Reference) {
		// If not, throw an error
		return sdk.ErrUnauthorized("The given account has no access to the document").Result()
	}

	// Checks whenever the given AccAddress is authorized to use the provided identity
	if !keeper.commercioIdKeeper.CanBeUsedBy(ctx, msg.Owner, msg.Identity) {
		return sdk.ErrUnauthorized("The provided identity cannot be used by the given account").Result()
	}

	// Store the document
	keeper.StoreDocument(ctx, msg.Owner, msg.Identity, msg.Reference, msg.Metadata)

	return sdk.Result{}
}

// ----------------------------------
// --- ShareDocument
// ----------------------------------

func handleShareDocument(ctx sdk.Context, keeper Keeper, msg MsgShareDocument) sdk.Result {

	// Checks if the the document already has an owner.
	// If it does, checks that msg sender is the same as the current owner
	if keeper.HasOwner(ctx, msg.Reference) && !keeper.IsOwner(ctx, msg.Owner, msg.Reference) {
		// If not, throw an error
		return sdk.ErrUnauthorized("The given account has no access to the document").Result()
	}

	// Checks whenever the given AccAddress is authorized to use the provided identity
	if !keeper.commercioIdKeeper.CanBeUsedBy(ctx, msg.Owner, msg.Sender) {
		return sdk.ErrUnauthorized("The provided sender identity cannot be used by the given account").Result()
	}

	// Share the document
	keeper.ShareDocument(ctx, msg.Reference, msg.Sender, msg.Receiver)

	return sdk.Result{}
}
