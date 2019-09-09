package types

const (
	ModuleName   = "docs"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeShareDocument              = "shareDocument"
	MsgTypeSendDocumentReceipt        = "sendDocumentReceipt"
	MsgTypeAddSupportedMetadataSchema = "addSupportedMetadataSchema"

	QuerySentDocuments     = "sent"
	QueryReceivedDocuments = "received"
	QueryReceipts          = "receipts"

	SupportedMetadataStoreKey = StoreKey + "supportedMetadata"

	SentDocumentsPrefix     = StoreKey + ":documents:sent:"
	ReceivedDocumentsPrefix = StoreKey + ":received:received:"

	SentDocumentsReceiptsPrefix     = StoreKey + ":receipts:sent:"
	ReceivedDocumentsReceiptsPrefix = StoreKey + ":receipts:received:"
)
