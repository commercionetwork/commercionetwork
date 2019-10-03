package types

import (
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
