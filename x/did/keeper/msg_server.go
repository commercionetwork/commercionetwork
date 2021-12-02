package keeper

import (
	"context"
	"time"

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

	if k.HasIdentity(ctx, msg.ID) {

	}

	var didDocumentNew = types.DidDocumentNew{
		Context:              msg.Context,
		ID:                   msg.ID,
		VerificationMethod:   []*types.VerificationMethod{},
		Service:              []*types.ServiceNew{},
		Authentication:       []*types.VerificationMethod{},
		AssertionMethod:      []*types.VerificationMethod{},
		CapabilityDelegation: []*types.VerificationMethod{},
		CapabilityInvocation: []*types.VerificationMethod{},
		KeyAgreement:         []*types.VerificationMethod{},
		Created:              &time.Time{},
		Updated:              &time.Time{},
	}

	id := k.AppendDid(ctx, didDocumentNew)

	return &types.MsgSetDidResponse{
		ID: id,
	}, nil
}

// AppendDid appends a didDocument in the store with given id
func (k Keeper) AppendDid(ctx sdk.Context, didDocumentNew types.DidDocumentNew) string {
	// Create the Document
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(sdk.AccAddress(didDocumentNew.ID)), k.cdc.MustMarshalBinaryBare(&didDocumentNew))
	return didDocumentNew.ID
}
