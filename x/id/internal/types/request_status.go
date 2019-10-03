package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RequestStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (status RequestStatus) Validate() sdk.Error {
	if err := ValidateStatus(status.Type); err != nil {
		return err
	}
	return nil
}

func ValidateStatus(status string) sdk.Error {
	statusType := strings.ToLower(status)
	if statusType != StatusRejected && statusType != StatusCanceled {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid status type: %s", status))
	}
	return nil
}
