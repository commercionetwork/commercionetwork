package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		// this line is used by starport scaffolding # ibc/keeper/attribute
		governmentKeeper government.Keeper
		UpgradeKeeper    upgrade.Keeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	// this line is used by starport scaffolding # ibc/keeper/parameter
	governmentKeeper government.Keeper,
	upgradeKeeper upgrade.Keeper,

) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		// this line is used by starport scaffolding # ibc/keeper/return
		governmentKeeper: governmentKeeper,
		UpgradeKeeper:    upgradeKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ScheduleUpgrade schedules an upgrade based on the specified plan.
// Only the government can schedule an upgrade.
// If there is another Plan already scheduled, it will overwrite it
// (implicitly cancelling the current plan)
func (k Keeper) ScheduleUpgradeGov(ctx sdk.Context, address sdk.AccAddress, plan upgradeTypes.Plan) error {
	if !address.Equals(k.governmentKeeper.GetGovernmentAddress(ctx)) {
		return fmt.Errorf("only the government address can schedule an upgrade")
	}

	return k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
}

// DeleteUpgrade clears any scheduled upgrade.
// Only the government can clear any scheduled upgrade.
func (k Keeper) DeleteUpgradeGov(ctx sdk.Context, address sdk.AccAddress) error {
	if !address.Equals(k.governmentKeeper.GetGovernmentAddress(ctx)) {
		return fmt.Errorf("only the government address can delete any upgrade")
	}

	k.UpgradeKeeper.ClearUpgradePlan(ctx)
	return nil
}

// GetUpgradePlan returns the currently scheduled Plan if any, setting havePlan to true if there is a scheduled
// upgrade or false if there is none
func (k Keeper) GetUpgradePlan(ctx sdk.Context) (plan upgradeTypes.Plan, havePlan bool) {
	return k.UpgradeKeeper.GetUpgradePlan(ctx)
}