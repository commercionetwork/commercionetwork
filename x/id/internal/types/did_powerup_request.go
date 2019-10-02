package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidPowerupRequestStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (status DidPowerupRequestStatus) Validate() sdk.Error {
	statusType := strings.ToLower(status.Type)
	if statusType != "accepted" && statusType != "rejected" && statusType != "canceled" {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Status type not valid: %s", status.Type))
	}

	return nil
}

type DidPowerupRequest struct {
	Status        *DidPowerupRequestStatus `json:"status"`
	Claimant      sdk.AccAddress           `json:"claimant"`
	Amount        sdk.Coins                `json:"amount"`
	Proof         string                   `json:"proof"`
	EncryptionKey string                   `json:"encryption_key"`
}

func (request DidPowerupRequest) Validate() sdk.Error {
	if request.Status != nil {
		if err := (*request.Status).Validate(); err != nil {
			return err
		}
	}

	if request.Claimant.Empty() {
		return sdk.ErrInvalidAddress(request.Claimant.String())
	}

	if request.Amount.Empty() {
		return sdk.ErrInvalidCoins("Powerup request amount cannot be empty")
	}

	if request.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins("Powerup request amount cannot contain negative values")
	}

	if err := ValidateProof(request.Proof); err != nil {
		return err
	}

	if err := ValidateEncryptionKey(request.EncryptionKey); err != nil {
		return err
	}

	return nil
}
