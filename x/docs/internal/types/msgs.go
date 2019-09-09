package types

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

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

type MsgShareDocument Document

func NewMsgShareDocument(document Document) MsgShareDocument {
	return MsgShareDocument(document)
}

// RouterKey Implements Msg.
func (msg MsgShareDocument) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgShareDocument) Type() string { return MsgTypeShareDocument }

func validateUuid(uuid string) bool {
	regex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	return regex.MatchString(uuid)
}

func validateDocMetadata(docMetadata DocumentMetadata) sdk.Error {
	if len(docMetadata.ContentUri) == 0 {
		return sdk.ErrUnknownRequest("MetadataSchema content URI can't be empty")
	}

	if (docMetadata.Schema == nil) && len(docMetadata.SchemaType) == 0 {
		return sdk.ErrUnknownRequest("Either schema or schema_type must be defined")
	}

	if docMetadata.Schema != nil {
		if len(docMetadata.Schema.Uri) == 0 {
			return sdk.ErrUnknownRequest("Schema URI can't be empty")
		}
		if len(docMetadata.Schema.Version) == 0 {
			return sdk.ErrUnknownRequest("Schema version can't be empty")
		}
	}

	if len(docMetadata.Proof) == 0 {
		return sdk.ErrUnknownRequest("Computation proof can't be empty")
	}
	return nil
}

func validateChecksum(checksum DocumentChecksum) sdk.Error {
	if len(checksum.Value) == 0 {
		return sdk.ErrUnknownRequest("Checksum value can't be empty")
	}
	if len(checksum.Algorithm) == 0 {
		return sdk.ErrUnknownRequest("Checksum algorithm can't be empty")
	}

	_, err := hex.DecodeString(checksum.Value)
	if err != nil {
		return sdk.ErrUnknownRequest("Invalid checksum value (must be hex)")
	}

	algorithm := strings.ToLower(checksum.Algorithm)

	// Check that the algorithm is valid
	length, ok := algorithms[algorithm]
	if !ok {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid algorithm type %s", algorithm))
	}

	// Check the validity of the checksum value
	if len(checksum.Value) != length {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid checksum length for algorithm %s", algorithm))
	}

	return nil
}

// ValidateBasic Implements Msg.
func (msg MsgShareDocument) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}
	if !validateUuid(msg.Uuid) {
		return sdk.ErrUnknownRequest("Invalid document UUID")
	}
	if len(msg.ContentUri) == 0 {
		return sdk.ErrUnknownRequest("Document content Uri can't be empty")
	}

	err := validateDocMetadata(msg.Metadata)
	if err != nil {
		return err
	}

	err = validateChecksum(msg.Checksum)
	if err != nil {
		return err
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

// ----------------------------------
// --- MsgAddSupportedMetadataSchema
// ----------------------------------

type MsgAddSupportedMetadataSchema struct {
	Signer sdk.AccAddress `json:"signer"`
	Schema MetadataSchema `json:"schema"`
}

// RouterKey Implements Msg.
func (msg MsgAddSupportedMetadataSchema) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgAddSupportedMetadataSchema) Type() string { return MsgTypeAddSupportedMetadataSchema }

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
