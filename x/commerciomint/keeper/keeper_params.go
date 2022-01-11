package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the total set params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// UpdateParams sets the total set of params
func (k Keeper) UpdateParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	k.paramSpace.SetParamSet(ctx, &params)
	return nil
}
