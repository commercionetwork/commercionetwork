package types

const (
	// ModuleName defines the module name
	ModuleName = "docs"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	/*// Version defines the current version the IBC module supports
	Version = "documents-1"

	// PortID is the default port id that module binds to
	PortID = "documents"*/

	DocumentStorePrefix = StoreKey + ":document:"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("documents-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	DocumentKey      = "Document-value-"
	DocumentCountKey = "Document-count-"

	SentDocumentsPrefix     = StoreKey + ":documents:sent:"
	ReceivedDocumentsPrefix = StoreKey + ":documents:received:"

	ReceiptsStorePrefix             = StoreKey + ":receipts:"
	SentDocumentsReceiptsPrefix     = StoreKey + ":receipts:sent:"
	ReceivedDocumentsReceiptsPrefix = StoreKey + ":receipts:received:"
	DocumentsReceiptsPrefix         = StoreKey + ":receipts:documents:"

	MsgTypeShareDocument       = "shareDocument"
	MsgTypeSendDocumentReceipt = "sendDocumentReceipt"

	QuerySentDocuments     = "sent"
	QueryReceivedDocuments = "received"
	QueryReceivedReceipts  = "receivedReceipts"
	QuerySentReceipts      = "sentReceipts"
)
