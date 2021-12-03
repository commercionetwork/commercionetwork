package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetDid(goCtx context.Context, msg *types.MsgSetDid) (*types.MsgSetDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	t := ctx.BlockTime()

	var DidDocument types.DidDocument

	if k.HasIdentity(ctx, msg.ID) {
		DidDocument, _ = k.Keeper.GetDdoByOwner(ctx, sdk.AccAddress(msg.ID))

		// update fields
		//
		//

		// update the timestamp for the fields that must be updated
		DidDocument.Updated = &t // &ctx.BlockTime()
	}

	DidDocument = types.DidDocument{
		Context:              msg.Context,
		ID:                   msg.ID,
		VerificationMethod:   msg.VerificationMethod,
		Service:              msg.Service,
		Authentication:       msg.Authentication,
		AssertionMethod:      msg.AssertionMethod,
		CapabilityDelegation: msg.CapabilityDelegation,
		CapabilityInvocation: msg.CapabilityInvocation,
		KeyAgreement:         msg.KeyAgreement,
		Created:              &t,
		Updated:              &t,
	}

	id := k.AppendDid(ctx, DidDocument)

	return &types.MsgSetDidResponse{
		ID: id,
	}, nil
}
