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

// getSentDocumentsStoreKey returns the byte representation of the key that should be used when updating the
// list of documents that the given user has sent
func (keeper Keeper) getSentDocumentsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.SentDocumentsPrefix + user.String())
}

// getReceivedDocumentsStoreKey returns the byte representation of the key that should be used when updating the
// list of documents that the given user has received
func (keeper Keeper) getReceivedDocumentsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.ReceivedDocumentsPrefix + user.String())
}

// ShareDocument allows the sharing of a document
func (keeper Keeper) ShareDocument(ctx sdk.Context, sender sdk.AccAddress, recipients []sdk.AccAddress, document types.Document) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)
	sentDocumentsStoreKey := keeper.getSentDocumentsStoreKey(sender)

	// Check any existing document
	if _, found := keeper.GetDocumentById(ctx, document.Uuid); found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("document with uuid %s already present", document.Uuid))
	}

	// Store the document object
	store.Set(keeper.getDocumentStoreKey(document.Uuid), keeper.cdc.MustMarshalBinaryBare(document))

	// Store the document as sent by the sender
	var sentDocsList types.DocumentIds
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(sentDocumentsStoreKey), &sentDocsList)
	if sentDocsList, success := sentDocsList.AppendIfMissing(document.Uuid); success {
		store.Set(sentDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(&sentDocsList))
	}

	// Store the documents as received for all the recipients
	for _, recipient := range recipients {
		receivedDocumentsStoreKey := keeper.getReceivedDocumentsStoreKey(recipient)

		var recipientDocsList types.DocumentIds
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(receivedDocumentsStoreKey), &recipientDocsList)
		if recipientDocsList, success := recipientDocsList.AppendIfMissing(document.Uuid); success {
			store.Set(receivedDocumentsStoreKey, keeper.cdc.MustMarshalBinaryBare(&recipientDocsList))
		}
	}

	return nil
}

func (keeper Keeper) GetDocumentById(ctx sdk.Context, id string) (document types.Document, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	documentKey := keeper.getDocumentStoreKey(document.Uuid)
	if !store.Has(documentKey) {
		return types.Document{}, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(store.Get(documentKey), &document)
	return document, true
}

// GetUserReceivedDocuments returns a list of all the documents that has been received from a user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) (types.Documents, sdk.Error) {
	store := ctx.KVStore(keeper.StoreKey)
	receivedDocumentsStoreKey := keeper.getReceivedDocumentsStoreKey(user)

	var receivedDocsIds types.DocumentIds
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

	var sentDocsIds types.DocumentIds
	sentDocsIdsBz := store.Get(keeper.getSentDocumentsStoreKey(user))
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

// ----------------------
// --- Receipts
// ----------------------

// getSentReceiptsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts that the given user has sent
func (keeper Keeper) getSentReceiptsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.SentDocumentsReceiptsPrefix + user.String())
}

// getReceivedReceiptsStoreKey returns the bytes representation of the key that should be used when
// updating the list of receipts that the given user has received
func (keeper Keeper) getReceivedReceiptsStoreKey(user sdk.Address) []byte {
	return []byte(types.ReceivedDocumentsReceiptsPrefix + user.String())
}

// SendDocumentReceipt allows to properly store the given receipt
func (keeper Keeper) SendDocumentReceipt(ctx sdk.Context, receipt types.DocumentReceipt) {
	store := ctx.KVStore(keeper.StoreKey)
	sentReceiptsStoreKey := keeper.getSentReceiptsStoreKey(receipt.Sender)
	receivedReceiptStoreKey := keeper.getReceivedReceiptsStoreKey(receipt.Recipient)

	// Store the receipt as sent
	var sentReceipts types.DocumentReceipts
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(sentReceiptsStoreKey), &sentReceipts)
	if sentReceipts, success := sentReceipts.AppendIfMissing(receipt); success {
		store.Set(sentReceiptsStoreKey, keeper.cdc.MustMarshalBinaryBare(&sentReceipts))
	}

	// Store the receipt as received
	var receivedReceipts types.DocumentReceipts
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(receivedReceiptStoreKey), &receivedReceipts)
	if receivedReceipts, success := receivedReceipts.AppendIfMissing(receipt); success {
		store.Set(receivedReceiptStoreKey, keeper.cdc.MustMarshalBinaryBare(&receivedReceipts))
	}
}

// GetUserReceivedReceipts returns the list of all the receipts that the given user has received
func (keeper Keeper) GetUserReceivedReceipts(ctx sdk.Context, user sdk.AccAddress) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)

	var receivedReceipts types.DocumentReceipts
	receiptsBz := store.Get(keeper.getReceivedReceiptsStoreKey(user))
	keeper.cdc.MustUnmarshalBinaryBare(receiptsBz, &receivedReceipts)

	return receivedReceipts
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
	sentDocs := store.Get(keeper.getSentReceiptsStoreKey(user))

	var sentReceipts types.DocumentReceipts
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentReceipts)

	return sentReceipts
}

// ----------------------
// --- Genesis utils
// ----------------------

// GetUsersSet returns the list of all the users that sent or received at least one document or receipt.
func (keeper Keeper) GetUsersSet(ctx sdk.Context) ctypes.Addresses {
	store := ctx.KVStore(keeper.StoreKey)

	prefixes := []string{
		types.SentDocumentsPrefix,
		types.ReceivedDocumentsPrefix,
		types.SentDocumentsReceiptsPrefix,
		types.ReceivedDocumentsReceiptsPrefix,
	}

	users := ctypes.Addresses{}
	for _, prefix := range prefixes {
		iterator := sdk.KVStorePrefixIterator(store, []byte(prefix))
		for ; iterator.Valid(); iterator.Next() {
			addressStr := strings.ReplaceAll(string(iterator.Key()), prefix, "")
			address, _ := sdk.AccAddressFromBech32(addressStr)
			users, _ = users.AppendIfMissing(address)
		}
	}

	return users
}

// SetUserDocuments should be used while initializing the genesis and allows to bulk update
// all the sent and received documents related to the given user
func (keeper Keeper) SetUserDocuments(ctx sdk.Context, user sdk.AccAddress, sentDocuments, receivedDocuments types.Documents) {
	store := ctx.KVStore(keeper.StoreKey)

	sentDocsBz := keeper.cdc.MustMarshalBinaryBare(&sentDocuments)
	if sentDocsBz != nil {
		store.Set(keeper.getSentDocumentsStoreKey(user), sentDocsBz)
	}

	receivedDocsBz := keeper.cdc.MustMarshalBinaryBare(&receivedDocuments)
	if receivedDocsBz != nil {
		store.Set(keeper.getReceivedDocumentsStoreKey(user), receivedDocsBz)
	}
}

// SetUserDocuments should be used while initializing the genesis and allows to bulk update
// all the sent and received receipts related to the given user
func (keeper Keeper) SetUserReceipts(ctx sdk.Context, user sdk.AccAddress, sentReceipts, receivedReceipts types.DocumentReceipts) {
	store := ctx.KVStore(keeper.StoreKey)

	sentReceiptsBz := keeper.cdc.MustMarshalBinaryBare(&sentReceipts)
	if sentReceiptsBz != nil {
		store.Set(keeper.getSentReceiptsStoreKey(user), sentReceiptsBz)
	}

	receivedReceiptsBz := keeper.cdc.MustMarshalBinaryBare(&receivedReceipts)
	if receivedReceiptsBz != nil {
		store.Set(keeper.getReceivedReceiptsStoreKey(user), receivedReceiptsBz)
	}
}
