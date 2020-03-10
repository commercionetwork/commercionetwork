package types

import (
	"encoding/base64"
	"fmt"

	uuid "github.com/satori/go.uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
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
func (msg MsgSetIdentity) ValidateBasic() error {
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
// --- MsgRequestDidPowerUp
// ---------------------------

type MsgRequestDidPowerUp struct {
	Claimant sdk.AccAddress `json:"claimant"`
	Amount   sdk.Coins      `json:"amount"`
	Proof    string         `json:"proof"`
	ID       string         `json:"id"`
}

// Route Implements Msg.
func (msg MsgRequestDidPowerUp) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgRequestDidPowerUp) Type() string { return MsgTypeRequestDidPowerUp }

// ValidateBasic Implements Msg.
func (msg MsgRequestDidPowerUp) ValidateBasic() error {
	if msg.Claimant.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (fmt.Sprintf("Invalid claimant: %s", msg.Claimant)))
	}

	if !msg.Amount.IsValid() || msg.Amount.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidCoins, (fmt.Sprintf("Power up amount not valid: %s", msg.Amount.String())))
	}

	if _, err := base64.StdEncoding.DecodeString(msg.Proof); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "proof must be base64-encoded")
	}

	if _, err := uuid.FromString(msg.ID); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "invalid ID, must be a valid UUID")
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
func (msg MsgPowerUpDid) ValidateBasic() error {
	if msg.Recipient.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (fmt.Sprintf("Invalid recipient address: %s", msg.Recipient)))
	}

	if msg.Amount.Empty() || !msg.Amount.IsValid() {
		return sdkErr.Wrap(sdkErr.ErrInvalidCoins, (fmt.Sprintf("Invalid power up amount: %s", msg.Amount)))
	}

	if !ValidateHex(msg.ActivationReference) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (fmt.Sprintf("Invalid activation_reference: %s", msg.ActivationReference)))
	}

	if msg.Signer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (fmt.Sprintf("Invalid signer address: %s", msg.Signer)))
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
func (msg MsgInvalidateDidPowerUpRequest) ValidateBasic() error {

	if msg.Editor.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (fmt.Sprintf("Invalid editor address: %s", msg.Editor)))
	}

	if !ValidateHex(msg.PowerUpProof) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (fmt.Sprintf("Invalid power_up_proof: %s", msg.PowerUpProof)))
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
