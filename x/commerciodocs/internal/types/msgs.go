package types

import (
	"encoding/hex"
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

/*
Author: Leonardo Bragagnolo
*/

const (
	MD5    = "md5"
	SHA1   = "sha-1"
	SHA224 = "sha-224"
	SHA256 = "sha-256"
	SHA384 = "sha-384"
	SHA512 = "sha-512"
)

// ----------------------------------
// --- ShareDocument
// ----------------------------------

type MsgShareDocument struct {
	types.Document
}

func NewMsgShareDocument(schemaUri string, schemaVersion string, docContentUri string, proof string, checksumValue string,
	checksumAlgorithm string, sender sdk.AccAddress, recipient sdk.AccAddress, metaContentUri string) MsgShareDocument {
	return MsgShareDocument{
		Document: types.Document{
			Sender:     sender,
			Recipient:  recipient,
			ContentUri: docContentUri,
			Metadata: types.DocumentMetadata{
				ContentUri: metaContentUri,
				Schema: types.DocumentMetadataSchema{
					Uri:     schemaUri,
					Version: schemaVersion,
				},
				Proof: proof,
			},
			Checksum: types.DocumentChecksum{
				Value:     checksumValue,
				Algorithm: checksumAlgorithm,
			},
		},
	}
}

// RouterKey Implements Msg.
func (msg MsgShareDocument) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgShareDocument) Type() string { return "share_document" }

//Basic validation of DocumentMetadata fields
func validateDocMetadata(docMetadata types.DocumentMetadata) sdk.Error {
	if len(docMetadata.ContentUri) == 0 {
		return sdk.ErrUnknownRequest("Metadata content uri can't be empty")
	}
	if len(docMetadata.Proof) == 0 {
		return sdk.ErrUnknownRequest("Computation Proof can't be empty")
	}
	if len(docMetadata.Schema.Uri) == 0 {
		return sdk.ErrUnknownRequest("Schema uri can't be empty")
	}
	if len(docMetadata.Schema.Version) == 0 {
		return sdk.ErrUnknownRequest("Schema version can't be empty")
	}
	return nil
}

//Checksum validation
func validateChecksum(checksum types.DocumentChecksum) sdk.Error {
	if len(checksum.Value) == 0 {
		return sdk.ErrUnknownRequest("Checksum value can't be empty")
	}
	if len(checksum.Algorithm) == 0 {
		return sdk.ErrUnknownRequest("Checksum algorithm can't be empty")
	}

	_, err := hex.DecodeString(checksum.Value)
	if err != nil {
		return sdk.ErrUnknownRequest("Invalid checksum value")
	}

	algorithm := strings.ToLower(checksum.Algorithm)
	if algorithm == MD5 && len(checksum.Value) != 32 {
		return sdk.ErrUnknownRequest("Invalid checksum length for MD5 hash")
	}
	if algorithm == SHA1 && len(checksum.Value) != 40 {
		return sdk.ErrUnknownRequest("Invalid checksum length for SHA1 hash")
	}
	if algorithm == SHA224 && len(checksum.Value) != 56 {
		return sdk.ErrUnknownRequest("Invalid checksum length for SHA224 hash")
	}
	if algorithm == SHA256 && len(checksum.Value) != 64 {
		return sdk.ErrUnknownRequest("Invalid checksum length for SHA256 hash")
	}
	if algorithm == SHA384 && len(checksum.Value) != 96 {
		return sdk.ErrUnknownRequest("Invalid checksum length for SHA384 hash")
	}
	if algorithm == SHA512 && len(checksum.Value) != 256 {
		return sdk.ErrUnknownRequest("Invalid checksum length for SHA512 hash")
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
