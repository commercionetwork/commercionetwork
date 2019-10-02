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
// --- MsgChangeDidDepositRequestStatus
// ---------------------------------------

type MsgChangeDidDepositRequestStatus struct {
	Editor       sdk.AccAddress          `json:"editor"`
	DepositProof string                  `json:"deposit_proof"`
	Status       DidDepositRequestStatus `json:"status"`
}

// Route Implements Msg.
func (msg MsgChangeDidDepositRequestStatus) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgChangeDidDepositRequestStatus) Type() string { return MsgTypeEditDidDepositRequest }

// ValidateBasic Implements Msg.
func (msg MsgChangeDidDepositRequestStatus) ValidateBasic() sdk.Error {
	if msg.Editor.Empty() {
		return sdk.ErrInvalidAddress(msg.Editor.String())
	}

	if err := ValidateProof(msg.DepositProof); err != nil {
		return err
	}

	if err := msg.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgChangeDidDepositRequestStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgChangeDidDepositRequestStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Editor}
}

// ---------------------------
// --- MsgRequestDidPowerUp
// ---------------------------

type MsgRequestDidPowerUp DidPowerUpRequest

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
// --- MsgChangeDidPowerUpRequestStatus
// ---------------------------------------

type MsgChangeDidPowerUpRequestStatus struct {
	PowerUpProof string                  `json:"PowerUp_proof"`
	Status       DidPowerUpRequestStatus `json:"status"`
	Editor       sdk.AccAddress          `json:"signer"`
}

// Route Implements Msg.
func (msg MsgChangeDidPowerUpRequestStatus) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgChangeDidPowerUpRequestStatus) Type() string { return MsgTypeEditDidPowerUpRequest }

// ValidateBasic Implements Msg.
func (msg MsgChangeDidPowerUpRequestStatus) ValidateBasic() sdk.Error {
	if err := ValidateProof(msg.PowerUpProof); err != nil {
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
func (msg MsgChangeDidPowerUpRequestStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgChangeDidPowerUpRequestStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Editor}
}
