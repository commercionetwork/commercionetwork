package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewPosition(owner sdk.AccAddress, deposit sdk.Int, liquidity sdk.Coin, id string, createdAt time.Time, exchangeRate sdk.DecProto) Position {

	return Position{
		Owner:        owner.String(),
		Collateral:   deposit.ToDec().RoundInt64(), // TODO FIX THIS
		Credits:      &liquidity,
		ID:           id,
		CreatedAt:    createdAt.String(), // TODO FIX THIS
		ExchangeRate: &exchangeRate,
	}
}

// Validate verifies that the data contained inside this position are all valid,
// returning an error is something isn't valid
func (pos Position) Validate() error {
	//if pos.Owner.Empty() {
	if pos.Owner == "" {
		return fmt.Errorf("invalid owner address: %s", pos.Owner)
	}
	if pos.Collateral == 0 || pos.Collateral < 0 {
		//return errors.New("invalid collateral amount")
		return fmt.Errorf("invalid collateral amount")
	}

	//TODO COMPLETE CONTROLS
	if !ValidateCredits(*pos.Credits) {
		return fmt.Errorf("invalid liquidity amount: %s", pos.Credits)
	}

	if pos.ExchangeRate.Dec.IsNegative() {
		return fmt.Errorf("exchange rate cannot be zero")
	}

	createdAt, err := time.Parse(time.RFC3339, pos.CreatedAt)
	if err != nil || createdAt == (time.Time{}) {
		return fmt.Errorf("cannot have empty creation time")
	}

	// TODO control repeated control
	/*
		if _, err := uuid.FromString(pos.ID); err != nil {
			return fmt.Errorf("id string must be a well-defined UUID")
		}*/

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
	posOwner := sdk.AccAddress(pos.Owner)
	etpOwner := sdk.AccAddress(etp.Owner)
	posCollateral := sdk.NewInt(pos.Collateral)
	etpCollateral := sdk.NewInt(etp.Collateral)

	return posOwner.Equals(etpOwner) &&
		posCollateral.Equal(etpCollateral) &&
		pos.Credits.IsEqual(*etp.Credits) &&
		pos.ID == etp.ID &&
		pos.ExchangeRate.Dec.Equal(etp.ExchangeRate.Dec) &&
		pos.CreatedAt == etp.CreatedAt
}
