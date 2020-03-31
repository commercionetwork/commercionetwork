package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ------------
// --- CDP
// ------------

// Cdp represents a Collateralized Debt position that is open from a user in order to convert
// any currently priced token into stable Commercio Cash Credits
type Cdp struct {
	Owner           sdk.AccAddress `json:"owner"`
	DepositedAmount sdk.Coins      `json:"deposited_amount"`
	CreditsAmount   sdk.Coins      `json:"credits_amount"`
	Timestamp       int64          `json:"timestamp"` // Block height at which the CDP has been created
}

func NewCdp(owner sdk.AccAddress, depositAmount sdk.Coins, liquidityAmount sdk.Coins, timeStamp int64) Cdp {
	return Cdp{
		Owner:           owner,
		DepositedAmount: depositAmount,
		CreditsAmount:   liquidityAmount,
		Timestamp:       timeStamp,
	}
}

// Validate verifies that the data contained inside this CDP are all valid,
// returning an error is something isn't valid
func (current Cdp) Validate() error {
	if current.Owner.Empty() {
		return fmt.Errorf("invalid owner address: %s", current.Owner)
	}
	if current.DepositedAmount.Empty() || current.DepositedAmount.IsAnyNegative() {
		return fmt.Errorf("invalid deposit amount: %s", current.DepositedAmount)
	}
	if current.CreditsAmount.Empty() || current.CreditsAmount.IsAnyNegative() {
		return fmt.Errorf("invalid liquidity amount: %s", current.CreditsAmount)
	}
	if current.Timestamp == 0 {
		return fmt.Errorf("invalid timestamp: %d", current.Timestamp)
	}
	return nil
}

// Equals returns true if and only if the two CDPs instances contain the same data
func (current Cdp) Equals(cdp Cdp) bool {
	return current.Owner.Equals(cdp.Owner) &&
		current.DepositedAmount.IsEqual(cdp.DepositedAmount) &&
		current.CreditsAmount.IsEqual(cdp.CreditsAmount) &&
		current.Timestamp == cdp.Timestamp
}

// -------------
// --- CDPs
// -------------

// Cdps represents a slice of CDP objects
type Cdps []Cdp

// AppendIfMissing appends the given cdp to the list of cdps if it does not exist inside it yet,
// returning also true if the object has been appended successfully
func (cdps Cdps) AppendIfMissing(cdp Cdp) (Cdps, bool) {
	for _, ele := range cdps {
		if ele.Equals(cdp) {
			return nil, false
		}
	}
	return append(cdps, cdp), true
}

// RemoveWhenFound filters a slice without allocating a new underlying array
func (cdps Cdps) RemoveWhenFound(timestamp int64) (Cdps, bool) {
	tmp := cdps[:0]
	removed := false
	for _, ele := range cdps {
		if ele.Timestamp != timestamp {
			tmp = append(tmp, ele)
		} else {
			removed = true
		}
	}
	return tmp, removed
}
