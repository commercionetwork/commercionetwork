package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryReceivedDocuments:
			return queryGetReceivedDocuments(ctx, path[1:], keeper)
		case types.QuerySentDocuments:
			return queryGetSentDocuments(ctx, path[1:], keeper)
		case types.QueryReceivedReceipts:
			return queryGetReceivedDocsReceipts(ctx, path[1:], keeper)
		case types.QuerySentReceipts:
			return queryGetSentDocsReceipts(ctx, path[1:], keeper)
		case types.QuerySupportedMetadataSchemes:
			return querySupportedMetadataSchemes(ctx, path[1:], keeper)
		case types.QueryTrustedMetadataProposers:
			return queryTrustedMetadataProposers(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

// ----------------------------------
// --- Documents
// ----------------------------------

func queryGetReceivedDocuments(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	ri := keeper.UserReceivedDocumentsIterator(ctx, address)
	defer ri.Close()

	receivedResult := []types.Document{}
	for ; ri.Valid(); ri.Next() {
		documentUUID := ""
		keeper.cdc.MustUnmarshalBinaryBare(ri.Value(), &documentUUID)

		document, err := keeper.GetDocumentByID(ctx, documentUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document with UUID %s even though the user has an associated received document",
					documentUUID,
				),
			)
		}

		receivedResult = append(receivedResult, document)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocuments(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	usdi := keeper.UserSentDocumentsIterator(ctx, address)
	defer usdi.Close()

	receivedResult := []types.Document{}
	for ; usdi.Valid(); usdi.Next() {
		documentUUID := ""
		keeper.cdc.MustUnmarshalBinaryBare(usdi.Value(), &documentUUID)

		document, err := keeper.GetDocumentByID(ctx, documentUUID)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document with UUID %s even though the user has an associated received document",
					documentUUID,
				),
			)
		}

		receivedResult = append(receivedResult, document)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// ----------------------------------
// --- Documents receipts
// ----------------------------------

func queryGetReceivedDocsReceipts(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	var uuid string
	if len(path) == 2 {
		uuid = path[1]
	}

	receipts := []types.DocumentReceipt{}

	ri := keeper.UserReceivedReceiptsIterator(ctx, address)
	defer ri.Close()

	for ; ri.Valid(); ri.Next() {
		rid := ""
		keeper.cdc.MustUnmarshalBinaryBare(ri.Value(), &rid)

		newReceipt, err := keeper.GetReceiptByID(ctx, rid)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document receipt with UUID %s even though the user has an associated received document with it",
					rid,
				),
			)
		}

		if uuid == "" {
			receipts = append(receipts, newReceipt)
			continue
		}

		if newReceipt.DocumentUUID == uuid {
			receipts = append(receipts, newReceipt)
		}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &receipts)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocsReceipts(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := path[0]
	address, err := sdk.AccAddressFromBech32(addr)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, addr)
	}

	receipts := []types.DocumentReceipt{}

	ri := keeper.UserSentReceiptsIterator(ctx, address)
	defer ri.Close()

	for ; ri.Valid(); ri.Next() {
		rid := ""
		keeper.cdc.MustUnmarshalBinaryBare(ri.Value(), &rid)

		newReceipt, err := keeper.GetReceiptByID(ctx, rid)
		if err != nil {
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf(
					"could not find document receipt with UUID %s even though the user has an associated received document with it",
					rid,
				),
			)
		}

		receipts = append(receipts, newReceipt)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &receipts)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// ----------------------------------
// --- Document metadata schemes
// ----------------------------------

func querySupportedMetadataSchemes(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	si := keeper.SupportedMetadataSchemesIterator(ctx)
	defer si.Close()

	schemes := []types.MetadataSchema{}
	for ; si.Valid(); si.Next() {
		var ms types.MetadataSchema
		keeper.cdc.MustUnmarshalBinaryBare(si.Value(), &ms)
		schemes = append(schemes, ms)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &schemes)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// -----------------------------------------
// --- Document metadata schemes proposers
// -----------------------------------------

func queryTrustedMetadataProposers(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	pi := keeper.TrustedSchemaProposersIterator(ctx)
	defer pi.Close()

	proposers := []sdk.AccAddress{}
	for ; pi.Valid(); pi.Next() {
		aa := sdk.AccAddress{}
		keeper.cdc.MustUnmarshalBinaryBare(pi.Value(), &aa)
		proposers = append(proposers, aa)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &proposers)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}
