package accreditations

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/memberships"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper, membershipKeeper memberships.Keeper, governmentKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgInviteUser:
			return handleMsgInviteUser(ctx, keeper, membershipKeeper, msg)
		case MsgSetUserVerified:
			return handleMsgSetUserVerified(ctx, keeper, msg)
		case MsgDepositIntoLiquidityPool:
			return handleMsgDepositIntoPool(ctx, keeper, msg)
		case MsgAddTrustedSigner:
			return handleMsgAddTrustedSigner(ctx, keeper, governmentKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgInviteUser(ctx sdk.Context, keeper Keeper, membershipKeeper memberships.Keeper, msg MsgInviteUser) sdk.Result {
	// Verify that the user that is inviting has already a membership
	if _, found := membershipKeeper.GetMembership(ctx, msg.Sender); !found {
		return sdk.ErrUnauthorized("Cannot send an invitation without having a membership").Result()
	}

	// Try inviting the user
	if err := keeper.InviteUser(ctx, msg.Recipient, msg.Sender); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgSetUserVerified(ctx sdk.Context, keeper Keeper, msg MsgSetUserVerified) sdk.Result {

	// Check the accreditation
	if !keeper.IsTrustedServiceProvider(ctx, msg.Verifier) {
		msg := fmt.Sprintf("User %s is not a valid TSP", msg.Verifier.String())
		return sdk.ErrUnauthorized(msg).Result()
	}

	// Create a credentials and store it
	credential := types.Credential{Timestamp: msg.Timestamp, User: msg.User, Verifier: msg.Verifier}
	keeper.SaveCredential(ctx, credential)

	return sdk.Result{}
}

func handleMsgDepositIntoPool(ctx sdk.Context, keeper Keeper, msg MsgDepositIntoLiquidityPool) sdk.Result {
	if err := keeper.DepositIntoPool(ctx, msg.Depositor, msg.Amount); err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}

	return sdk.Result{}
}

func handleMsgAddTrustedSigner(ctx sdk.Context, keeper Keeper, governmentKeeper government.Keeper, msg MsgAddTrustedSigner) sdk.Result {
	if !governmentKeeper.GetGovernmentAddress(ctx).Equals(msg.Government) {
		return sdk.ErrInvalidAddress("invalid government address").Result()
	}

	keeper.AddTrustedServiceProvider(ctx, msg.Tsp)
	return sdk.Result{}
}
