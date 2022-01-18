package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyCheckMembershipsEpochIdentifier = []byte("CheckMembershipsEpochIdentifier")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(checkMembershipsEpochIdentifier string) Params {
	return Params{
		CheckMembershipsEpochIdentifier: checkMembershipsEpochIdentifier,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		CheckMembershipsEpochIdentifier: EpochDay,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateCheckMembershipsEpochIdentifier(p.CheckMembershipsEpochIdentifier); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCheckMembershipsEpochIdentifier, &p.CheckMembershipsEpochIdentifier, validateCheckMembershipsEpochIdentifier),
	}
}

func validateCheckMembershipsEpochIdentifier(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("empty check memberships epoch identifier: %+v", i)
	}

	return nil
}
