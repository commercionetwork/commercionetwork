package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
)

func NewPosition(owner sdk.AccAddress, deposit sdk.Int, liquidity sdk.Coin, id string, createdAt time.Time, exchangeRate sdk.Dec) Position {

	return Position{
		Owner:        owner.String(),
		Collateral:   deposit.ToDec().RoundInt64(), // TODO FIX THIS
		Credits:      &liquidity,
		ID:           id,
		CreatedAt:    &createdAt, // TODO FIX THIS
		ExchangeRate: exchangeRate,
	}
}

// Validate verifies that the data contained inside this position are all valid,
// returning an error is something isn't valid
func (pos Position) Validate() error {
	if _, err := sdk.AccAddressFromBech32(pos.Owner); err != nil {
		return err
	}
	if pos.Collateral <= 0 {
		return fmt.Errorf("invalid collateral amount")
	}

	//TODO COMPLETE CONTROLS
	if !ValidateCredits(*pos.Credits) {
		return fmt.Errorf("invalid liquidity amount: %s", pos.Credits)
	}

	if pos.ExchangeRate.IsNegative() {
		return fmt.Errorf("exchange rate cannot be negative")
	}

	if pos.CreatedAt != nil {
		if *pos.CreatedAt == (time.Time{}) {
			return fmt.Errorf("cannot have empty creation time")
		}
	}

	if _, err := uuid.FromString(pos.ID); err != nil {
		return fmt.Errorf("id string must be a well-defined UUID")
	}

	return nil
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

func ValidateCredits(credits sdk.Coin) bool {
	if credits.IsValid() && credits.IsPositive() {
		return true
	}
	return false
}

// Equals returns true if and only if the two Position instances are equal.
func (pos Position) Equals(etp Position) bool {
	posOwner, _ := sdk.AccAddressFromBech32(pos.Owner)
	etpOwner, _ := sdk.AccAddressFromBech32(etp.Owner)
	posCollateral := sdk.NewInt(pos.Collateral)
	etpCollateral := sdk.NewInt(etp.Collateral)

	return posOwner.Equals(etpOwner) &&
		posCollateral.Equal(etpCollateral) &&
		pos.Credits.IsEqual(*etp.Credits) &&
		pos.ID == etp.ID &&
		pos.ExchangeRate.Equal(etp.ExchangeRate) &&
		pos.CreatedAt == etp.CreatedAt
}
