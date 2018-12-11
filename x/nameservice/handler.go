package nameservice

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
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetName
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {

	// Checks if the the msg sender is the same as the current owner
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.NameID)) {
		// If not, throw an error
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}

	// If so, set the name to the value specified in the msg.
	keeper.SetName(ctx, msg.NameID, msg.Value)

	// return
	return sdk.Result{}
}

// Handle MsgBuyName
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	// NOTE: This handler uses functions from the coinKeeper to perform currency operations.
	// If your application is performing currency operations you man want to take a look at the godocs for this module
	// to see what functions it exposes.

	// Checks if the the bid price is greater than the price paid by the current owner
	if keeper.GetPrice(ctx, msg.NameID).IsAllGT(msg.Bid) {
		// If not, throw an error
		return sdk.ErrInsufficientCoins("Bid not high enough").Result()
	}

	// Checks if someone owns the name
	if keeper.HasOwner(ctx, msg.NameID) {
		// If so, send the coin from the buyer to the owner of the name
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.NameID), msg.Bid)
		if err != nil {
			// If not, throw an error
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		// If not, just subtract the amount from the buyer without sending it to anyone
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.NameID, msg.Buyer)
	keeper.SetPrice(ctx, msg.NameID, msg.Bid)
	return sdk.Result{}
}