package custombank

import (
	"fmt"

	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(h sdk.Handler, k Keeper, govKeeper governmentKeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case bank.MsgSend:
			return handleMsgSend(ctx, k, h, msg)

		case bank.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, h, msg)

		case MsgBlockAddressSend:
			return handleMsgBlockAddressSend(ctx, k, govKeeper, msg)

		case MsgUnlockAddressSend:
			return handleMsgUnlockAddressSend(ctx, k, govKeeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, h sdk.Handler, msg bank.MsgSend) (*sdk.Result, error) {

	// Check if address is blocked
	if k.IsAddressBlocked(ctx, msg.FromAddress) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("Account %s is blocked", msg.FromAddress.String()))
	}

	return h(ctx, msg)
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k Keeper, h sdk.Handler, msg bank.MsgMultiSend) (*sdk.Result, error) {

	// Check if the sender is blocked
	for _, out := range msg.Outputs {
		if k.IsAddressBlocked(ctx, out.Address) {
			return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("Account %s is blocked", out.Address.String()))
		}
	}

	return h(ctx, msg)
}

// Handle MsgBlockAccountSend.
func handleMsgBlockAddressSend(ctx sdk.Context, k Keeper, govKeeper governmentKeeper.Keeper, msg MsgBlockAddressSend) (*sdk.Result, error) {

	// Check the signer
	if !govKeeper.GetGovernmentAddress(ctx).Equals(msg.Signer) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, "Cannot block an address without being the government")
	}

	// Block the address
	k.AddBlockedAddresses(ctx, msg.Address)

	return &sdk.Result{}, nil
}

// Handle MsgUnlockAccountSend.
func handleMsgUnlockAddressSend(ctx sdk.Context, k Keeper, govKeeper governmentKeeper.Keeper, msg MsgUnlockAddressSend) (*sdk.Result, error) {

	// Check the signer
	if !govKeeper.GetGovernmentAddress(ctx).Equals(msg.Signer) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, "Cannot block an address without being the government")
	}

	// Block the address
	k.RemoveBlockedAddress(ctx, msg.Address)

	return &sdk.Result{}, nil
}
