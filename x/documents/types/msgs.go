package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (msg *MsgShareDocument) ValidateBasic() error {
	var document = Document(*msg)

	if err := document.Validate(); err != nil {
		return fmt.Errorf("invalid document: %e", err)
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
	receipt := DocumentReceipt(*msg)

	if err := receipt.Validate(); err != nil {
		return fmt.Errorf("invalid document receipt: %e", err)
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
