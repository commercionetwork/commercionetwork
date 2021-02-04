package keeper

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
)

const (
	eventNewMetadataScheme = "new_metadata_scheme"
	eventNewTSP            = "new_tsp"
	eventSavedDocument     = "new_saved_document"
	eventSavedReceipt      = "new_saved_receipt"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	StoreKey sdk.StoreKey

	GovernmentKeeper government.Keeper

	cdc *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, gKeeper government.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:         storeKey,
		GovernmentKeeper: gKeeper,
		cdc:              cdc,
	}
}

// ----------------------
// --- Metadata schemes
// ----------------------

// AddSupportedMetadataScheme allows to add or update the given metadata scheme definition as a supported metadata
// scheme that will be accepted into document sending transactions
func (keeper Keeper) AddSupportedMetadataScheme(ctx sdk.Context, metadataSchema types.MetadataSchema) {
	store := ctx.KVStore(keeper.StoreKey)

	msk := metadataSchemaKey(metadataSchema)

	store.Set(msk, keeper.cdc.MustMarshalBinaryBare(metadataSchema))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventNewMetadataScheme,
		sdk.NewAttribute("version", metadataSchema.Version),
		sdk.NewAttribute("type", metadataSchema.Type),
		sdk.NewAttribute("uri", metadataSchema.SchemaURI),
	))
}

// IsMetadataSchemeTypeSupported returns true iff the given metadata scheme type is supported
// as an official one
func (keeper Keeper) IsMetadataSchemeTypeSupported(ctx sdk.Context, metadataSchemaType string) bool {
	i := keeper.SupportedMetadataSchemesIterator(ctx)
	defer i.Close()

	for ; i.Valid(); i.Next() {
		var ms types.MetadataSchema
		keeper.cdc.MustUnmarshalBinaryBare(i.Value(), &ms)

		if ms.Type == metadataSchemaType {
			return true
		}
	}

	return false
}

// ------------------------------
// --- Metadata schema proposers
// ------------------------------

// AddTrustedSchemaProposer adds the given proposer to the list of trusted addresses
// that can propose new metadata schemes as officially recognized
func (keeper Keeper) AddTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	store.Set(metadataSchemaProposerKey(proposer), keeper.cdc.MustMarshalBinaryBare(proposer))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventNewTSP,
		sdk.NewAttribute("address", proposer.String()),
	))
}

// IsTrustedSchemaProposer returns true iff the given proposer is a trusted one
func (keeper Keeper) IsTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.StoreKey)

	return store.Has(metadataSchemaProposerKey(proposer))
}

// ----------------------
// --- Documents
// ----------------------

// SaveDocument allows the sharing of a document
func (keeper Keeper) SaveDocument(ctx sdk.Context, document types.Document) error {
	// Check the id validity
	if len(strings.TrimSpace(document.UUID)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid document id: %s", document.UUID))
	}

	// Check any existing document
	// when err is nil, we found a document with said document.UUID
	if _, err := keeper.GetDocumentByID(ctx, document.UUID); err == nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Document with uuid %s already present", document.UUID))
	}

	// Store the document object
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(getDocumentStoreKey(document.UUID), keeper.cdc.MustMarshalBinaryBare(&document))

	// Store the document as sent by the sender

	// Idea: SentDocumentsPrefix + address + document.UUID -> document.UUID
	sentDocumentsStoreKey := getSentDocumentsIdsUUIDStoreKey(document.Sender, document.UUID)

	store.Set(sentDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(document.UUID))

	// Store the documents as received for all the recipients
	for _, recipient := range document.Recipients {
		receivedDocumentsStoreKey := getReceivedDocumentsIdsUUIDStoreKey(recipient, document.UUID)

		store.Set(receivedDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(document.UUID))
	}

	attributes := make([]sdk.Attribute, 0, len(document.Recipients)+2)

	attributes = append(attributes,
		sdk.NewAttribute("sender", document.Sender.String()),
		sdk.NewAttribute("doc_id", document.UUID),
	)

	for i, r := range document.Recipients {
		attributes = append(attributes, sdk.NewAttribute(
			fmt.Sprintf("receiver_%d", i),
			r.String(),
		))
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSavedDocument,
		attributes...,
	))

	return nil
}

// GetDocumentByID returns the document having the given id, or false if no document has been found
func (keeper Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, error) {
	store := ctx.KVStore(keeper.StoreKey)

	documentKey := getDocumentStoreKey(id)
	if !store.Has(documentKey) {
		return types.Document{}, fmt.Errorf("cannot find document with uuid %s", id)
	}

	var document types.Document
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentKey), &document)
	return document, nil
}

// ----------------------
// --- Receipts
// ----------------------

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) error {
	// Check the id
	if len(strings.TrimSpace(receipt.UUID)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid document receipt id: %s", receipt.UUID))
	}

	if _, err := keeper.GetDocumentByID(ctx, receipt.DocumentUUID); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "recepit points to a non-existing document UUID")
	}

	store := ctx.KVStore(keeper.StoreKey)
	sentReceiptsIdsStoreKey := getSentReceiptsIdsUUIDStoreKey(receipt.Sender, receipt.DocumentUUID)
	receivedReceiptIdsStoreKey := getReceivedReceiptsIdsUUIDStoreKey(receipt.Recipient, receipt.DocumentUUID)

	marshaledRecepit := keeper.cdc.MustMarshalBinaryBare(receipt)
	marshaledRecepitID := keeper.cdc.MustMarshalBinaryBare(receipt.UUID)

	// Store the receipt as sent
	if store.Has(sentReceiptsIdsStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("sent receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}
	store.Set(sentReceiptsIdsStoreKey, marshaledRecepitID)

	// Store the receipt as received
	if store.Has(receivedReceiptIdsStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("received receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}
	store.Set(receivedReceiptIdsStoreKey, marshaledRecepitID)

	// Store the receipt
	store.Set(getReceiptStoreKey(receipt.UUID), marshaledRecepit)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSavedReceipt,
		sdk.NewAttribute("receipt_id", receipt.UUID),
		sdk.NewAttribute("document_id", receipt.DocumentUUID),
		sdk.NewAttribute("sender", receipt.Sender.String()),
		sdk.NewAttribute("recipient", receipt.Recipient.String()),
	))

	return nil
}

// GetReceiptByID returns the document receipt having the given id, or false if such receipt could not be found
func (keeper Keeper) GetReceiptByID(ctx sdk.Context, id string) (types.DocumentReceipt, error) {
	store := ctx.KVStore(keeper.StoreKey)
	key := getReceiptStoreKey(id)

	if !store.Has(key) {
		return types.DocumentReceipt{}, fmt.Errorf("cannot find receipt with uuid %s", id)
	}

	var receipt types.DocumentReceipt
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &receipt)
	return receipt, nil
}

// ExtractDocument returns a Document slice instance and its UUID given an iterator byte stream value.
func (keeper Keeper) ExtractDocument(ctx sdk.Context, keyVal []byte) (types.Document, string, error) {
	documentUUID := string(keyVal[len(types.DocumentStorePrefix):])

	document, err := keeper.GetDocumentByID(ctx, documentUUID)
	return document, documentUUID, err
}

// ExtractReceipt returns a DocumentReceipt slice instance and its UUID given an iterator byte stream value.
func (keeper Keeper) ExtractReceipt(ctx sdk.Context, iterVal []byte) (types.DocumentReceipt, string, error) {
	rid := ""
	keeper.cdc.MustUnmarshalBinaryBare(iterVal, &rid)

	newReceipt, err := keeper.GetReceiptByID(ctx, rid)
	return newReceipt, rid, err
}

// ExtractMetadataSchema returns a MetadataSchema slice instance and given an iterator byte stream value.
func (keeper Keeper) ExtractMetadataSchema(iterVal []byte) types.MetadataSchema {
	ms := types.MetadataSchema{}

	keeper.cdc.MustUnmarshalBinaryBare(iterVal, &ms)
	return ms
}

// ExtractTrustedSchemaProposer returns a sdk.AccAddress slice instance given an iterator byte stream value.
func (keeper Keeper) ExtractTrustedSchemaProposer(iterVal []byte) sdk.AccAddress {
	tsp := sdk.AccAddress{}

	keeper.cdc.MustUnmarshalBinaryBare(iterVal, &tsp)
	return tsp
}
