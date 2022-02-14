package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gofrs/uuid"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/documents/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryReceivedDocuments:
			return queryGetReceivedDocuments(ctx, path[1:], k, legacyQuerierCdc)
		case types.QuerySentDocuments:
			return queryGetSentDocuments(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryReceivedReceipts:
			return queryGetReceivedDocsReceipts(ctx, path[1:], k, legacyQuerierCdc)
		case types.QuerySentReceipts:
			return queryGetSentDocsReceipts(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryDocumentReceipts:
			return queryGetDocumentsReceipts(ctx, path[1:], k, legacyQuerierCdc)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("unknown %s query endpoint: %s", types.ModuleName, path[0]))
		}
	}
}

// ----------------------------------
// --- Documents
// ----------------------------------

func queryGetReceivedDocuments(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	addr := path[0]
	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, addr)
	}

	ri := k.UserReceivedDocumentsIterator(ctx, address)
	defer ri.Close()

	documents := []types.Document{}
	for ; ri.Valid(); ri.Next() {
		documentUUID := string(ri.Value())

		document, err := k.GetDocumentByID(ctx, documentUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document with UUID %s even though the user has an associated received document",
					documentUUID,
				),
			)
		}

		documents = append(documents, document)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, documents)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocuments(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	addr := path[0]
	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, addr)
	}

	usdi := k.UserSentDocumentsIterator(ctx, address)
	defer usdi.Close()

	documents := []types.Document{}
	for ; usdi.Valid(); usdi.Next() {
		documentUUID := string(usdi.Value())

		document, err := k.GetDocumentByID(ctx, documentUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document with UUID %s even though the user has an associated received document",
					documentUUID,
				),
			)
		}

		documents = append(documents, document)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, documents)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return bz, nil
}

// ----------------------------------
// --- Documents Receipts
// ----------------------------------

func queryGetReceivedDocsReceipts(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	addr := path[0]
	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, addr)
	}

	ri := k.UserReceivedReceiptsIterator(ctx, address)
	defer ri.Close()

	receipts := []types.DocumentReceipt{}
	for ; ri.Valid(); ri.Next() {
		receiptUUID := string(ri.Value())

		receipt, err := k.GetReceiptByID(ctx, receiptUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document receipt with UUID %s even though the user has an associated received document with it",
					receiptUUID,
				),
			)
		}
		receipts = append(receipts, receipt)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, receipts)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocsReceipts(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	addr := path[0]
	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, addr)
	}

	receipts := []types.DocumentReceipt{}

	ri := k.UserSentReceiptsIterator(ctx, address)
	defer ri.Close()

	for ; ri.Valid(); ri.Next() {
		receiptUUID := string(ri.Value())

		receipt, err := k.GetReceiptByID(ctx, receiptUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document receipt with UUID %s even though the user has an associated received document with it",
					receiptUUID,
				),
			)
		}

		receipts = append(receipts, receipt)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, receipts)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetDocumentsReceipts(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	documentUUID, err := uuid.FromString(path[0])
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("invalid UUID: %s", path[0]))
	}

	receipts := []types.DocumentReceipt{}

	ri := k.UUIDDocumentsReceiptsIterator(ctx, documentUUID.String())
	defer ri.Close()

	for ; ri.Valid(); ri.Next() {
		receiptUUID := string(ri.Value())

		newReceipt, err := k.GetReceiptByID(ctx, receiptUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document receipt with UUID %s even though there is a document associated with it",
					receiptUUID,
				),
			)
		}

		receipts = append(receipts, newReceipt)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, receipts)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return bz, nil
}
