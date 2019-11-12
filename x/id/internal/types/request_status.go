package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RequestStatus represents the status that can be associated to a
// Did Deposit request or a Did Power Up request
type RequestStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewRequestStatus(statusType, message string) RequestStatus {
	return RequestStatus{
		Type:    statusType,
		Message: message,
	}
}

// Validate returns an error if something present inside the status is wrong
func (status RequestStatus) Validate() sdk.Error {
	statusType := strings.ToLower(status.Type)
	if statusType != StatusRejected && statusType != StatusCanceled {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid status type: %s", status.Type))
	}
	return nil
}
