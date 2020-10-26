package types

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
)

// denom used by the minted tokens
const CreditsDenom = "uccc"

// Position represents a exchange trade position that is open from a user in order to convert
// any currently priced token into Commercio Cash Credits.
type Position struct {
	Owner        sdk.AccAddress `json:"owner"`
	Collateral   sdk.Int        `json:"collateral"`
	Credits      sdk.Coin       `json:"credits"`
	CreatedAt    time.Time      `json:"created_at"`
	ID           string         `json:"id"`
	ExchangeRate sdk.Int        `json:"exchange_rate"`
}

func NewPosition(owner sdk.AccAddress, deposit sdk.Int, liquidity sdk.Coin, id string, createdAt time.Time, exchangeRate sdk.Int) Position {
	return Position{
		Owner:        owner,
		Collateral:   deposit,
		Credits:      liquidity,
		ID:           id,
		CreatedAt:    createdAt,
		ExchangeRate: exchangeRate,
	}
}

// Validate verifies that the data contained inside this position are all valid,
// returning an error is something isn't valid
func (pos Position) Validate() error {
	if pos.Owner.Empty() {
		return fmt.Errorf("invalid owner address: %s", pos.Owner)
	}
	if pos.Collateral.IsZero() || pos.Collateral.IsNegative() {
		return errors.New("invalid collateral amount")
	}
	if !ValidateCredits(pos.Credits) {
		return fmt.Errorf("invalid liquidity amount: %s", pos.Credits)
	}

	if pos.ExchangeRate.IsNegative() {
		return fmt.Errorf("exchange rate cannot be zero")
	}

	if pos.CreatedAt == (time.Time{}) {
		return fmt.Errorf("cannot have empty creation time")
	}

	if _, err := uuid.FromString(pos.ID); err != nil {
		return fmt.Errorf("id string must be a well-defined UUID")
	}

	return nil
}

// Equals returns true if and only if the two Position instances are equal.
func (pos Position) Equals(etp Position) bool {
	return pos.Owner.Equals(etp.Owner) &&
		pos.Collateral.Equal(etp.Collateral) &&
		pos.Credits.IsEqual(etp.Credits) &&
		pos.ID == etp.ID &&
		pos.ExchangeRate.Equal(etp.ExchangeRate) &&
		pos.CreatedAt.Equal(etp.CreatedAt)
}

func ValidateCredits(credits sdk.Coin) bool {
	if credits.IsValid() && credits.IsPositive() {
		return true
	}
	return false
}

func ValidateDeposit(deposit sdk.Coins) bool {
	for _, coin := range deposit {
		if coin.Denom != CreditsDenom {
			return false
		}
	}

	if !deposit.IsValid() || deposit.Empty() || !deposit.IsAllPositive() {
		return false
	}

	return true
}
