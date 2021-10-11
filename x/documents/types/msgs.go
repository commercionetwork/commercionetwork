package types

import (
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gofrs/uuid"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

// ----------------------------------
// --- MsgShareDocument
// ----------------------------------

var _ sdk.Msg = &MsgShareDocument{}

func NewMsgShareDocument(document Document) *MsgShareDocument {
	return &MsgShareDocument{
		Sender:         document.Sender,
		Recipients:     document.Recipients,
		UUID:           document.UUID,
		Metadata:       document.Metadata,
		ContentURI:     document.ContentURI,
		Checksum:       document.Checksum,
		EncryptionData: document.EncryptionData,
		DoSign:         document.DoSign,
	}
}

func (msg *MsgShareDocument) Route() string {
	return ModuleName
}

func (msg *MsgShareDocument) Type() string {
	return MsgTypeShareDocument
}

func (msg *MsgShareDocument) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgShareDocument) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

/*
func (msg *MsgShareDocument) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkErr.Wrapf(sdkErr.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}*/

func (msg *MsgShareDocument) ValidateBasic() error {
	if msg.Sender == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Sender)
	}

	if len(msg.Recipients) == 0 {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Recipients cannot be empty")
	}

	for _, recipient := range msg.Recipients {
		if recipient == "" {
			return sdkErr.Wrap(sdkErr.ErrInvalidAddress, recipient)
		}
	}

	if !validateUUID(msg.UUID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid document UUID: %s", msg.UUID))
	}

	err := msg.Metadata.Validate()
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	if msg.Checksum != nil {
		err = msg.Checksum.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
		}
	}

	if msg.EncryptionData != nil {
		err = msg.EncryptionData.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
		}
	}

	if msg.EncryptionData != nil {

		// check that each document recipient have some encrypted data
		for _, recipient := range msg.Recipients {
			recipientAccAddr, _ := sdk.AccAddressFromBech32(recipient)
			if !msg.EncryptionData.ContainsRecipient(recipientAccAddr) {
				errMsg := fmt.Sprintf(
					"%s is a recipient inside the document but not in the encryption data",
					recipient,
				)
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
			}
		}

		// check that there are no spurious encryption data recipients not present
		// in the document recipient list
		for _, encAdd := range msg.EncryptionData.Keys {
			var recipient ctypes.Strings = msg.Recipients
			if !recipient.Contains(encAdd.Recipient) {
				errMsg := fmt.Sprintf(
					"%s is a recipient inside encryption data but not inside the message",
					encAdd.Recipient,
				)
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
			}
		}

		// Check that the `encrypted_data' field name is actually present in msg
		fNotPresent := func(s string) error {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("field \"%s\" not present in document, but marked as encrypted", s),
			)
		}

		for _, fieldName := range msg.EncryptionData.EncryptedData {
			switch fieldName {
			case "content_uri":
				if msg.ContentURI == "" {
					return fNotPresent("content_uri")
				}
			case "metadata.schema.uri":
				if msg.Metadata.Schema == nil || msg.Metadata.Schema.URI == "" {
					return fNotPresent("metadata.schema.uri")
				}
			}
		}

	}

	if msg.DoSign != nil {
		if msg.Checksum == nil {
			return sdkErr.Wrap(
				sdkErr.ErrUnknownRequest,
				"field \"checksum\" not present in document, but required when using do_sign",
			)
		}

		if msg.ContentURI == "" {
			return sdkErr.Wrap(
				sdkErr.ErrUnknownRequest,
				"field \"content_uri\" not present in document, but required when using do_sign",
			)
		}
		var snData SdnData = msg.DoSign.SdnData
		err := snData.Validate()
		if err != nil {
			return sdkErr.Wrap(
				sdkErr.ErrUnknownRequest,
				err.Error(),
			)
		}
	}

	msgToDocument := Document{
		Sender:         msg.Sender,
		Recipients:     msg.Recipients,
		UUID:           msg.UUID,
		Metadata:       msg.Metadata,
		ContentURI:     msg.ContentURI,
		Checksum:       msg.Checksum,
		EncryptionData: msg.EncryptionData,
		DoSign:         msg.DoSign,
	}
	if err := msgToDocument.lengthLimits(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest,
			err.Error(),
		)
	}

	return nil
}

// ----------------------------------
// --- MsgSendDocumentReceipt
// ----------------------------------

var _ sdk.Msg = &MsgSendDocumentReceipt{}

func NewMsgSendDocumentReceipt(uuid string, sender string, recipient string, txHash string, documentUUID string, proof string) *MsgSendDocumentReceipt {
	return &MsgSendDocumentReceipt{
		UUID:         uuid,
		Sender:       sender,
		Recipient:    recipient,
		TxHash:       txHash,
		DocumentUUID: documentUUID,
		Proof:        proof,
	}
}

// RouterKey Implements Msg.
func (msg *MsgSendDocumentReceipt) Route() string {
	return ModuleName
}

// Type Implements Msg.
func (msg *MsgSendDocumentReceipt) Type() string {
	return MsgTypeSendDocumentReceipt
}

// ValidateBasic Implements Msg.
func (msg *MsgSendDocumentReceipt) ValidateBasic() error {
	_, err := uuid.FromString(msg.UUID)

	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid uuid: %s", msg.UUID))
	}

	if msg.Sender == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Sender)
	}

	if msg.Recipient == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Recipient)
	}

	if len(strings.TrimSpace(msg.TxHash)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Send Document's Transaction Hash can't be empty")
	}

	_, err = uuid.FromString(msg.DocumentUUID)

	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid document UUID")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg *MsgSendDocumentReceipt) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg *MsgSendDocumentReceipt) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
