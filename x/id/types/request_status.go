package types

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// RequestStatus represents the status that can be associated to a
// Did Credits request or a Did Power Up request
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

	switch statusType {
	case StatusRejected, StatusApproved, StatusCanceled:
		return nil
	default:
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid status type: %s", status.Type))

	}
}
