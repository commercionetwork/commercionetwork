package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// UpdateConversionRate stores the conversion rate.
func (k Keeper) UpdateConversionRate(ctx sdk.Context, rate sdk.DecProto) error {
	if err := types.ValidateConversionRate(rate.Dec); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	// TODO work around the Marshal() method of Dec object
	store.Set([]byte(types.CollateralRateKey), []byte(rate.String()))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSetConversionRate,
		sdk.NewAttribute("rate", rate.String()),
	))

	return nil
}

// GetConversionRate retrieves the conversion rate.
// TODO CONTROL
func (k Keeper) GetConversionRate(ctx sdk.Context) sdk.DecProto {
	store := ctx.KVStore(k.storeKey)
	rate, _ := sdk.NewDecFromStr(string(store.Get([]byte(types.CollateralRateKey))))
	var returnRate sdk.DecProto
	returnRate.Dec = rate
	return returnRate
}

// UpdateFreezePeriod stores the freeze period in seconds.
func (k Keeper) UpdateFreezePeriod(ctx sdk.Context, freezePeriod time.Duration) error {
	if err := types.ValidateFreezePeriod(freezePeriod); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	strFreezePeriod := freezePeriod.String()
	store.Set([]byte(types.FreezePeriodKey), []byte(strFreezePeriod)) // TODO CHECK IF THIS IS CORRECT
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSetFreezePeriod,
		sdk.NewAttribute("freeze_period", freezePeriod.String()),
	))

	return nil
}

// GetFreezePeriod retrieves the freeze period.
func (k Keeper) GetFreezePeriod(ctx sdk.Context) time.Duration {
	store := ctx.KVStore(k.storeKey)
	var freezePeriod string
	//k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.FreezePeriodKey)), &freezePeriod)
	freezePeriod = string(store.Get([]byte(types.FreezePeriodKey)))
	freezePeriodDuration, _ := time.ParseDuration(freezePeriod) // Catch the error
	return freezePeriodDuration
}
