package types

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- MsgShareDocument
// ----------------------------------

type MsgShareDocument Document

func NewMsgShareDocument(document Document) MsgShareDocument {
	return MsgShareDocument(document)
}

// RouterKey Implements Msg.
func (msg MsgShareDocument) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgShareDocument) Type() string { return MsgTypeShareDocument }

// ValidateBasic Implements Msg.
func (msg MsgShareDocument) ValidateBasic() error {
	return Document(msg).Validate()
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

func NewMsgSendDocumentReceipt(receipt DocumentReceipt) MsgSendDocumentReceipt {
	return MsgSendDocumentReceipt(receipt)
}

// RouterKey Implements Msg.
func (msg MsgSendDocumentReceipt) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgSendDocumentReceipt) Type() string { return MsgTypeSendDocumentReceipt }

// ValidateBasic Implements Msg.
func (msg MsgSendDocumentReceipt) ValidateBasic() error {
	if !validateUUID(msg.UUID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid uuid: %s", msg.UUID))
	}

	if msg.Sender.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Sender.String())
	}

	if msg.Recipient.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Recipient.String())
	}

	if len(strings.TrimSpace(msg.TxHash)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Send Document's Transaction Hash can't be empty")
	}

	if !validateUUID(msg.DocumentUUID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid document UUID")
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
func (msg MsgAddSupportedMetadataSchema) ValidateBasic() error {
	if msg.Signer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Signer.String())
	}
	if err := msg.Schema.Validate(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
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
func (msg MsgAddTrustedMetadataSchemaProposer) ValidateBasic() error {
	if msg.Proposer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Proposer.String())
	}
	if msg.Signer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Signer.String())
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
