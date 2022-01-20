package keeper

import (
	"context"
	"fmt"

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

// SetIdentity
func (k msgServer) UpdateIdentity(goCtx context.Context, msg *types.MsgSetDidDocument) (*types.MsgSetDidDocumentResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	timestamp := ctx.BlockTime().Format(types.ComplaintW3CTime)

	identity := types.Identity{
		DidDocument: msg.DidDocument,
		Metadata: &types.Metadata{
			Updated: timestamp,
		},
	}

	previousIdentity, err := k.GetLastIdentityOfAddress(ctx, msg.DidDocument.ID)
	if err != nil {
		// create new identity
		identity.Metadata.Created = timestamp
	} else {
		// use last identity info
		if msg.DidDocument.Equal(previousIdentity.DidDocument) {
			return nil, fmt.Errorf("cannot update the identity with the same DID document as the last one stored")
		}
		identity.Metadata.Created = previousIdentity.Metadata.Created
	}

	k.SetIdentity(ctx, identity)

	return &types.MsgSetDidDocumentResponse{}, nil
}
