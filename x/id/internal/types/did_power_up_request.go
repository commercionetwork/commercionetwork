package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidPowerUpRequestStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (status DidPowerUpRequestStatus) Validate() sdk.Error {
	statusType := strings.ToLower(status.Type)
	if err := ValidateStatus(statusType); err != nil {
		return err
	}

	return nil
}

type DidPowerUpRequest struct {
	Status        *DidPowerUpRequestStatus `json:"status"`
	Claimant      sdk.AccAddress           `json:"claimant"`
	Amount        sdk.Coins                `json:"amount"`
	Proof         string                   `json:"proof"`
	EncryptionKey string                   `json:"encryption_key"`
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

	if request.Amount.Empty() {
		return sdk.ErrInvalidCoins("PowerUp request amount cannot be empty")
	}

	if request.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins("PowerUp request amount cannot contain negative values")
	}

	if err := ValidateProof(request.Proof); err != nil {
		return err
	}

	if err := ValidateEncryptionKey(request.EncryptionKey); err != nil {
		return err
	}

	return nil
}
