package types

/*
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultRewardTax      sdk.NewDecWithPrec(1, 10)
)


// Parameter keys
var (
	RewardTax = []byte("RewardTax")
)

var _ paramtypes.ParamSet = &Params{}

// NewParams creates a new Params object
func NewParams(
	rewardTax sdk.Dec,
) Params {
	return Params{
		RewardTax: rewardTax,
	}
}

// ParamKeyTable for vbr module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of vbr module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(RewardTax, &p.RewardTax, validateRewardTax),
	}
}


// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		RewardTax:      DefaultRewardTax,
	}
}

// validateRewardTax controls if RewardTax is valid.
func validateRewardTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid reward tax: %s", v)
	}
	return nil
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateRewardTax(p.RewardTax); err != nil {
		return err
	}

	return nil
}
*/
