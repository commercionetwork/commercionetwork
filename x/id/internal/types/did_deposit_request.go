package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidDepositRequestStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (status DidDepositRequestStatus) Validate() sdk.Error {
	if len(strings.TrimSpace(status.Type)) == 0 {
		return sdk.ErrUnknownRequest("Status type cannot be empty")
	}

	if status.Type != "approved" && status.Type != "rejected" && status.Type != "canceled" {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid status type: %s", status.Type))
	}

	return nil
}

type DidDepositRequest struct {
	Status        *DidDepositRequestStatus `json:"status"`         // Type of the request
	Recipient     sdk.AccAddress           `json:"recipient"`      // Address that should be funded
	Amount        sdk.Coins                `json:"amount"`         // Amount that should be taken
	DepositProof  string                   `json:"deposit_proof"`  // Proof of the deposit, encrypted using an AES-256 key and hex encoded
	EncryptionKey string                   `json:"encryption_key"` // AES-256 key encrypted using reader's public key and hex encoded
	FromAddress   sdk.AccAddress           `json:"from_address"`   // Address from which the funds should be taken
}

func (request DidDepositRequest) Validate() sdk.Error {
	if request.Recipient.Empty() {
		return sdk.ErrInvalidAddress(request.Recipient.String())
	}

	if request.Amount.Empty() {
		return sdk.ErrInvalidCoins("Deposit amount cannot be empty")
	}

	if request.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins("Deposit amount cannot be contain negative values")
	}

	if err := ValidateDepositProof(request.DepositProof); err != nil {
		return err
	}

	encryptionKey := strings.TrimSpace(request.EncryptionKey)
	if len(encryptionKey) == 0 {
		return sdk.ErrUnknownRequest("Encryption key cannot be empty")
	}

	if _, err := hex.DecodeString(encryptionKey); err != nil {
		return sdk.ErrUnknownRequest("Encryption key must be hex encoded")
	}

	if request.FromAddress.Empty() {
		return sdk.ErrInvalidAddress(request.FromAddress.String())
	}

	return nil
}

func ValidateDepositProof(proof string) sdk.Error {

	depositProof := strings.TrimSpace(proof)
	if len(depositProof) == 0 {
		return sdk.ErrUnknownRequest("Deposit proof cannot be empty")
	}

	if _, err := hex.DecodeString(depositProof); err != nil {
		return sdk.ErrUnknownRequest("Deposit proof must be hex encoded")
	}

	return nil
}
