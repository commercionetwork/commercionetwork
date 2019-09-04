package types

const (
	ModuleName   = "docs"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeShareDocument       = "shareDocument"
	MsgTypeSendDocumentReceipt = "documentReceipt"

	QuerySentDocuments     = "sent"
	QueryReceivedDocuments = "received"
	QueryReceipts          = "receipts"

	SentDocumentsPrefix     = StoreKey + ":documents:sent:"
	ReceivedDocumentsPrefix = StoreKey + ":received:received:"

	SentDocumentsReceiptsPrefix     = StoreKey + ":receipts:sent:"
	ReceivedDocumentsReceiptsPrefix = StoreKey + ":receipts:received:"
)
