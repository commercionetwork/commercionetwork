package memberships

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/accreditations"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(stableCreditDenom string, keeper keeper.Keeper,
	accreditationKeeper accreditations.Keeper, bankKeeper bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgBuyMembership:
			return handleMsgBuyMembership(ctx, stableCreditDenom, keeper, accreditationKeeper, bankKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

var membershipCosts = map[string]int64{
	types.MembershipTypeBronze: 25,
	types.MembershipTypeSilver: 250,
	types.MembershipTypeGold:   2500,
	types.MembershipTypeBlack:  25000,
}

// handleMsgBuyMembership allows to handle a MsgBuyMembership message.
// In order to be able to buy a membership the following requirements must be met.
// 1. The user has been invited from a member already having a membership
// 2. The user has been verified from a TSP
// 3. The membership must be valid
// 4. The user has enough stable credits in his wallet
func handleMsgBuyMembership(ctx sdk.Context, stableCreditsDenom string,
	keeper keeper.Keeper, accKeeper accreditations.Keeper, bankKeeper bank.Keeper,
	msg types.MsgBuyMembership) sdk.Result {

	// -------------
	// --- 1. Check the invitation and the invitee membership type

	invitation, found := accKeeper.GetInvite(ctx, msg.Buyer)
	if !found {
		return sdk.ErrUnauthorized("Cannot buy a membership without being invited").Result()
	}

	if _, found := keeper.GetMembership(ctx, invitation.Sender); !found {
		msg := "Cannot buy a membership after being invited by a user having a Green membership"
		return sdk.ErrUnauthorized(msg).Result()
	}

	// -------------
	// --- 2. Make sure the user has properly being verified

	if credentials := accKeeper.GetUserCredentials(ctx, msg.Buyer); len(credentials) == 0 {
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

	// -------------
	// --- 4. Get the tokens from the user account
	membershipPrice := membershipCosts[msg.MembershipType] * 1000000 // Always multiply by one million
	membershipCost := sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, membershipPrice))
	if _, err := bankKeeper.SubtractCoins(ctx, msg.Buyer, membershipCost); err != nil {
		return err.Result()
	}

	// -------------
	// --- Assign the membership
	if _, err := keeper.AssignMembership(ctx, msg.Buyer, msg.MembershipType); err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	return sdk.Result{}
}
