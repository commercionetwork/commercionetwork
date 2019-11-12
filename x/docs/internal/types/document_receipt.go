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

// ------------------------
// --- DocumentReceipts
// ------------------------

// DocumentReceipts represents a slice of DocumentReceipt
type DocumentReceipts []DocumentReceipt

// IsEmpty returns true if the given slice is empty
func (receipts DocumentReceipts) IsEmpty() bool {
	return len(receipts) == 0
}

// AppendIfMissing returns a new slice containing the given receipt only if it was missing
// from the receipts array, or returns the default array.
// The second value returned tells if the new receipt was properly appended or not.
func (receipts DocumentReceipts) AppendIfMissing(receipt DocumentReceipt) (DocumentReceipts, bool) {
	for _, ele := range receipts {
		if ele.Equals(receipt) {
			return receipts, false
		}
	}
	return append(receipts, receipt), true
}

// AppendAllIfMissing appends each receipt contained inside the other slice into the
// receipts slice, only if they are not already present inside the given slice.
func (receipts DocumentReceipts) AppendAllIfMissing(other DocumentReceipts) DocumentReceipts {
	result := receipts
	for _, receipt := range other {
		result, _ = result.AppendIfMissing(receipt)
	}
	return result
}

// FindByDocumentID returns all the receipts having the given document ID.
func (receipts DocumentReceipts) FindByDocumentID(docID string) DocumentReceipts {
	var foundReceipts DocumentReceipts
	for _, ele := range receipts {
		if ele.DocumentUUID == docID {
			foundReceipts = append(foundReceipts, ele)
		}
	}
	return foundReceipts
}
