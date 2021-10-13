package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Identity(c context.Context, req *types.QueryResolveDidRequest) (*types.QueryResolveDidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var didDocument types.DidDocument
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasIdentity(ctx, req.ID) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := ctx.KVStore(k.storeKey)
	k.cdc.MustUnmarshalBinaryBare(store.Get(getIdentityStoreKey(sdk.AccAddress(req.ID))), &didDocument)

	return &types.QueryResolveDidResponse{DidDocument: &didDocument}, nil
}
