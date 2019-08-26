package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/**
Document receipt contains the generic information about the proof that a shared document (identified by the Uuid field)
has been read by the recipient.
It contains information about the sender of the receipt's proof, the recipient,
the original shared Document transaction's hash, a uuid that identifies which document has been received and
the proof that document been read by its recipient.
To be valid, all the fields of document receipt can't be empty and in particular
transaction hash need to be referred to an existing transaction on chain and the uuid associated with document
has to be non-empty and unique.
*/
type DocumentReceipt struct {
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	TxHash    string         `json:"tx_hash"`
	Uuid      string         `json:"uuid"`
	Proof     string         `json:"proof"`
}

func (receipt DocumentReceipt) Equals(rec DocumentReceipt) bool {
	if !receipt.Sender.Equals(rec.Sender) {
		return false
	}
	if !receipt.Recipient.Equals(rec.Recipient) {
		return false
	}
	if receipt.TxHash != rec.TxHash {
		return false
	}
	if receipt.Uuid != rec.Uuid {
		return false
	}
	if receipt.Proof != rec.Proof {
		return false
	}
	return true
}
