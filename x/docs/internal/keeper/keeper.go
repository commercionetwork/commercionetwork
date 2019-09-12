package keeper

import (
	"strings"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	StoreKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey: storeKey,
		cdc:      cdc,
	}
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
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) {
	store := ctx.KVStore(keeper.StoreKey)

	// Store the document as sent
	var sentDocsList types.Documents
	sentDocs := store.Get(keeper.getSentDocumentsStoreKey(document.Sender))
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)

	sentDocsList = sentDocsList.AppendIfMissing(document)

	store.Set(
		keeper.getSentDocumentsStoreKey(document.Sender),
		keeper.cdc.MustMarshalBinaryBare(&sentDocsList),
	)

	// Store the documents as received
	var recipientDocsList types.Documents
	receivedDocs := store.Get(keeper.getReceivedDocumentsStoreKey(document.Recipient))
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &recipientDocsList)

	recipientDocsList = recipientDocsList.AppendIfMissing(document)

	store.Set(
		keeper.getReceivedDocumentsStoreKey(document.Recipient),
		keeper.cdc.MustMarshalBinaryBare(&recipientDocsList),
	)
}

// GetUserReceivedDocuments returns a list of all the documents that has been received from a user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) types.Documents {

	store := ctx.KVStore(keeper.StoreKey)
	receivedDocs := store.Get([]byte(types.ReceivedDocumentsPrefix + user.String()))

	var receivedDocsList types.Documents
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &receivedDocsList)

	return receivedDocsList
}

// GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) types.Documents {
	store := ctx.KVStore(keeper.StoreKey)
	sentDocs := store.Get([]byte(types.SentDocumentsPrefix + user.String()))

	var sentDocsList types.Documents
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)

	return sentDocsList
}

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

	// Store the receipt as sent
	var sentReceipts types.DocumentReceipts
	sentReceiptBz := store.Get(keeper.getSentReceiptsStoreKey(receipt.Sender))
	keeper.cdc.MustUnmarshalBinaryBare(sentReceiptBz, &sentReceipts)

	sentReceipts = sentReceipts.AppendReceiptIfMissing(receipt)

	store.Set(
		keeper.getSentReceiptsStoreKey(receipt.Sender),
		keeper.cdc.MustMarshalBinaryBare(&sentReceipts),
	)

	// Store the receipt as received
	var receivedReceipts types.DocumentReceipts
	receivedReceiptsBz := store.Get(keeper.getReceivedReceiptsStoreKey(receipt.Recipient))
	keeper.cdc.MustUnmarshalBinaryBare(receivedReceiptsBz, &receivedReceipts)

	receivedReceipts = receivedReceipts.AppendReceiptIfMissing(receipt)

	store.Set(
		keeper.getReceivedReceiptsStoreKey(receipt.Recipient),
		keeper.cdc.MustMarshalBinaryBare(&receivedReceipts),
	)
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
	sentDocs := store.Get([]byte(types.ReceivedDocumentsReceiptsPrefix + user.String()))

	var sentReceipts types.DocumentReceipts
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentReceipts)

	return sentReceipts
}

// --------------------
// --- Genesis utils
// --------------------

// GetUsersSet returns the list of all the users that sent or received at least one document or receipt.
func (keeper Keeper) GetUsersSet(ctx sdk.Context) ([]sdk.AccAddress, error) {
	prefixes := []string{
		types.SentDocumentsPrefix,
		types.ReceivedDocumentsPrefix,
		types.SentDocumentsReceiptsPrefix,
		types.ReceivedDocumentsReceiptsPrefix,
	}

	var err error
	users := ctypes.Addresses{}
	for _, prefix := range prefixes {
		users, err = keeper.addAccountsWithPrefix(ctx, prefix, users)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (keeper Keeper) addAccountsWithPrefix(ctx sdk.Context, prefix string, existingAccounts ctypes.Addresses) (ctypes.Addresses, error) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(prefix))

	for ; iterator.Valid(); iterator.Next() {
		stringKey := strings.ReplaceAll(string(iterator.Key()), prefix, "")
		address, err := sdk.AccAddressFromBech32(stringKey)
		if err != nil {
			return nil, err
		}

		existingAccounts = existingAccounts.AppendIfMissing(address)
	}

	return existingAccounts, nil
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
