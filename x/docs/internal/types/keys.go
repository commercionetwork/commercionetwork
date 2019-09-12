package types

const (
	ModuleName   = "docs"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	SupportedMetadataSchemesStoreKey = StoreKey + "supportedMetadata"
	MetadataSchemaProposersStoreKey  = StoreKey + "metadataSchemaProposers"

	SentDocumentsPrefix     = StoreKey + ":documents:sent:"
	ReceivedDocumentsPrefix = StoreKey + ":received:received:"

	SentDocumentsReceiptsPrefix     = StoreKey + ":receipts:sent:"
	ReceivedDocumentsReceiptsPrefix = StoreKey + ":receipts:received:"

	MsgTypeShareDocument                    = "shareDocument"
	MsgTypeSendDocumentReceipt              = "sendDocumentReceipt"
	MsgTypeAddSupportedMetadataSchema       = "addSupportedMetadataSchema"
	MsgTypeAddTrustedMetadataSchemaProposer = "addTrustedMetadataSchemaProposer"

	QuerySentDocuments            = "sent"
	QueryReceivedDocuments        = "received"
	QueryReceivedReceipts         = "receivedReceipts"
	QuerySentReceipts             = "sentReceipts"
	QuerySupportedMetadataSchemes = "supportedMetadataSchemes"
	QueryTrustedMetadataProposers = "trustedMetadataProposers"
)
