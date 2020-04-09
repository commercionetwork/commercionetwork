package types

import (
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"
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
	ProofKey string         `json:"proof_key"`
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

	if _, err := base64.StdEncoding.DecodeString(msg.ProofKey); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "proof key must be base64-encoded")
	}

	if _, err := uuid.FromString(msg.ID); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "invalid PowerUpID, must be a valid UUID")
	}

	if msg.ProofKey == "" {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "proof key cannot be empty")
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

// ---------------------------
// --- MsgChangePowerUpStatus
// ---------------------------

type MsgChangePowerUpStatus struct {
	Recipient sdk.AccAddress `json:"recipient"`
	PowerUpID string         `json:"id"`
	Status    RequestStatus  `json:"status"`
	Signer    sdk.AccAddress `json:"signer"`
}

// Route Implements Msg.
func (msg MsgChangePowerUpStatus) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgChangePowerUpStatus) Type() string { return MsgTypeChangePowerUpStatus }

// ValidateBasic Implements Msg.
func (msg MsgChangePowerUpStatus) ValidateBasic() error {
	if msg.Recipient.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (fmt.Sprintf("Invalid recipient address: %s", msg.Recipient)))
	}

	if _, err := uuid.FromString(msg.PowerUpID); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "invalid PowerUpID, must be a valid UUID")
	}

	if msg.Signer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (fmt.Sprintf("Invalid signer address: %s", msg.Signer)))
	}

	if err := msg.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgChangePowerUpStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgChangePowerUpStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
