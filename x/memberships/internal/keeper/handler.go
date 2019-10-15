package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper, governmentKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgInviteUser:
			return handleMsgInviteUser(ctx, keeper, msg)
		case types.MsgSetUserVerified:
			return handleMsgSetUserVerified(ctx, keeper, msg)
		case types.MsgDepositIntoLiquidityPool:
			return handleMsgDepositIntoPool(ctx, keeper, msg)
		case types.MsgAddTsp:
			return handleMsgAddTrustedSigner(ctx, keeper, governmentKeeper, msg)
		case types.MsgBuyMembership:
			return handleMsgBuyMembership(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgInviteUser(ctx sdk.Context, keeper Keeper, msg types.MsgInviteUser) sdk.Result {
	// Verify that the user that is inviting has already a membership
	if _, found := keeper.GetMembership(ctx, msg.Sender); !found {
		return sdk.ErrUnauthorized("Cannot send an invitation without having a membership").Result()
	}

	// Try inviting the user
	if err := keeper.InviteUser(ctx, msg.Recipient, msg.Sender); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgSetUserVerified(ctx sdk.Context, keeper Keeper, msg types.MsgSetUserVerified) sdk.Result {

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

func handleMsgDepositIntoPool(ctx sdk.Context, keeper Keeper, msg types.MsgDepositIntoLiquidityPool) sdk.Result {
	if err := keeper.DepositIntoPool(ctx, msg.Depositor, msg.Amount); err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}

	return sdk.Result{}
}

func handleMsgAddTrustedSigner(ctx sdk.Context, keeper Keeper, governmentKeeper government.Keeper, msg types.MsgAddTsp) sdk.Result {
	if !governmentKeeper.GetGovernmentAddress(ctx).Equals(msg.Government) {
		return sdk.ErrInvalidAddress("invalid government address").Result()
	}

	keeper.AddTrustedServiceProvider(ctx, msg.Tsp)
	return sdk.Result{}
}

// handleMsgBuyMembership allows to handle a MsgBuyMembership message.
// In order to be able to buy a membership the following requirements must be met.
// 1. The user has been invited from a member already having a membership
// 2. The user has been verified from a TSP
// 3. The membership must be valid
// 4. The user has enough stable credits in his wallet
func handleMsgBuyMembership(ctx sdk.Context, keeper Keeper, msg types.MsgBuyMembership) sdk.Result {

	// -------------
	// --- 1. Check the invitation and the invitee membership type

	_, found := keeper.GetInvite(ctx, msg.Buyer)
	if !found {
		return sdk.ErrUnauthorized("Cannot buy a membership without being invited").Result()
	}

	// -------------
	// --- 2. Make sure the user has properly being verified

	if credentials := keeper.GetUserCredentials(ctx, msg.Buyer); len(credentials) == 0 {
		msg := "User has not yet been verified by a Trusted Service Provider"
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// -------------
	// --- 3. Verify the membership validity

	// Check the type
	if !types.IsMembershipTypeValid(msg.MembershipType) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid membership type: %s", msg.MembershipType)).Result()
	}

	// Get the current membership
	membership, found := keeper.GetMembership(ctx, msg.Buyer)
	if found && !types.CanUpgrade(keeper.GetMembershipType(membership), msg.MembershipType) {
		errMsg := fmt.Sprintf("Cannot upgrade from %s membership to %s", keeper.GetMembershipType(membership), msg.MembershipType)
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	if err := keeper.BuyMembership(ctx, msg.Buyer, msg.MembershipType); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
