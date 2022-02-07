package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// DocumentReceipt contains the generic information about the proof that a shared document (identified by the DocumentUUID field)
// has been read by the recipient.
// It contains information about the sender of the receipt's proof, the recipient,
// the original shared Document transaction's hash, a uuid that identifies which document has been received and
// the proof that the document been read by its recipient.
// To be valid, all the fields of document receipt can't be empty and in particular
// transaction hash needs to be referred to an existing transaction on chain and the uuid associated with document
// has to be non-empty and unique.

// Equals implements equatable
func (receipt DocumentReceipt) Equals(rec DocumentReceipt) bool {
	return receipt.UUID == rec.UUID &&
		receipt.Sender == rec.Sender &&
		receipt.Recipient == rec.Recipient &&
		receipt.TxHash == rec.TxHash &&
		receipt.DocumentUUID == rec.DocumentUUID &&
		receipt.Proof == rec.Proof
}

// TODO: add similar method to package legacy for v220
func (receipt DocumentReceipt) Validate() error {
	if receipt.UUID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "UUID cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(receipt.Sender); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, receipt.Sender)
	}

	if _, err := sdk.AccAddressFromBech32(receipt.Recipient); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, receipt.Recipient)
	}

	if receipt.TxHash == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "TxHash cannot be empty")
	}

	if receipt.DocumentUUID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "document UUID cannot be empty")
	}

	if receipt.Proof == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Proof cannot be empty")
	}

	return nil
}
