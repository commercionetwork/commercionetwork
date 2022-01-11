package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultConversionRate               = sdk.NewDec(1)
	DefaultFreezePeriod   time.Duration = time.Hour * 24 * 7 * 3
	KeyConversionRate                   = []byte("ConversionRate")
	KeyFreezePeriod                     = []byte("FreezePeriod")
)

const (
	DefaultCreditsDenom = "uccc"
)

// ParamTable for commerciomint module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(conversionRate sdk.Dec, freezePeriod *time.Duration) Params {
	return Params{
		ConversionRate: conversionRate,
		FreezePeriod:   freezePeriod,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		ConversionRate: DefaultConversionRate,
		FreezePeriod:   &DefaultFreezePeriod,
	}
}

func (p *Params) Validate() error {

	if err := ValidateConversionRate(p.ConversionRate); err != nil {
		return err
	}

	if err := ValidateFreezePeriod(*p.FreezePeriod); err != nil {
		return err
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
	v, ok := i.(sdk.Dec)
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
