package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateRewardRate(rate sdk.Dec) error {
	if rate.IsNil() {
		return fmt.Errorf("reward rate must be not nil")
	}
	if !rate.IsPositive() {
		return fmt.Errorf("reward rate must be positive: %s", rate)
	}
	return nil
}
