package keeper

import (
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
}

// IsTrustedSchemaProposer returns true iff the given proposer is a trusted one
func (keeper Keeper) IsTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.StoreKey)

	return store.Has(metadataSchemaProposerKey(proposer))
}

func (keeper Keeper) SupportedMetadataSchemesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SupportedMetadataSchemesStoreKey))
}

func (keeper Keeper) TrustedSchemaProposersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.MetadataSchemaProposersStoreKey))
}

// ----------------------
// --- Documents
// ----------------------

// SaveDocument allows the sharing of a document
func (keeper Keeper) SaveDocument(ctx sdk.Context, document types.Document) sdk.Error {
	// Check the id validity
	if len(strings.TrimSpace(document.UUID)) == 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid document id: %s", document.UUID))
	}

	// Check any existing document
	if _, err := keeper.GetDocumentByID(ctx, document.UUID); err != nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Document with uuid %s already present", document.UUID))
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

func (keeper Keeper) UserReceivedDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getReceivedDocumentsIdsStoreKey(user))
}

func (keeper Keeper) UserSentDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getSentDocumentsIdsStoreKey(user))
}

// GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) (types.Documents, sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey)

	var sentDocsIds types.DocumentIDs
	sentDocsIdsBz := store.Get(getSentDocumentsIdsStoreKey(user))
	keeper.cdc.MustUnmarshalBinaryBare(sentDocsIdsBz, &sentDocsIds)

	docs := types.Documents{}
	for _, docUUID := range sentDocsIds {

		// Read the document
		var document types.Document
		documentStoreKey := getDocumentStoreKey(docUUID)
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentStoreKey), &document)

		// Append it to the list
		docs = docs.AppendIfMissing(document)
	}

	return docs, nil
}

// GetDocuments returns all the documents stored inside the given context
func (keeper Keeper) GetDocuments(ctx sdk.Context) types.Documents {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DocumentStorePrefix))

	documents := types.Documents{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var document types.Document
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &document)
		documents = documents.AppendIfMissing(document)
	}

	return documents
}

func (keeper Keeper) DocumentsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.DocumentStorePrefix))
}

// ----------------------
// --- Receipts
// ----------------------

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) sdk.Error {
	// Check the id
	if len(strings.TrimSpace(receipt.UUID)) == 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid document receipt id: %s", receipt.UUID))
	}

	store := ctx.KVStore(keeper.StoreKey)
	sentReceiptsIdsStoreKey := getSentReceiptsIdsUUIDStoreKey(receipt.Sender, receipt.UUID)
	receivedReceiptIdsStoreKey := getReceivedReceiptsIdsUUIDStoreKey(receipt.Recipient, receipt.UUID)

	marshaledRecepit := keeper.cdc.MustMarshalBinaryBare(receipt)

	// Store the receipt as sent
	if store.Has(sentReceiptsIdsStoreKey) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Receipt's UUID already present: %s", receipt.UUID))
	}
	store.Set(sentReceiptsIdsStoreKey, marshaledRecepit)

	// Store the receipt as received
	if store.Has(receivedReceiptIdsStoreKey) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Receipt's UUID already present: %s", receipt.UUID))
	}
	store.Set(receivedReceiptIdsStoreKey, marshaledRecepit)

	// Store the receipt
	store.Set(getReceiptStoreKey(receipt.UUID), marshaledRecepit)
	return nil
}

// GetReceiptByID returns the document receipt having the given id, or false if such receipt could not be found
func (keeper Keeper) GetReceiptByID(ctx sdk.Context, id string) (types.DocumentReceipt, bool) {
	store := ctx.KVStore(keeper.StoreKey)
	key := getReceiptStoreKey(id)

	if !store.Has(key) {
		return types.DocumentReceipt{}, false
	}

	var receipt types.DocumentReceipt
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &receipt)
	return receipt, true
}

func (keeper Keeper) UserReceivedReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getReceivedReceiptsIdsStoreKey(user))
}

func (keeper Keeper) UserSentReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, getSentReceiptsIdsStoreKey(user))
}

func (keeper Keeper) ReceiptsIterators(ctx sdk.Context) (sdk.Iterator, sdk.Iterator) {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SentDocumentsReceiptsPrefix)),
		sdk.KVStorePrefixIterator(store, []byte(types.ReceivedDocumentsReceiptsPrefix))
}
