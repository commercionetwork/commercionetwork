package types

const (
	ModuleName   = "commerciodocs"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	MsgTypeShareDocument   = "shareDocument"
	MsgTypeDocumentReceipt = "documentReceipt"

	QuerySentDocuments     = "sent"
	QueryReceivedDocuments = "received"
	QueryReceipt           = "receipt"
)
