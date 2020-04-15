package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DocumentReceipt contains the generic information about the proof that a shared document (identified by the DocumentUUID field)
// has been read by the recipient.
// It contains information about the sender of the receipt's proof, the recipient,
// the original shared Document transaction's hash, a uuid that identifies which document has been received and
// the proof that document been read by its recipient.
// To be valid, all the fields of document receipt can't be empty and in particular
// transaction hash need to be referred to an existing transaction on chain and the uuid associated with document
// has to be non-empty and unique.
type DocumentReceipt struct {
	UUID         string         `json:"uuid"`
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    sdk.AccAddress `json:"recipient"`
	TxHash       string         `json:"tx_hash"`
	DocumentUUID string         `json:"document_uuid"`
	Proof        string         `json:"proof"` // Optional
}

// Equals implements equatable
func (receipt DocumentReceipt) Equals(rec DocumentReceipt) bool {
	return receipt.UUID == rec.UUID &&
		receipt.Sender.Equals(rec.Sender) &&
		receipt.Recipient.Equals(rec.Recipient) &&
		receipt.TxHash == rec.TxHash &&
		receipt.DocumentUUID == rec.DocumentUUID &&
		receipt.Proof == rec.Proof
}
