package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CDPRequest struct {
	Signer          sdk.AccAddress `json:"signer"`
	DepositedAmount sdk.Coins      `json:"deposited_amount"`
	Timestamp       string         `json:"timestamp"`
}

func (cdpr CDPRequest) Equals(request CDPRequest) bool {
	return cdpr.Signer.Equals(request.Signer) &&
		cdpr.DepositedAmount.IsEqual(request.DepositedAmount) &&
		cdpr.Timestamp == request.Timestamp
}
