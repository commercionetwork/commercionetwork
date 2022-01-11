package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultConversionRate = sdk.NewDec(1)
)

const (
	DefaultCreditsDenom               = "uccc"
	DefaultFreezePeriod time.Duration = time.Hour * 24 * 7 * 3
)

// Parameter store keys
var (
	KeyConversionRate = []byte("CollateralRate")
	KeyFreezePeriod   = []byte("FreezePeriod")
)

func (p *Params) ValidateBasic() error {

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
