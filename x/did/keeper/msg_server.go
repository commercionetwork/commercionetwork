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

// SetDidDocument
func (k msgServer) SetDidDocument(goCtx context.Context, msg *types.MsgSetDidDocument) (*types.MsgSetDidDocumentResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	timestamp := ctx.BlockTime().Format(types.ComplaintW3CTime)

	ddo := types.DidDocument{
		Context:              msg.Context,
		ID:                   msg.ID,
		VerificationMethod:   msg.VerificationMethod,
		Service:              msg.Service,
		Authentication:       msg.Authentication,
		AssertionMethod:      msg.AssertionMethod,
		CapabilityDelegation: msg.CapabilityDelegation,
		CapabilityInvocation: msg.CapabilityInvocation,
		KeyAgreement:         msg.KeyAgreement,
	}

	previousDDO, err := k.GetDidDocumentOfAddress(ctx, msg.ID)
	if err != nil {
		ddo.Created = timestamp
	} else {
		ddo.Created = previousDDO.Created
	}
	ddo.Updated = timestamp

	id := k.UpdateDidDocument(ctx, ddo)

	return &types.MsgSetDidDocumentResponse{ID: id}, nil
}
