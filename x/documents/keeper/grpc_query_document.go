package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
func (k Keeper) DocumentAll(c context.Context, req *types.QueryAllDocumentRequest) (*types.QueryAllDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var documents []*types.Document
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, types.KeyPrefix(types.DocumentKey))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		var document types.Document
		if err := k.cdc.UnmarshalBinaryBare(value, &document); err != nil {
			return err
		}

		documents = append(documents, &document)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDocumentResponse{Document: documents, Pagination: pageRes}, nil
}
*/
func (k Keeper) Document(c context.Context, req *types.QueryGetDocumentRequest) (*types.QueryGetDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var document types.Document
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasDocument(ctx, req.UUID) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DocumentKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(getDocumentStoreKey(req.UUID)), &document)

	return &types.QueryGetDocumentResponse{Document: &document}, nil
}

func (k Keeper) SentDocuments(c context.Context, req *types.QueryGetSentDocumentsRequest) (*types.QueryGetSentDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var documents []*types.Document
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, types.KeyPrefix(types.SentDocumentsPrefix))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		var document types.Document

		if err := k.cdc.UnmarshalBinaryBare(value, &document); err != nil {
			return err
		}
		if document.Sender == req.Address {
			documents = append(documents, &document)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetSentDocumentsResponse{Document: documents, Pagination: pageRes}, nil
}

func (k Keeper) ReceivedDocument(c context.Context, req *types.QueryGetReceivedDocumentRequest) (*types.QueryGetReceivedDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receivedDocuments []*types.Document
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, types.KeyPrefix(types.ReceivedDocumentsPrefix))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		/*var receivedDocument types.Document
		if err := k.cdc.UnmarshalBinaryBare(value, &receivedDocument); err != nil {
			return err
		}*/
		useAddr, _ := sdk.AccAddressFromBech32(req.Address)
		rdi := k.UserReceivedDocumentsIterator(ctx, useAddr)
		defer rdi.Close()

		//documents := []types.Document{}
		for ; rdi.Valid(); rdi.Next() {
			id := string(rdi.Value())
			doc, err := k.GetDocumentByID(ctx, id)
			if err != nil {
				return err
			}

			receivedDocuments = append(receivedDocuments, &doc)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetReceivedDocumentResponse{ReceivedDocument: receivedDocuments, Pagination: pageRes}, nil
}

func (k Keeper) SentDocumentsReceipts(c context.Context, req *types.QueryGetSentDocumentsReceiptsRequest) (*types.QueryGetSentDocumentsReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receipts []*types.DocumentReceipt
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, types.KeyPrefix(types.SentDocumentsReceiptsPrefix))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		var receipt types.DocumentReceipt

		if err := k.cdc.UnmarshalBinaryBare(value, &receipt); err != nil {
			return err
		}

		if receipt.Sender == req.Address {
			receipts = append(receipts, &receipt)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetSentDocumentsReceiptsResponse{Receipt: receipts, Pagination: pageRes}, nil
}

func (k Keeper) ReceivedDocumentsReceipts(c context.Context, req *types.QueryGetReceivedDocumentsReceiptsRequest) (*types.QueryGetReceivedDocumentsReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receivedReceipts []*types.DocumentReceipt
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	documentStore := prefix.NewStore(store, types.KeyPrefix(types.ReceivedDocumentsReceiptsPrefix))

	pageRes, err := query.Paginate(documentStore, req.Pagination, func(key []byte, value []byte) error {
		userAddr, _ := sdk.AccAddressFromBech32(req.Address)
		urri := k.UserReceivedReceiptsIterator(ctx, userAddr)
		defer urri.Close()

		for ; urri.Valid(); urri.Next() {
			rid := string(urri.Value())

			r, err := k.GetReceiptByID(ctx, rid)
			if err != nil {
				return err
			}

			receivedReceipts = append(receivedReceipts, &r)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetReceivedDocumentsReceiptsResponse{ReceiptReceived: receivedReceipts, Pagination: pageRes}, nil
}
