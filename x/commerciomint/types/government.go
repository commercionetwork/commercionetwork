package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateConversionRate(rate sdk.Dec) error {
	if rate.IsZero() {
		return fmt.Errorf("conversion rate cannot be zero")
	}
	if rate.IsNegative() {
		return fmt.Errorf("conversion rate must be positive")
	}
	return nil
}

func ValidateFreezePeriod(freezePeriod time.Duration) error {
	if freezePeriod.Seconds() < 0 {
		return fmt.Errorf("freeze rate cannot be lower than zero")
	}
	return nil
}
