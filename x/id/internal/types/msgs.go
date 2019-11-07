package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetIdentity DidDocument

func NewMsgSetIdentity(document DidDocument) MsgSetIdentity {
	return MsgSetIdentity(document)
}

// Route Implements Msg.
func (msg MsgSetIdentity) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgSetIdentity) Type() string { return MsgTypeSetIdentity }

// ValidateBasic Implements Msg.
func (msg MsgSetIdentity) ValidateBasic() sdk.Error {
	return DidDocument(msg).Validate()
}

// GetSignBytes Implements Msg.
func (msg MsgSetIdentity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ID}
}

// ---------------------------
// --- MsgRequestDidDeposit
// ---------------------------

type MsgRequestDidDeposit struct {
	Recipient     sdk.AccAddress `json:"recipient"`      // Address that should be funded
	Amount        sdk.Coins      `json:"amount"`         // Amount that should be taken
	Proof         string         `json:"proof"`          // Proof of the deposit, encrypted using an AES-256 key and hex encoded
	EncryptionKey string         `json:"encryption_key"` // AES-256 key encrypted using reader's public key and hex encoded
	FromAddress   sdk.AccAddress `json:"from_address"`   // Address from which the funds should be taken
}

// Route Implements Msg.
func (msg MsgRequestDidDeposit) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgRequestDidDeposit) Type() string { return MsgTypeRequestDidDeposit }

// ValidateBasic Implements Msg.
func (msg MsgRequestDidDeposit) ValidateBasic() sdk.Error {
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid recipient: %s", msg.Recipient))
	}

	if !msg.Amount.IsValid() || msg.Amount.Empty() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("Deposit amount not valid: %s", msg.Amount))
	}

	if !ValidateHex(msg.Proof) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid proof: %s", msg.Proof))
	}

	if !ValidateHex(msg.EncryptionKey) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid encryption key value: %s", msg.EncryptionKey))
	}

	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid from_address: %s", msg.FromAddress))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRequestDidDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRequestDidDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// ------------------------
// --- MsgMoveDeposit
// ------------------------

type MsgMoveDeposit struct {
	DepositProof string         `json:"deposit_proof"`
	Signer       sdk.AccAddress `json:"signer"`
}

func NewMsgMoveDeposit(proof string, signer sdk.AccAddress) MsgMoveDeposit {
	return MsgMoveDeposit{
		DepositProof: proof,
		Signer:       signer,
	}
}

// Route Implements Msg.
func (msg MsgMoveDeposit) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgMoveDeposit) Type() string { return MsgTypeMoveDeposit }

// ValidateBasic Implements Msg.
func (msg MsgMoveDeposit) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid signer address: %s", msg.Signer))
	}

	if !ValidateHex(msg.DepositProof) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid deposit_proof: %s", msg.DepositProof))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMoveDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMoveDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// --------------------------------------
// --- MsgInvalidateDidDepositRequest
// ---------------------------------------

type MsgInvalidateDidDepositRequest struct {
	Editor       sdk.AccAddress `json:"editor"`
	DepositProof string         `json:"deposit_proof"`
	Status       RequestStatus  `json:"status"`
}

func NewMsgInvalidateDidDepositRequest(status RequestStatus, proof string, editor sdk.AccAddress) MsgInvalidateDidDepositRequest {
	return MsgInvalidateDidDepositRequest{
		Editor:       editor,
		DepositProof: proof,
		Status:       status,
	}
}

// Route Implements Msg.
func (msg MsgInvalidateDidDepositRequest) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgInvalidateDidDepositRequest) Type() string { return MsgTypeInvalidateDidDepositRequest }

// ValidateBasic Implements Msg.
func (msg MsgInvalidateDidDepositRequest) ValidateBasic() sdk.Error {
	if msg.Editor.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid editor address: %s", msg.Editor))
	}

	if !ValidateHex(msg.DepositProof) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid deposit_proof: %s", msg.DepositProof))
	}

	if err := msg.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgInvalidateDidDepositRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgInvalidateDidDepositRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Editor}
}

// ---------------------------
// --- MsgRequestDidPowerUp
// ---------------------------

type MsgRequestDidPowerUp struct {
	Claimant      sdk.AccAddress `json:"claimant"`
	Amount        sdk.Coins      `json:"amount"`
	Proof         string         `json:"proof"`
	EncryptionKey string         `json:"encryption_key"`
}

// Route Implements Msg.
func (msg MsgRequestDidPowerUp) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgRequestDidPowerUp) Type() string { return MsgTypeRequestDidPowerUp }

// ValidateBasic Implements Msg.
func (msg MsgRequestDidPowerUp) ValidateBasic() sdk.Error {
	if msg.Claimant.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid claimant: %s", msg.Claimant))
	}

	if !msg.Amount.IsValid() || msg.Amount.Empty() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("Power up amount not valid: %s", msg.Amount.String()))
	}

	if !ValidateHex(msg.Proof) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid proof: %s", msg.Proof))
	}

	if !ValidateHex(msg.EncryptionKey) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid encryption key value: %s", msg.EncryptionKey))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRequestDidPowerUp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRequestDidPowerUp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Claimant}
}

// ------------------------
// --- MsgPowerUpDid
// ------------------------

type MsgPowerUpDid struct {
	Recipient           sdk.AccAddress `json:"recipient"`
	Amount              sdk.Coins      `json:"amount"`
	ActivationReference string         `json:"activation_reference"`
	Signer              sdk.AccAddress `json:"signer"`
}

// Route Implements Msg.
func (msg MsgPowerUpDid) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgPowerUpDid) Type() string { return MsgTypePowerUpDid }

// ValidateBasic Implements Msg.
func (msg MsgPowerUpDid) ValidateBasic() sdk.Error {
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress("Power up recipient cannot be empty")
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress("Power up signer cannot be empty")
	}

	if msg.Amount.Empty() || !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("Invalid power up amount: %s", msg.Amount))
	}

	if !ValidateHex(msg.ActivationReference) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid proof: %s", msg.ActivationReference))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgPowerUpDid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgPowerUpDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// ---------------------------------------
// --- MsgInvalidateDidPowerUpRequest
// ---------------------------------------

type MsgInvalidateDidPowerUpRequest struct {
	PowerUpProof string         `json:"power_up_proof"`
	Status       RequestStatus  `json:"status"`
	Editor       sdk.AccAddress `json:"editor"`
}

func NewMsgInvalidateDidPowerUpRequest(status RequestStatus, proof string, editor sdk.AccAddress) MsgInvalidateDidPowerUpRequest {
	return MsgInvalidateDidPowerUpRequest{
		Editor:       editor,
		PowerUpProof: proof,
		Status:       status,
	}
}

// Route Implements Msg.
func (msg MsgInvalidateDidPowerUpRequest) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgInvalidateDidPowerUpRequest) Type() string { return MsgTypeInvalidateDidPowerUpRequest }

// ValidateBasic Implements Msg.
func (msg MsgInvalidateDidPowerUpRequest) ValidateBasic() sdk.Error {

	if msg.Editor.Empty() {
		return sdk.ErrInvalidAddress(fmt.Sprintf("Invalid editor address: %s", msg.Editor))
	}

	if !ValidateHex(msg.PowerUpProof) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid power_up_proof: %s", msg.PowerUpProof))
	}

	if err := msg.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgInvalidateDidPowerUpRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgInvalidateDidPowerUpRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Editor}
}
