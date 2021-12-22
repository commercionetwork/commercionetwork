package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	timestamp := obtainTimestamp(ctx)

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

	if !k.HasDidDocument(ctx, msg.ID) {
		ddo.Created = timestamp
		// ddo.Updated = NO // "The updated property is omitted if an Update operation has never been performed on the DID document"
	} else {
		previousDDO, err := k.GetDidDocumentOfAddress(ctx, msg.ID)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot update DDO: %e", err)
		}
		ddo.Created = previousDDO.Created
		ddo.Updated = timestamp
	}

	id := k.UpdateDidDocument(ctx, ddo)

	return &types.MsgSetDidDocumentResponse{ID: id}, nil
}
