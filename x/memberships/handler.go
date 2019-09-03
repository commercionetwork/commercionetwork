package memberships

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgAssignMembership:
			return handleMsgAssignMembership(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgAssignMembership allows to handle a MsgAssignMembership checking that the user that wants to set an
// identity is the real owner of that identity.
// If the user is not allowed to use that identity, returns an error.
func handleMsgAssignMembership(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAssignMembership) sdk.Result {
	// Check the signer
	if !keeper.GetTrustedMinters(ctx).Contains(msg.Signer) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid minter address: %s", msg.Signer.String())).Result()
	}

	// Check the type
	if !types.IsMembershipTypeValid(msg.MembershipType) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid membership type: %s", msg.MembershipType)).Result()
	}

	// Get the current membership
	membership, found := keeper.GetMembership(ctx, msg.User)
	if found && !types.CanUpgrade(keeper.GetMembershipType(membership), msg.MembershipType) {
		errMsg := fmt.Sprintf("Cannot upgrade from %s membership to %s", keeper.GetMembershipType(membership), msg.MembershipType)
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	// Assign the membership
	if _, err := keeper.AssignMembership(ctx, msg.User, msg.MembershipType); err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	return sdk.Result{}
}
