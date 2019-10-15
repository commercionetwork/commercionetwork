package custombank

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(h sdk.Handler, k Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case bank.MsgSend:
			return handleMsgSend(ctx, k, h, msg)

		case bank.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, h, msg)

		case MsgBlockAddressSend:
			return handleMsgBlockAddressSend(ctx, k, govKeeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, h sdk.Handler, msg bank.MsgSend) sdk.Result {

	// Check if address is blocked
	if k.IsAddressBlocked(ctx, msg.FromAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("Account %s is blocked", msg.FromAddress.String())).Result()
	}

	return h(ctx, msg)
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k Keeper, h sdk.Handler, msg bank.MsgMultiSend) sdk.Result {

	// Check if the sender is blocked
	for _, out := range msg.Outputs {
		if k.IsAddressBlocked(ctx, out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("Account %s is blocked", out.Address.String())).Result()
		}
	}

	return h(ctx, msg)
}

// Handle MsgBlockAccountSend.
func handleMsgBlockAddressSend(ctx sdk.Context, k Keeper, govKeeper government.Keeper, msg MsgBlockAddressSend) sdk.Result {

	// Check the signer
	if !govKeeper.GetGovernmentAddress(ctx).Equals(msg.Signer) {
		return sdk.ErrUnauthorized("Cannot block an address without being the government").Result()
	}

	// Block the address
	k.AddBlockedAddresses(ctx, msg.Address)

	return sdk.Result{}
}
