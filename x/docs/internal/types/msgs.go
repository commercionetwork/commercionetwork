package types

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var algorithms = map[string]int{
	"md5":     32,
	"sha-1":   40,
	"sha-224": 56,
	"sha-256": 64,
	"sha-384": 96,
	"sha-512": 128,
}

// ----------------------------------
// --- MsgShareDocument
// ----------------------------------

type MsgShareDocument struct {
	Sender     sdk.AccAddress  `json:"sender"`
	Recipients types.Addresses `json:"recipients"`
	Document   Document        `json:"document"`
}

func NewMsgShareDocument(sender sdk.AccAddress, recipients types.Addresses, document Document) MsgShareDocument {
	return MsgShareDocument{
		Sender:     sender,
		Recipients: recipients,
		Document:   document,
	}
}

// RouterKey Implements Msg.
func (msg MsgShareDocument) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgShareDocument) Type() string { return MsgTypeShareDocument }

// ValidateBasic Implements Msg.
// TODO: Test more
func (msg MsgShareDocument) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.Recipients.Empty() {
		return sdk.ErrUnknownRequest("Recipients cannot be empty")
	}
	for _, recipient := range msg.Recipients {
		if recipient.Empty() {
			return sdk.ErrInvalidAddress(recipient.String())
		}
	}

	err := msg.Document.Validate()
	if err != nil {
		return nil
	}

	if msg.Document.EncryptionData != nil {

		for _, recipient := range msg.Recipients {
			found := false

			// Check that each address inside the EncryptionData object is contained inside the list of addresses
			for _, encAdd := range msg.Document.EncryptionData.Keys {

				// Check that each recipient has an encrypted data associated to it
				if recipient.Equals(encAdd.Recipient) {
					found = true
				}

				if !msg.Recipients.Contains(encAdd.Recipient) {
					errMsg := fmt.Sprintf(
						"%s is a recipient inside encryption data but not inside the message",
						encAdd.Recipient.String(),
					)
					return sdk.ErrInvalidAddress(errMsg)
				}
			}

			if !found {
				// The recipient is not found inside the list of encrypted data recipients
				errMsg := fmt.Sprintf("%s is a recipient but has no encryption data specified", recipient.String())
				return sdk.ErrInvalidAddress(errMsg)
			}
		}

	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgShareDocument) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgShareDocument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ----------------------------------
// --- MsgSendDocumentReceipt
// ----------------------------------

type MsgSendDocumentReceipt DocumentReceipt

func NewMsgDocumentReceipt(receipt DocumentReceipt) MsgSendDocumentReceipt {
	return MsgSendDocumentReceipt(receipt)
}

// RouterKey Implements Msg.
func (msg MsgSendDocumentReceipt) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgSendDocumentReceipt) Type() string { return MsgTypeSendDocumentReceipt }

// ValidateBasic Implements Msg.
func (msg MsgSendDocumentReceipt) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}
	if len(msg.TxHash) == 0 {
		return sdk.ErrUnknownRequest("Send Document's Transaction Hash can't be empty")
	}
	if !validateUuid(msg.DocumentUuid) {
		return sdk.ErrUnknownRequest("Invalid document UUID")
	}
	if len(msg.Proof) == 0 {
		return sdk.ErrUnknownRequest("Receipt proof can't be empty")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSendDocumentReceipt) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSendDocumentReceipt) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ------------------------------------
// --- MsgAddSupportedMetadataSchema
// ------------------------------------

type MsgAddSupportedMetadataSchema struct {
	Signer sdk.AccAddress `json:"signer"`
	Schema MetadataSchema `json:"schema"`
}

// RouterKey Implements Msg.
func (msg MsgAddSupportedMetadataSchema) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgAddSupportedMetadataSchema) Type() string { return MsgTypeAddSupportedMetadataSchema }

// ValidateBasic Implements Msg.
func (msg MsgAddSupportedMetadataSchema) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if err := msg.Schema.Validate(); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddSupportedMetadataSchema) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddSupportedMetadataSchema) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// -----------------------------------------
// --- MsgAddTrustedMetadataSchemaProposer
// -----------------------------------------

type MsgAddTrustedMetadataSchemaProposer struct {
	Proposer sdk.AccAddress `json:"proposer"`
	Signer   sdk.AccAddress `json:"signer"`
}

// RouterKey Implements Msg.
func (msg MsgAddTrustedMetadataSchemaProposer) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgAddTrustedMetadataSchemaProposer) Type() string {
	return MsgTypeAddTrustedMetadataSchemaProposer
}

// ValidateBasic Implements Msg.
func (msg MsgAddTrustedMetadataSchemaProposer) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddTrustedMetadataSchemaProposer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddTrustedMetadataSchemaProposer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
