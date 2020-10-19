package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/docs/types"
)

// metadataSchemaKey generates an unique store key for a MetadataSchema
func metadataSchemaKey(ms types.MetadataSchema) []byte {
	return append([]byte(types.SupportedMetadataSchemesStoreKey), []byte(ms.SchemaURI)...)
}

// metadataSchemaProposerKey generates an unique store key for a Schema Proposer
func metadataSchemaProposerKey(addr sdk.AccAddress) []byte {
	return append([]byte(types.MetadataSchemaProposersStoreKey), addr...)
}

// getDocumentStoreKey generates an unique store key for a Document UUID
func getDocumentStoreKey(uuid string) []byte {
	return []byte(types.DocumentStorePrefix + uuid)
}

// getSentDocumentsIdsUUIDStoreKey generates a SentDocumentID for a given user and document UUID
func getSentDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(user, []byte(":"+documentUUID)...)

	return append([]byte(types.SentDocumentsPrefix), userPart...)
}

// getReceivedDocumentsIdsUUIDStoreKey generates a ReceivedDocumentID for a given user and document UUID
func getReceivedDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(user, []byte(":"+documentUUID)...)

	return append([]byte(types.ReceivedDocumentsPrefix), userPart...)
}

// getReceivedDocumentsIdsStoreKey generates a ReceivedDocumentsID store key for a given user
func getReceivedDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsPrefix), user...)
}

// getSentDocumentsIdsStoreKey generates a ReceivedDocumentsID store key for a given user
func getSentDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsPrefix), user...)
}

// getReceiptStoreKey generates a store key for a document UUID
func getReceiptStoreKey(id string) []byte {
	return []byte(types.ReceiptsStorePrefix + id)
}

// getSentReceiptsIdsStoreKey generates a SentReceiptsID store key for a given user
func getSentReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsReceiptsPrefix), user...)
}

// getReceivedReceiptsIdsStoreKey generates a ReceivedReceiptsID store key for a given user
func getReceivedReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), user...)
}

// getSentReceiptsIdsUUIDStoreKey generates a SentReceiptsID store key for a given user
func getSentReceiptsIdsUUIDStoreKey(user sdk.AccAddress, recepitUUID string) []byte {
	recepitPart := append(user, []byte(":"+recepitUUID)...)

	return append([]byte(types.SentDocumentsReceiptsPrefix), recepitPart...)
}

// getReceivedReceiptsIdsUUIDStoreKey generates a ReceivedReceiptsID store key for a given user
func getReceivedReceiptsIdsUUIDStoreKey(user sdk.AccAddress, recepitUUID string) []byte {
	recepitPart := append(user, []byte(":"+recepitUUID)...)

	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), recepitPart...)
}
