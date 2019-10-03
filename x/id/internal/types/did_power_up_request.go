package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DidDepositRequest represents the request that is sent from a user when he wants to send
// something to his pairwise Did. This request will be read and unencrypted from a central
// identity that will later update the status and send the funds to the pairwise Did
type DidPowerUpRequest struct {
	Status        *RequestStatus `json:"status"`
	Claimant      sdk.AccAddress `json:"claimant"`
	Amount        sdk.Coins      `json:"amount"`
	Proof         string         `json:"proof"`
	EncryptionKey string         `json:"encryption_key"`
}

func (request DidPowerUpRequest) Validate() sdk.Error {
	if request.Status != nil {
		if err := (*request.Status).Validate(); err != nil {
			return err
		}
	}

	if request.Claimant.Empty() {
		return sdk.ErrInvalidAddress(request.Claimant.String())
	}

	if !request.Amount.IsValid() || request.Amount.Empty() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("PowerUp request amount not valid: %s", request.Amount.String()))
	}

	if request.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins("PowerUp request amount cannot contain negative values")
	}

	if err := ValidateHex(request.Proof); err != nil {
		return err
	}

	if err := ValidateEncryptionKey(request.EncryptionKey); err != nil {
		return err
	}

	return nil
}
