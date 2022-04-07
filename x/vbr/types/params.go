package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyDistrEpochIdentifier = []byte("DistrEpochIdentifier")
	KeyEarnRate             = []byte("EarnRate")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(distrEpochIdentifier string, earnRate sdk.Dec) Params {
	return Params{
		DistrEpochIdentifier: distrEpochIdentifier,
		EarnRate:             earnRate,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		DistrEpochIdentifier: EpochDay,
		EarnRate:             sdk.NewDecWithPrec(5, 1),
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateDistrEpochIdentifier(p.DistrEpochIdentifier); err != nil {
		return err
	}
	if err := validateEarnRate(p.EarnRate); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistrEpochIdentifier, &p.DistrEpochIdentifier, validateDistrEpochIdentifier),
		paramtypes.NewParamSetPair(KeyEarnRate, &p.EarnRate, validateEarnRate),
	}
}

func validateDistrEpochIdentifier(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("empty distribution epoch identifier: %+v", i)
	}

	return nil
}

func validateEarnRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("earn rate cannot be negative: %+v", i)
	}

	return nil
}
