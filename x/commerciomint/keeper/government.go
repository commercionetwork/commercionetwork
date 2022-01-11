package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// UpdateConversionRate stores the conversion rate.
func (k Keeper) UpdateConversionRate(ctx sdk.Context, conversionRate sdk.Dec) error {
	if err := types.ValidateConversionRate(conversionRate); err != nil {
		return err
	}

	p := k.GetParams(ctx)
	p.ConversionRate = conversionRate

	k.UpdateParams(ctx, p)

	// store := ctx.KVStore(k.storeKey)
	// // TODO work around the Marshal() method of Dec object
	// store.Set([]byte(types.CollateralRateKey), []byte(conversionRate.String()))

	// ctx.EventManager().EmitEvent(sdk.NewEvent(
	// 	eventSetConversionRate,
	// 	sdk.NewAttribute("rate", conversionRate.String()),
	// ))

	return nil
}

// GetConversionRate retrieves the conversion rate.
// TODO CONTROL
func (k Keeper) GetConversionRate(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).ConversionRate
	// store := ctx.KVStore(k.storeKey)
	// rate, _ := sdk.NewDecFromStr(string(store.Get([]byte(types.CollateralRateKey))))
	// return rate
}

// UpdateFreezePeriod stores the freeze period in seconds.
func (k Keeper) UpdateFreezePeriod(ctx sdk.Context, freezePeriod time.Duration) error {
	if err := types.ValidateFreezePeriod(freezePeriod); err != nil {
		return err
	}

	p := k.GetParams(ctx)
	p.FreezePeriod = freezePeriod

	k.UpdateParams(ctx, p)

	// store := ctx.KVStore(k.storeKey)
	// strFreezePeriod := freezePeriod.String()
	// store.Set([]byte(types.FreezePeriodKey), []byte(strFreezePeriod)) // TODO CHECK IF THIS IS CORRECT
	// ctx.EventManager().EmitEvent(sdk.NewEvent(
	// 	eventSetFreezePeriod,
	// 	sdk.NewAttribute("freeze_period", freezePeriod.String()),
	// ))

	return nil
}

// GetFreezePeriod retrieves the freeze period.
func (k Keeper) GetFreezePeriod(ctx sdk.Context) time.Duration {
	return k.GetParams(ctx).FreezePeriod

	// store := ctx.KVStore(k.storeKey)
	// freezePeriod := string(store.Get([]byte(types.FreezePeriodKey)))
	// freezePeriodDuration, _ := time.ParseDuration(freezePeriod) // Catch the error
	// return freezePeriodDuration
}
