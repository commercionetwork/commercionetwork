package types

import (
	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func NewParams(distrEpochIdentifier string, earnRate math.LegacyDec) Params {
	return Params{
		DistrEpochIdentifier: distrEpochIdentifier,
		EarnRate:             earnRate,
	}
}

// Parameter store keys
var (
	KeyDistrEpochIdentifier = []byte("DistrEpochIdentifier")
	KeyEarnRate             = []byte("EarnRate")
)

// ParamTable for vbr module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// default vbr module params
func DefaultParams() Params {
	return NewParams(EpochDay, math.LegacyNewDecWithPrec(5, 1))
}

// Params validation
func (p Params) Validate() error {
	if err := ValidateDistrEpochIdentifier(p.DistrEpochIdentifier); err != nil {
		return err
	}
	if err := ValidateEarnRate(p.EarnRate); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistrEpochIdentifier, &p.DistrEpochIdentifier, validateDistrEpochIdentifierParamSetPairs),
		paramtypes.NewParamSetPair(KeyEarnRate, &p.EarnRate, validateEarnRateParamSetPairs),
	}
}

func validateDistrEpochIdentifierParamSetPairs(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := ValidateDistrEpochIdentifier(v); err != nil {
		return err
	}

	return nil
}

func ValidateDistrEpochIdentifier(i string) error {
	switch i {
	case EpochDay, EpochWeek, EpochHour, EpochMinute:
		return nil
	}

	return fmt.Errorf("invalid distr epoch identifier: %s", i)
}

func validateEarnRateParamSetPairs(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := ValidateEarnRate(v); err != nil {
		return err
	}

	return nil
}

func ValidateEarnRate(e math.LegacyDec) error {
	if e.IsNegative() {
		return fmt.Errorf("earn rate cannot be negative: %+v", e)
	}

	return nil
}
