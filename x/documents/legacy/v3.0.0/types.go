package v3_0_0

import (
	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

const (
	ModuleName = types.ModuleName

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

const (
	DocumentStorePrefix = StoreKey + ":document:"

	DocumentPrefix          = StoreKey + ":documents:"
	SentDocumentsPrefix     = DocumentPrefix + "sent:"
	ReceivedDocumentsPrefix = DocumentPrefix + "received:"

	ReceiptsStorePrefix = StoreKey + ":receipt:"

	ReceiptsPrefix                  = StoreKey + ":receipts:"
	SentDocumentsReceiptsPrefix     = ReceiptsPrefix + "sent:"
	ReceivedDocumentsReceiptsPrefix = ReceiptsPrefix + "received:"
	DocumentsReceiptsPrefix         = ReceiptsPrefix + "documents:"

	MsgTypeShareDocument       = "shareDocument"
	MsgTypeSendDocumentReceipt = "sendDocumentReceipt"
)
