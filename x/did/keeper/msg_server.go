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

func (k msgServer) SetDid(goCtx context.Context, msg *types.MsgSetDid) (*types.MsgSetDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	t := ctx.BlockTime()

	var didDocumentNew types.DidDocumentNew

	if k.HasIdentity(ctx, msg.ID) {
		didDocumentNew, _ = k.Keeper.GetDdoByOwner(ctx, sdk.AccAddress(msg.ID))

		// update fields
		//
		//

		// update the timestamp for the fields that must be updated
		didDocumentNew.Updated = &t // &ctx.BlockTime()
	}

	didDocumentNew = types.DidDocumentNew{
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

// GetDdoByOwner returns the DID Document reference associated to a given DID.
// If the given DID has no DID Document reference associated, returns nil.
func (k Keeper) GetDdoByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocumentNew, error) {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocumentNew{}, fmt.Errorf("did document with owner %s not found", owner.String())
	}

	var didDocumentNew types.DidDocumentNew
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocumentNew)
	return didDocumentNew, nil
}
