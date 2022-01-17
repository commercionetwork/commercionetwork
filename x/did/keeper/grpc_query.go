package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Identity(c context.Context, req *types.QueryResolveDidDocumentRequest) (*types.QueryResolveDidDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	didDocument, err := k.GetDidDocumentOfAddress(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.QueryResolveDidDocumentResponse{DidDocument: &didDocument}, nil
}
