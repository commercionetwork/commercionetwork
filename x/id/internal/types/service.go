package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Service represents any type of service the subject wishes to advertise,
// including decentralized identity management services for further discovery,
// authentication, authorization, or interaction.
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// Equals returns true service and other contain the same data
func (service Service) Equals(other Service) bool {
	return service.ID == other.ID &&
		service.Type == other.Type &&
		service.ServiceEndpoint == other.ServiceEndpoint
}

// Validate checks the data present inside service and returns an
// error if something is invalid
func (service Service) Validate() sdk.Error {
	if len(strings.TrimSpace(service.ID)) == 0 {
		return sdk.ErrUnknownRequest("Service id cannot be empty")
	}

	if len(strings.TrimSpace(service.Type)) == 0 {
		return sdk.ErrUnknownRequest("Service type cannot be empty")
	}

	if len(strings.TrimSpace(service.ServiceEndpoint)) == 0 {
		return sdk.ErrUnknownRequest("Service endpoint cannot be empty")
	}

	return nil
}

// Services represents a list of Service objects
type Services []Service

// Equals allows to easily tells if two Services objects contain the same
// data in the same order
func (services Services) Equals(other Services) bool {
	if len(services) != len(other) {
		return false
	}

	for index, service := range services {
		if !service.Equals(other[index]) {
			return false
		}
	}
	return true
}
