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

// AppendDid appends a didDocument in the store with given id
func (k Keeper) AppendDid(ctx sdk.Context, DidDocument types.DidDocument) string {
	// Create the Document
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(sdk.AccAddress(DidDocument.ID)), k.cdc.MustMarshalBinaryBare(&DidDocument))
	return DidDocument.ID
}

// GetDdoByOwner returns the DID Document reference associated to a given DID.
// If the given DID has no DID Document reference associated, returns nil.
func (k Keeper) GetDdoByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, fmt.Errorf("did document with owner %s not found", owner.String())
	}

	var DidDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &DidDocument)
	return DidDocument, nil
}
