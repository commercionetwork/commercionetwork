package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Position represents a collateralized debt position that is open from a user in order to convert
// any currently priced token into stable Commercio Cash Credits.
type Position struct {
	Owner     sdk.AccAddress `json:"owner"`
	Deposit   sdk.Coins      `json:"deposit"`
	Credits   sdk.Coin       `json:"credits"`
	CreatedAt int64          `json:"timestamp"` // Block height at which the CDP has been created
}

func NewPosition(owner sdk.AccAddress, deposit sdk.Coins, liquidity sdk.Coin, timeStamp int64) Position {
	return Position{
		Owner:     owner,
		Deposit:   deposit,
		Credits:   liquidity,
		CreatedAt: timeStamp,
	}
}

// Validate verifies that the data contained inside this position are all valid,
// returning an error is something isn't valid
func (pos Position) Validate() error {
	if pos.Owner.Empty() {
		return fmt.Errorf("invalid owner address: %s", pos.Owner)
	}
	if !ValidateDeposit(pos.Deposit) {
		return fmt.Errorf("invalid deposit amount: %s", pos.Deposit)
	}
	if !ValidateCredits(pos.Credits) {
		return fmt.Errorf("invalid liquidity amount: %s", pos.Credits)
	}
	if pos.CreatedAt < 1 {
		return fmt.Errorf("invalid timestamp: %d", pos.CreatedAt)
	}
	return nil
}

// Equals returns true if and only if the two Position instances are equal.
func (pos Position) Equals(cdp Position) bool {
	return pos.Owner.Equals(cdp.Owner) &&
		pos.Deposit.IsEqual(cdp.Deposit) &&
		pos.Credits.IsEqual(cdp.Credits) &&
		pos.CreatedAt == cdp.CreatedAt
}

func ValidateCredits(credits sdk.Coin) bool {
	if credits.IsValid() && credits.IsPositive() {
		return true
	}
	return false
}

func ValidateDeposit(deposit sdk.Coins) bool {
	if !deposit.IsValid() || deposit.Empty() || !deposit.IsAllPositive() {
		return false
	}
	return true
}
