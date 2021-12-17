package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyDistrEpochIdentifier = []byte("DistrEpochIdentifier")
	KeyVbrEarnRate = []byte("VbrEarnRate")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(distrEpochIdentifier string, vbrEarnRate sdk.Dec) Params {
	return Params{
		DistrEpochIdentifier: distrEpochIdentifier,
		VbrEarnRate: vbrEarnRate,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		DistrEpochIdentifier: /*EpochMinute*/EpochDay,
		VbrEarnRate: sdk.NewDec(int64(50)),
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateDistrEpochIdentifier(p.DistrEpochIdentifier); err != nil {
		return err
	}
	if err := validateVbrEarnRate(p.VbrEarnRate); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistrEpochIdentifier, &p.DistrEpochIdentifier, validateDistrEpochIdentifier),
		paramtypes.NewParamSetPair(KeyVbrEarnRate, &p.VbrEarnRate, validateVbrEarnRate),
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

func validateVbrEarnRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() || v.GT(sdk.NewDec(1)) {
		return fmt.Errorf("invalid vbr earn rate(must be between 0 and 1): %+v", i)
	}

	return nil
}