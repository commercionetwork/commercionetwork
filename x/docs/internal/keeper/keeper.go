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

func (keeper Keeper) metadataSchemaKey(ms types.MetadataSchema) []byte {
	return append([]byte(types.SupportedMetadataSchemesStoreKey), []byte(ms.SchemaURI)...)
}

func (keeper Keeper) metadataSchemaProposerKey(addr sdk.AccAddress) []byte {
	return append([]byte(types.MetadataSchemaProposersStoreKey), keeper.cdc.MustMarshalBinaryBare(addr)...)
}

func (keeper Keeper) SupportedMetadataSchemesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SupportedMetadataSchemesStoreKey))
}

// AddSupportedMetadataScheme allows to add or update the given metadata scheme definition as a supported metadata
// scheme that will be accepted into document sending transactions
func (keeper Keeper) AddSupportedMetadataScheme(ctx sdk.Context, metadataSchema types.MetadataSchema) {
	store := ctx.KVStore(keeper.StoreKey)

	msk := keeper.metadataSchemaKey(metadataSchema)

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

// GetSupportedMetadataSchemes returns the list of all the officially supported metadata schemes
// TODO: remove
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
// TODO: check with alessio
func (keeper Keeper) AddTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	store.Set(keeper.metadataSchemaProposerKey(proposer), keeper.cdc.MustMarshalBinaryBare(proposer))
}

// IsTrustedSchemaProposer returns true iff the given proposer is a trusted one
func (keeper Keeper) IsTrustedSchemaProposer(ctx sdk.Context, proposer sdk.AccAddress) bool {
	store := ctx.KVStore(keeper.StoreKey)

	return store.Has(keeper.metadataSchemaProposerKey(proposer))
}

// GetTrustedSchemaProposers returns the list of all the trusted addresses
// that can propose new metadata schemes as officially recognized
// TODO: remove
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
func (keeper Keeper) getSentDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(keeper.cdc.MustMarshalBinaryBare(user), []byte(":"+documentUUID)...)

	return append([]byte(types.SentDocumentsPrefix), userPart...)
}

// getReceivedDocumentsIdsStoreKey returns the byte representation of the key that should be used when updating the
// list of documents that the given user has received
func (keeper Keeper) getReceivedDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(keeper.cdc.MustMarshalBinaryBare(user), []byte(":"+documentUUID)...)

	return append([]byte(types.ReceivedDocumentsPrefix), userPart...)
}

func (keeper Keeper) getReceivedDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsPrefix), keeper.cdc.MustMarshalBinaryBare(user)...)
}

func (keeper Keeper) getSentDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsPrefix), keeper.cdc.MustMarshalBinaryBare(user)...)
}

// SaveDocument allows the sharing of a document
func (keeper Keeper) SaveDocument(ctx sdk.Context, document types.Document) sdk.Error {
	// Check the id validity
	if len(strings.TrimSpace(document.UUID)) == 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid document id: %s", document.UUID))
	}

	// Check any existing document
	if _, found := keeper.GetDocumentByID(ctx, document.UUID); found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Document with uuid %s already present", document.UUID))
	}

	// Store the document object
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(keeper.getDocumentStoreKey(document.UUID), keeper.cdc.MustMarshalBinaryBare(&document))

	// Store the document as sent by the sender

	// Idea: SentDocumentsPrefix + address + document.UUID -> document.UUID
	sentDocumentsStoreKey := keeper.getSentDocumentsIdsUUIDStoreKey(document.Sender, document.UUID)

	store.Set(sentDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(document.UUID))

	// Store the documents as received for all the recipients
	for _, recipient := range document.Recipients {
		receivedDocumentsStoreKey := keeper.getReceivedDocumentsIdsUUIDStoreKey(recipient, document.UUID)

		store.Set(receivedDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(document.UUID))
	}

	return nil
}

// GetDocumentByID returns the document having the given id, or false if no document has been found
func (keeper Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, bool) {
	store := ctx.KVStore(keeper.StoreKey)

	documentKey := keeper.getDocumentStoreKey(id)
	if !store.Has(documentKey) {
		return types.Document{}, false
	}

	var document types.Document
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentKey), &document)
	return document, true
}

func (keeper Keeper) UserReceivedDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, keeper.getReceivedDocumentsIdsStoreKey(user))
}

// GetUserReceivedDocuments returns a list of all the documents that has been received from a user
// TODO: remove
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) (types.Documents, sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey)
	receivedDocumentsStoreKey := keeper.getReceivedDocumentsIdsStoreKey(user)

	var receivedDocsIds types.DocumentIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(receivedDocumentsStoreKey), &receivedDocsIds)

	docs := types.Documents{}
	for _, docUUID := range receivedDocsIds {

		// Read the document
		var document types.Document
		documentStoreKey := keeper.getDocumentStoreKey(docUUID)
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentStoreKey), &document)

		// Append it to the list
		docs = docs.AppendIfMissing(document)
	}

	return docs, nil
}

func (keeper Keeper) UserSentDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, keeper.getSentDocumentsIdsStoreKey(user))
}

// GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) (types.Documents, sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey)

	var sentDocsIds types.DocumentIDs
	sentDocsIdsBz := store.Get(keeper.getSentDocumentsIdsStoreKey(user))
	keeper.cdc.MustUnmarshalBinaryBare(sentDocsIdsBz, &sentDocsIds)

	docs := types.Documents{}
	for _, docUUID := range sentDocsIds {

		// Read the document
		var document types.Document
		documentStoreKey := keeper.getDocumentStoreKey(docUUID)
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

// getReceiptStoreKey returns the bytes representation of the key that should be used when
// storing a document receipt
func (keeper Keeper) getReceiptStoreKey(id string) []byte {
	return []byte(types.ReceiptsStorePrefix + id)
}

// getSentReceiptsIdsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts ids that the given user has sent
func (keeper Keeper) getSentReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsReceiptsPrefix), keeper.cdc.MustMarshalBinaryBare(user)...)
}

// getReceivedReceiptsIdsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts ids that the given user has received
func (keeper Keeper) getReceivedReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), keeper.cdc.MustMarshalBinaryBare(user)...)
}

func (keeper Keeper) getSentReceiptsIdsUUIDStoreKey(user sdk.AccAddress, recepitUUID string) []byte {
	recepitPart := append(keeper.cdc.MustMarshalBinaryBare(user), []byte(":"+recepitUUID)...)

	return append([]byte(types.SentDocumentsReceiptsPrefix), recepitPart...)
}

// getReceivedReceiptsIdsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts ids that the given user has received
func (keeper Keeper) getReceivedReceiptsIdsUUIDStoreKey(user sdk.AccAddress, recepitUUID string) []byte {
	recepitPart := append(keeper.cdc.MustMarshalBinaryBare(user), []byte(":"+recepitUUID)...)

	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), recepitPart...)
}

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) sdk.Error {
	// Check the id
	if len(strings.TrimSpace(receipt.UUID)) == 0 {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid document receipt id: %s", receipt.UUID))
	}

	store := ctx.KVStore(keeper.StoreKey)
	sentReceiptsIdsStoreKey := keeper.getSentReceiptsIdsUUIDStoreKey(receipt.Sender, receipt.UUID)
	receivedReceiptIdsStoreKey := keeper.getReceivedReceiptsIdsUUIDStoreKey(receipt.Recipient, receipt.UUID)

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
	store.Set(keeper.getReceiptStoreKey(receipt.UUID), marshaledRecepit)
	return nil
}

// GetReceiptByID returns the document receipt having the given id, or false if such receipt could not be found
func (keeper Keeper) GetReceiptByID(ctx sdk.Context, id string) (types.DocumentReceipt, bool) {
	store := ctx.KVStore(keeper.StoreKey)
	key := keeper.getReceiptStoreKey(id)

	if !store.Has(key) {
		return types.DocumentReceipt{}, false
	}

	var receipt types.DocumentReceipt
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &receipt)
	return receipt, true
}

func (keeper Keeper) UserReceivedReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, keeper.getReceivedReceiptsIdsStoreKey(user))
}

// GetUserReceivedReceipts returns the list of all the receipts that the given user has received
func (keeper Keeper) GetUserReceivedReceipts(ctx sdk.Context, user sdk.AccAddress) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	var ids types.DocumentReceiptsIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getReceivedReceiptsIdsStoreKey(user)), &ids)

	receipts := types.DocumentReceipts{}
	for _, id := range ids {
		if receipt, found := keeper.GetReceiptByID(ctx, id); found {
			receipts, _ = receipts.AppendIfMissing(receipt)
		}
	}

	return receipts
}

// GetUserReceivedReceiptsForDocument returns the receipts that the given recipient has received for the document having the
// given uuid
// TODO: find a way to rework this
func (keeper Keeper) GetUserReceivedReceiptsForDocument(ctx sdk.Context, recipient sdk.AccAddress, docUUID string) types.DocumentReceipts {
	receivedReceipts := keeper.GetUserReceivedReceipts(ctx, recipient)
	return receivedReceipts.FindByDocumentID(docUUID)
}

// GetUserSentDocuments returns a list of all documents sent by user
// TODO: remove
func (keeper Keeper) GetUserSentReceipts(ctx sdk.Context, user sdk.AccAddress) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	var ids types.DocumentReceiptsIDs
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getSentReceiptsIdsStoreKey(user)), &ids)

	receipts := types.DocumentReceipts{}
	for _, id := range ids {
		if receipt, found := keeper.GetReceiptByID(ctx, id); found {
			receipts, _ = receipts.AppendIfMissing(receipt)
		}
	}

	return receipts
}

func (keeper Keeper) UserSentReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, keeper.getSentReceiptsIdsStoreKey(user))
}

// GetReceipts returns all the receipts that are stored inside the current context
// TODO: remove
func (keeper Keeper) GetReceipts(ctx sdk.Context) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	receipts := types.DocumentReceipts{}

	// Iterate over just the sent receipts as the received ones are the same but saved in to different places
	sentReceiptsIterator := sdk.KVStorePrefixIterator(store, []byte(types.SentDocumentsReceiptsPrefix))
	defer sentReceiptsIterator.Close()
	for ; sentReceiptsIterator.Valid(); sentReceiptsIterator.Next() {
		var sentReceipts types.DocumentReceipts
		keeper.cdc.MustUnmarshalBinaryBare(sentReceiptsIterator.Value(), &sentReceipts)
		receipts = receipts.AppendAllIfMissing(sentReceipts)
	}

	receivedReceiptsIterator := sdk.KVStorePrefixIterator(store, []byte(types.ReceivedDocumentsReceiptsPrefix))
	defer receivedReceiptsIterator.Close()
	for ; receivedReceiptsIterator.Valid(); receivedReceiptsIterator.Next() {
		var receivedReceipts types.DocumentReceipts
		keeper.cdc.MustUnmarshalBinaryBare(receivedReceiptsIterator.Value(), &receivedReceipts)
		receipts = receipts.AppendAllIfMissing(receivedReceipts)
	}

	return receipts
}

func (keeper Keeper) ReceiptsIterators(ctx sdk.Context) (sdk.Iterator, sdk.Iterator) {
	store := ctx.KVStore(keeper.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SentDocumentsReceiptsPrefix)),
		sdk.KVStorePrefixIterator(store, []byte(types.ReceivedDocumentsReceiptsPrefix))
}
