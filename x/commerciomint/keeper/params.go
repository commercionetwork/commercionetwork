package keeper

import (
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateParams sets the module params.
func (k Keeper) UpdateParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	k.paramSpace.SetParamSet(ctx, &params)
	return nil
}

// GetParams retrieves the module params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// GetConversionRate retrieves the conversion rate param.
func (k Keeper) GetConversionRate(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).ConversionRate
}

// GetFreezePeriod retrieves the freeze period param.
func (k Keeper) GetFreezePeriod(ctx sdk.Context) time.Duration {
	return k.GetParams(ctx).FreezePeriod
}
