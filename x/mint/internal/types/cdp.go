package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Cdp stands for Collateralized Debt Position
type Cdp struct {
	Owner           sdk.AccAddress `json:"owner"`
	DepositedAmount sdk.Coins      `json:"deposited_amount"`
	CreditsAmount   sdk.Coins      `json:"credits_amount"`
	Timestamp       time.Time      `json:"timestamp"`
}

func (current Cdp) Validate() error {
	if current.Owner.Empty() {
		return sdk.ErrInvalidAddress(current.Owner.String())
	}
	if current.DepositedAmount.Empty() || current.DepositedAmount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(current.DepositedAmount.String())
	}
	if current.CreditsAmount.Empty() || current.CreditsAmount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(current.CreditsAmount.String())
	}
	if current.Timestamp.IsZero() {
		return sdk.ErrUnknownRequest("timestamp not valid")
	}
	return nil
}

func NewCdp(request CdpRequest, liquidityAmount sdk.Coins) Cdp {
	return Cdp{
		Owner:           request.Signer,
		DepositedAmount: request.DepositedAmount,
		CreditsAmount:   liquidityAmount,
		Timestamp:       request.Timestamp,
	}
}

func (current Cdp) Equals(cdp Cdp) bool {
	return current.Owner.Equals(cdp.Owner) &&
		current.DepositedAmount.IsEqual(cdp.DepositedAmount) &&
		current.CreditsAmount.IsEqual(cdp.CreditsAmount) &&
		current.Timestamp == cdp.Timestamp
}

type Cdps []Cdp

func (cdps Cdps) AppendIfMissing(cdp Cdp) (Cdps, bool) {
	for _, ele := range cdps {
		if ele.Equals(cdp) {
			return nil, true
		}
	}
	return append(cdps, cdp), false
}

//This method filters a slice without allocating a new underlying array
func (cdps Cdps) RemoveWhenFound(timestamp time.Time) (Cdps, bool) {
	tmp := cdps[:0]
	removed := false
	for _, ele := range cdps {
		if !ele.Timestamp.Equal(timestamp) {
			tmp = append(tmp, ele)
		} else {
			removed = true
		}
	}
	return tmp, removed
}
