package keeper

import (
	"fmt"
	"strings"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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

// AddSupportedMetadataScheme allows to add the given metadata scheme definition as a supported metadata
// scheme that will be accepted into document sending transactions
func (keeper Keeper) AddSupportedMetadataScheme(ctx sdk.Context, metadataSchema types.MetadataSchema) {
	store := ctx.KVStore(keeper.StoreKey)

	// Read and update
	schemes := keeper.GetSupportedMetadataSchemes(ctx)
	if schemes, success := schemes.AppendIfMissing(metadataSchema); success {
		store.Set([]byte(types.SupportedMetadataSchemesStoreKey), keeper.cdc.MustMarshalBinaryBare(&schemes))
	}
}

// IsMetadataSchemeTypeSupported returns true iff the given metadata scheme type is supported
// as an official one
func (keeper Keeper) IsMetadataSchemeTypeSupported(ctx sdk.Context, metadataSchemaType string) bool {
	schemes := keeper.GetSupportedMetadataSchemes(ctx)
	return schemes.IsTypeSupported(metadataSchemaType)
}

// GetSupportedMetadataSchemes returns the list of all the officially supported metadata schemes
func (keeper Keeper) GetSupportedMetadataSchemes(ctx sdk.Context) types.MetadataSchemes {
	store := ctx.KVStore(keeper.StoreKey)

	var schemes types.MetadataSchemes
	schemesBz := store.Get([]byte(types.SupportedMetadataSchemesStoreKey))
	keeper.cdc.MustUnmarshalBinaryBare(schemesBz, &schemes)

	return schemes
}

// ------------------------------
// --- Metadata schema proposers
// ------------------------------

// AddTrustedSchemaProposer adds the given proposer to the list of trusted addresses
// that can propose new metadata schemes as officially recognized
func (keeper Keeper) AddTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	// Read and update
	proposers := keeper.GetTrustedSchemaProposers(ctx)
	if proposers, success := proposers.AppendIfMissing(proposer); success {
		proposersBz := keeper.cdc.MustMarshalBinaryBare(&proposers)
		store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)
	}
}

// IsTrustedSchemaProposer returns true iff the given proposer is a trusted one
func (keeper Keeper) IsTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) bool {
	return keeper.GetTrustedSchemaProposers(ctx).Contains(proposer)
}

// GetTrustedSchemaProposers returns the list of all the trusted addresses
// that can propose new metadata schemes as officially recognized
func (keeper Keeper) GetTrustedSchemaProposers(ctx sdk.Context) ctypes.Addresses {
	store := ctx.KVStore(keeper.StoreKey)

	var proposers ctypes.Addresses
	proposersBz := store.Get([]byte(types.MetadataSchemaProposersStoreKey))
	keeper.cdc.MustUnmarshalBinaryBare(proposersBz, &proposers)
	return proposers
}

// ----------------------
// --- Documents
// ----------------------

func (keeper Keeper) getDocumentStoreKey(uuid string) []byte {
	return []byte(types.DocumentStorePrefix + uuid)
}

// getSentDocumentsIdsStoreKey returns the byte representation of the key that should be used when updating the
// list of documents that the given user has sent
func (keeper Keeper) getSentDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.SentDocumentsPrefix + user.String())
}

// getReceivedDocumentsIdsStoreKey returns the byte representation of the key that should be used when updating the
// list of documents that the given user has received
func (keeper Keeper) getReceivedDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.ReceivedDocumentsPrefix + user.String())
}

// SaveDocument allows the sharing of a document
func (keeper Keeper) SaveDocument(ctx sdk.Context, document types.Document) sdk.Error {
	// Check the id validity
	if len(strings.TrimSpace(document.UUID)) == 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid document id: %s", document.UUID))
	}

	// Check any existing document
	if _, found := keeper.GetDocumentById(ctx, document.UUID); found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Document with uuid %s already present", document.UUID))
	}

	// Store the document object
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(keeper.getDocumentStoreKey(document.UUID), keeper.cdc.MustMarshalBinaryBare(&document))

	// Store the document as sent by the sender
	sentDocumentsStoreKey := keeper.getSentDocumentsIdsStoreKey(document.Sender)

	var sentDocsList types.DocumentIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(sentDocumentsStoreKey), &sentDocsList)
	if sentDocsList, success := sentDocsList.AppendIfMissing(document.UUID); success {
		store.Set(sentDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(&sentDocsList))
	}

	// Store the documents as received for all the recipients
	for _, recipient := range document.Recipients {
		receivedDocumentsStoreKey := keeper.getReceivedDocumentsIdsStoreKey(recipient)

		var recipientDocsList types.DocumentIDs
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(receivedDocumentsStoreKey), &recipientDocsList)
		if recipientDocsList, success := recipientDocsList.AppendIfMissing(document.UUID); success {
			store.Set(receivedDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(&recipientDocsList))
		}
	}

	return nil
}

// GetDocumentById returns the document having the given id, or false if no document has been found
func (keeper Keeper) GetDocumentById(ctx sdk.Context, id string) (types.Document, bool) {
	store := ctx.KVStore(keeper.StoreKey)

	documentKey := keeper.getDocumentStoreKey(id)
	if !store.Has(documentKey) {
		return types.Document{}, false
	}

	var document types.Document
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentKey), &document)
	return document, true
}

// GetUserReceivedDocuments returns a list of all the documents that has been received from a user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) (types.Documents, sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey)
	receivedDocumentsStoreKey := keeper.getReceivedDocumentsIdsStoreKey(user)

	var receivedDocsIds types.DocumentIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(receivedDocumentsStoreKey), &receivedDocsIds)

	docs := types.Documents{}
	for _, docUuid := range receivedDocsIds {

		// Read the document
		var document types.Document
		documentStoreKey := keeper.getDocumentStoreKey(docUuid)
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentStoreKey), &document)

		// Append it to the list
		docs = docs.AppendIfMissing(document)
	}

	return docs, nil
}

// GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) (types.Documents, sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey)

	var sentDocsIds types.DocumentIDs
	sentDocsIdsBz := store.Get(keeper.getSentDocumentsIdsStoreKey(user))
	keeper.cdc.MustUnmarshalBinaryBare(sentDocsIdsBz, &sentDocsIds)

	docs := types.Documents{}
	for _, docUuid := range sentDocsIds {

		// Read the document
		var document types.Document
		documentStoreKey := keeper.getDocumentStoreKey(docUuid)
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

// ----------------------
// --- Receipts
// ----------------------

// getReceiptStoreKey returns the bytes representation of the key that should be used when
// storing a document receipt
func (keeper Keeper) getReceiptStoreKey(id string) []byte {
	return []byte(types.ReceiptsStorePrefix + id)
}

// getSentReceiptsIdsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts ids that the given user has sent
func (keeper Keeper) getSentReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.SentDocumentsReceiptsPrefix + user.String())
}

// getReceivedReceiptsIdsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts ids that the given user has received
func (keeper Keeper) getReceivedReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.ReceivedDocumentsReceiptsPrefix + user.String())
}

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) sdk.Error {
	// Check the id
	if len(strings.TrimSpace(receipt.UUID)) == 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid document receipt id: %s", receipt.UUID))
	}

	store := ctx.KVStore(keeper.StoreKey)
	sentReceiptsIdsStoreKey := keeper.getSentReceiptsIdsStoreKey(receipt.Sender)
	receivedReceiptIdsStoreKey := keeper.getReceivedReceiptsIdsStoreKey(receipt.Recipient)

	// Store the receipt as sent
	var sentReceiptsIds types.DocumentReceiptsIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(sentReceiptsIdsStoreKey), &sentReceiptsIds)
	if newIds, success := sentReceiptsIds.AppendIfMissing(receipt.UUID); success {
		store.Set(sentReceiptsIdsStoreKey, keeper.cdc.MustMarshalBinaryBare(&newIds))
	}

	// Store the receipt as received
	var receivedReceiptsIds types.DocumentReceiptsIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(receivedReceiptIdsStoreKey), &receivedReceiptsIds)
	if newIds, success := receivedReceiptsIds.AppendIfMissing(receipt.UUID); success {
		store.Set(receivedReceiptIdsStoreKey, keeper.cdc.MustMarshalBinaryBare(&newIds))
	}

	// Store the receipt
	store.Set(keeper.getReceiptStoreKey(receipt.UUID), keeper.cdc.MustMarshalBinaryBare(&receipt))
	return nil
}

// GetReceiptById returns the document receipt having the given id, or false if such receipt could not be found
func (keeper Keeper) GetReceiptById(ctx sdk.Context, id string) (types.DocumentReceipt, bool) {
	store := ctx.KVStore(keeper.StoreKey)
	key := keeper.getReceiptStoreKey(id)

	if !store.Has(key) {
		return types.DocumentReceipt{}, false
	}

	var receipt types.DocumentReceipt
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &receipt)
	return receipt, true
}

// GetUserReceivedReceipts returns the list of all the receipts that the given user has received
func (keeper Keeper) GetUserReceivedReceipts(ctx sdk.Context, user sdk.AccAddress) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	var ids types.DocumentReceiptsIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getReceivedReceiptsIdsStoreKey(user)), &ids)

	receipts := types.DocumentReceipts{}
	for _, id := range ids {
		if receipt, found := keeper.GetReceiptById(ctx, id); found {
			receipts, _ = receipts.AppendIfMissing(receipt)
		}
	}

	return receipts
}

// GetUserReceivedReceiptsForDocument returns the receipts that the given recipient has received for the document having the
// given uuid
func (keeper Keeper) GetUserReceivedReceiptsForDocument(ctx sdk.Context, recipient sdk.AccAddress, docUuid string) types.DocumentReceipts {
	receivedReceipts := keeper.GetUserReceivedReceipts(ctx, recipient)
	return receivedReceipts.FindByDocumentId(docUuid)
}

// GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentReceipts(ctx sdk.Context, user sdk.AccAddress) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	var ids types.DocumentReceiptsIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getSentReceiptsIdsStoreKey(user)), &ids)

	receipts := types.DocumentReceipts{}
	for _, id := range ids {
		if receipt, found := keeper.GetReceiptById(ctx, id); found {
			receipts, _ = receipts.AppendIfMissing(receipt)
		}
	}

	return receipts
}

// GetReceipts returns all the receipts that are stored inside the current context
func (keeper Keeper) GetReceipts(ctx sdk.Context) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	receipts := types.DocumentReceipts{}

	// Iterate over just the sent receipts as the received ones are the same but saved in to different places
	sentReceiptsIterator := sdk.KVStorePrefixIterator(store, []byte(types.SentDocumentsReceiptsPrefix))
	for ; sentReceiptsIterator.Valid(); sentReceiptsIterator.Next() {
		var sentReceipts types.DocumentReceipts
		keeper.cdc.MustUnmarshalBinaryBare(sentReceiptsIterator.Value(), &sentReceipts)
		receipts = receipts.AppendAllIfMissing(sentReceipts)
	}

	receivedReceiptsIterator := sdk.KVStorePrefixIterator(store, []byte(types.ReceivedDocumentsReceiptsPrefix))
	for ; receivedReceiptsIterator.Valid(); receivedReceiptsIterator.Next() {
		var receivedReceipts types.DocumentReceipts
		keeper.cdc.MustUnmarshalBinaryBare(receivedReceiptsIterator.Value(), &receivedReceipts)
		receipts = receipts.AppendAllIfMissing(receivedReceipts)
	}

	return receipts
}
