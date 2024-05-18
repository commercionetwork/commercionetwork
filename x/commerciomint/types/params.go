package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultConversionRate               = math.LegacyNewDec(1)
	DefaultFreezePeriod   time.Duration = time.Hour * 24 * 7 * 3
	KeyConversionRate                   = []byte("ConversionRate")
	KeyFreezePeriod                     = []byte("FreezePeriod")
)

// ParamKeyTable for commerciomint module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(conversionRate math.LegacyDec, freezePeriod time.Duration) Params {
	return Params{
		ConversionRate: conversionRate,
		FreezePeriod:   freezePeriod,
	}
}

func (p *Params) Validate() error {

	if err := ValidateConversionRate(p.ConversionRate); err != nil {
		return fmt.Errorf("invalid conversion rate: %e", err)
	}

	if err := ValidateFreezePeriod(p.FreezePeriod); err != nil {
		return fmt.Errorf("invalid freeze period: %e", err)
	}

	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyConversionRate, &p.ConversionRate, validateConversionRateParamSetPairs),
		paramtypes.NewParamSetPair(KeyFreezePeriod, &p.FreezePeriod, validateFreezePeriodParamSetPairs),
	}
}

func validateConversionRateParamSetPairs(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return ValidateConversionRate(v)
}

func validateFreezePeriodParamSetPairs(i interface{}) error {
	fp, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return ValidateFreezePeriod(fp)
}

func ValidateConversionRate(conversionRate math.LegacyDec) error {
	if !conversionRate.IsPositive() {
		return fmt.Errorf("conversion rate must be positive")
	}

	return nil
}

func ValidateFreezePeriod(freezePeriod time.Duration) error {
	if freezePeriod.Seconds() < 0 {
		return fmt.Errorf("freeze rate cannot be lower than zero seconds")
	}
	return nil
}
