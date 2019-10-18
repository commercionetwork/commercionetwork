package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CdpRequest struct {
	Signer          sdk.AccAddress `json:"signer"`
	DepositedAmount sdk.Coins      `json:"deposit_amount"`
	Timestamp       time.Time      `json:"timestamp"`
}

func (cdpr CdpRequest) Equals(request CdpRequest) bool {
	return cdpr.Signer.Equals(request.Signer) &&
		cdpr.DepositedAmount.IsEqual(request.DepositedAmount) &&
		cdpr.Timestamp == request.Timestamp
}
