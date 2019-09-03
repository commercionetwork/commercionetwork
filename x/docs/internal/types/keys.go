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

	//KVStore prefix
	SentDocumentsPrefix     = StoreKey + "sentBy:"
	ReceivedDocumentsPrefix = StoreKey + "received:"
	DocumentReceiptPrefix   = StoreKey + "receiptOf:"
)
