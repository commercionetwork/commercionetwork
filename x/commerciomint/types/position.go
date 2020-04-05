package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Position represents a collateralized debt position that is open from a user in order to convert
// any currently priced token into stable Commercio Cash Credits.
type Position struct {
	Owner     sdk.AccAddress `json:"owner"`
	Deposit   sdk.Coin       `json:"deposit"`
	Credits   sdk.Coins      `json:"credits"`
	CreatedAt int64          `json:"timestamp"` // Block height at which the CDP has been created
}

func NewPosition(owner sdk.AccAddress, deposit sdk.Coin, liquidity sdk.Coins, timeStamp int64) Position {
	return Position{
		Owner:     owner,
		Deposit:   deposit,
		Credits:   liquidity,
		CreatedAt: timeStamp,
	}
}

// Validate verifies that the data contained inside this position are all valid,
// returning an error is something isn't valid
func (current Position) Validate() error {
	if current.Owner.Empty() {
		return fmt.Errorf("invalid owner address: %s", current.Owner)
	}
	if !current.Deposit.IsValid() || !current.Deposit.IsPositive() {
		return fmt.Errorf("invalid deposit amount: %s", current.Deposit)
	}
	if current.Credits.Empty() || current.Credits.IsAnyNegative() {
		return fmt.Errorf("invalid liquidity amount: %s", current.Credits)
	}
	if current.CreatedAt < 1 {
		return fmt.Errorf("invalid timestamp: %d", current.CreatedAt)
	}
	return nil
}

// Equals returns true if and only if the two Position instances are equal.
func (current Position) Equals(cdp Position) bool {
	return current.Owner.Equals(cdp.Owner) &&
		current.Deposit.IsEqual(cdp.Deposit) &&
		current.Credits.IsEqual(cdp.Credits) &&
		current.CreatedAt == cdp.CreatedAt
}
