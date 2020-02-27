package types

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
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
func (status RequestStatus) Validate() error {
	statusType := strings.ToLower(status.Type)
	if statusType != StatusRejected && statusType != StatusCanceled {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (fmt.Sprintf("Invalid status type: %s", status.Type)))
	}
	return nil
}
