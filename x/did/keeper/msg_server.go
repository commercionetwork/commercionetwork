package keeper

import (
	"context"
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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

// UpdateIdentity updates an Identity using the current block time for the Metadata Updated field
// If there is no Identity associated to the DID document ID, the Metadata Created field is set with the current block time
// Otherwise, the timestamp contained in the last Identity is used
// If the DID document in the message is the same one as the one contained in the last Identity, returns an error
func (k msgServer) UpdateIdentity(goCtx context.Context, msg *types.MsgSetIdentity) (*types.MsgSetIdentityResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

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
	ctypes.EmitCommonEvents(ctx, msg.DidDocument.ID)

	return &types.MsgSetIdentityResponse{}, nil
}
