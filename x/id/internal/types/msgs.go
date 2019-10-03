package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------
// --- MsgSetIdentity
// ----------------------

type MsgSetIdentity struct {
	Owner       sdk.AccAddress `json:"owner"`
	DidDocument DidDocument    `json:"did_document"`
}

func NewMsgSetIdentity(owner sdk.AccAddress, document DidDocument) MsgSetIdentity {
	return MsgSetIdentity{
		Owner:       owner,
		DidDocument: document,
	}
}

// Route Implements Msg.
func (msg MsgSetIdentity) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgSetIdentity) Type() string { return MsgTypeSetIdentity }

// ValidateBasic Implements Msg.
func (msg MsgSetIdentity) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	if err := msg.DidDocument.Validate(); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetIdentity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ---------------------------
// --- MsgRequestDidDeposit
// ---------------------------

type MsgRequestDidDeposit DidDepositRequest

func NewMsgRequestDidDeposit(request DidDepositRequest) MsgRequestDidDeposit {
	return MsgRequestDidDeposit(request)
}

// Route Implements Msg.
func (msg MsgRequestDidDeposit) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgRequestDidDeposit) Type() string { return MsgTypeRequestDidDeposit }

// ValidateBasic Implements Msg.
func (msg MsgRequestDidDeposit) ValidateBasic() sdk.Error {
	return DidDepositRequest(msg).Validate()
}

// GetSignBytes Implements Msg.
func (msg MsgRequestDidDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRequestDidDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// --------------------------------------
// --- MsgInvalidateDidDepositRequest
// ---------------------------------------

type MsgInvalidateDidDepositRequest struct {
	Editor       sdk.AccAddress `json:"editor"`
	DepositProof string         `json:"deposit_proof"`
	Status       RequestStatus  `json:"status"`
}

func NewMsgInvalidateDidDepositRequest(status RequestStatus, proof string,
	editor sdk.AccAddress) MsgInvalidateDidDepositRequest {
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
		return sdk.ErrInvalidAddress(msg.Editor.String())
	}

	if err := ValidateHex(msg.DepositProof); err != nil {
		return err
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

type MsgRequestDidPowerUp DidPowerUpRequest

func NewMsgRequestDidPowerUp(request DidPowerUpRequest) MsgRequestDidPowerUp {
	return MsgRequestDidPowerUp(request)
}

// Route Implements Msg.
func (msg MsgRequestDidPowerUp) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgRequestDidPowerUp) Type() string { return MsgTypeRequestDidPowerUp }

// ValidateBasic Implements Msg.
func (msg MsgRequestDidPowerUp) ValidateBasic() sdk.Error {
	return DidPowerUpRequest(msg).Validate()
}

// GetSignBytes Implements Msg.
func (msg MsgRequestDidPowerUp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRequestDidPowerUp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Claimant}
}

// ---------------------------------------
// --- MsgInvalidateDidPowerUpRequest
// ---------------------------------------

type MsgInvalidateDidPowerUpRequest struct {
	PowerUpProof string         `json:"power_up_proof"`
	Status       RequestStatus  `json:"status"`
	Editor       sdk.AccAddress `json:"editor"`
}

func NewMsgInvalidateDidPowerUpRequest(status RequestStatus, proof string,
	editor sdk.AccAddress) MsgInvalidateDidPowerUpRequest {
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
	if err := ValidateHex(msg.PowerUpProof); err != nil {
		return err
	}

	if err := msg.Status.Validate(); err != nil {
		return err
	}

	if msg.Editor.Empty() {
		return sdk.ErrInvalidAddress(msg.Editor.String())
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
func (msg MsgMoveDeposit) Type() string { return MsgTypeWithdrawDeposit }

// ValidateBasic Implements Msg.
func (msg MsgMoveDeposit) ValidateBasic() sdk.Error {
	if err := ValidateHex(msg.DepositProof); err != nil {
		return err
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
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
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}

	if msg.Amount.IsValid() || msg.Amount.Empty() || msg.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	if err := ValidateHex(msg.ActivationReference); err != nil {
		return err
	}

	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
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
