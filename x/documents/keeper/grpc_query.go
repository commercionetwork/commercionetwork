package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Document(c context.Context, req *types.QueryGetDocumentRequest) (*types.QueryGetDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	document, err := k.GetDocumentByID(ctx, req.UUID)
	if err != nil {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetDocumentResponse{Document: &document}, nil
}

func (k Keeper) SentDocuments(c context.Context, req *types.QueryGetSentDocumentsRequest) (*types.QueryGetSentDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// consider using documents := []*types.Document{}
	var documents []*types.Document
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	userAddress, e := sdk.AccAddressFromBech32(req.Address)
	if e != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", req.Address))
	}
	documentStore := prefix.NewStore(store, getSentDocumentsIdsStoreKey(userAddress))

	pageRes, err := query.Paginate(
		documentStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			sentDocument, err := k.GetDocumentByID(ctx, string(value))
			if err != nil {
				return err
			}
			documents = append(documents, &sentDocument)

			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetSentDocumentsResponse{
		Document:   documents,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) UUIDDocuments(c context.Context, req *types.QueryGetUUIDDocumentsRequest) (*types.QueryGetUUIDDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var documents []string

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	userAddress, e := sdk.AccAddressFromBech32(req.Address)
	if e != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", req.Address))
	}
	documentStore := prefix.NewStore(store, getSentDocumentsIdsStoreKey(userAddress))

	pageRes, err := query.Paginate(
		documentStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			documents = append(documents, string(value))
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetUUIDDocumentsResponse{
		UUIDs:      documents,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) ReceivedDocument(c context.Context, req *types.QueryGetReceivedDocumentRequest) (*types.QueryGetReceivedDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receivedDocuments []*types.Document
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	userAddress, e := sdk.AccAddressFromBech32(req.Address)
	if e != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", req.Address))
	}
	documentStore := prefix.NewStore(store, getReceivedDocumentsIdsStoreKey(userAddress))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		receivedDocument, err := k.GetDocumentByID(ctx, string(value))
		if err != nil {
			return err
		}
		receivedDocuments = append(receivedDocuments, &receivedDocument)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetReceivedDocumentResponse{
		ReceivedDocument: receivedDocuments,
		Pagination:       pageRes,
	}, nil
}

func (k Keeper) SentDocumentsReceipts(c context.Context, req *types.QueryGetSentDocumentsReceiptsRequest) (*types.QueryGetSentDocumentsReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receipts []*types.DocumentReceipt
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	userAddress, e := sdk.AccAddressFromBech32(req.Address)
	if e != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", e))
	}
	documentStore := prefix.NewStore(store, getSentReceiptsIdsStoreKey(userAddress))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		sentReceipt, err := k.GetReceiptByID(ctx, string(value))
		if err != nil {
			return err
		}
		receipts = append(receipts, &sentReceipt)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetSentDocumentsReceiptsResponse{
		Receipt:    receipts,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) ReceivedDocumentsReceipts(c context.Context, req *types.QueryGetReceivedDocumentsReceiptsRequest) (*types.QueryGetReceivedDocumentsReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receivedReceipts []*types.DocumentReceipt
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	userAddress, e := sdk.AccAddressFromBech32(req.Address)
	if e != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", req.Address))
	}
	documentStore := prefix.NewStore(store, getReceivedReceiptsIdsStoreKey(userAddress))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		receivedReceipt, err := k.GetReceiptByID(ctx, string(value))
		if err != nil {
			return err
		}
		receivedReceipts = append(receivedReceipts, &receivedReceipt)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetReceivedDocumentsReceiptsResponse{
		ReceiptReceived: receivedReceipts,
		Pagination:      pageRes,
	}, nil
}

func (k Keeper) DocumentsReceipts(c context.Context, req *types.QueryGetDocumentsReceiptsRequest) (*types.QueryGetDocumentsReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receivedReceipts []*types.DocumentReceipt
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, getDocumentReceiptsIdsStoreKey(req.UUID))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		receivedReceipt, err := k.GetReceiptByID(ctx, string(value))
		if err != nil {
			return err
		}
		receivedReceipts = append(receivedReceipts, &receivedReceipt)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetDocumentsReceiptsResponse{
		Receipts:   receivedReceipts,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) DocumentsUUIDReceipts(c context.Context, req *types.QueryGetDocumentsUUIDReceiptsRequest) (*types.QueryGetDocumentsUUIDReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receivedReceipts []string
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, getDocumentReceiptsIdsStoreKey(req.UUID))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		receivedReceipts = append(receivedReceipts, string(value))
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryGetDocumentsUUIDReceiptsResponse{
		UUIDs:      receivedReceipts,
		Pagination: pageRes,
	}, nil
}
