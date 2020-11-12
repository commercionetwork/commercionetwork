package keeper

import (
	"fmt"

	government "github.com/commercionetwork/commercionetwork/x/government/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

type Keeper struct {
	cdc              *codec.Codec
	governmentKeeper government.Keeper
	UpgradeKeeper    upgrade.Keeper
}

func NewKeeper(cdc *codec.Codec, governmentKeeper government.Keeper, upgradeKeeper upgrade.Keeper) Keeper {
	return Keeper{
		cdc:              cdc,
		governmentKeeper: governmentKeeper,
		UpgradeKeeper:    upgradeKeeper,
	}
}

// ScheduleUpgrade schedules an upgrade based on the specified plan.
// Only the government can schedule an upgrade.
// If there is another Plan already scheduled, it will overwrite it
// (implicitly cancelling the current plan)
func (k Keeper) ScheduleUpgrade(ctx sdk.Context, address sdk.AccAddress, plan upgrade.Plan) error {
	if !address.Equals(k.governmentKeeper.GetGovernmentAddress(ctx)) {
		return fmt.Errorf("only the government address can schedule an upgrade")
	}

	return k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
}

// GetUpgradePlan returns the currently scheduled Plan if any, setting havePlan to true if there is a scheduled
// upgrade or false if there is none
func (k Keeper) GetUpgradePlan(ctx sdk.Context) (plan upgrade.Plan, havePlan bool) {
	return k.UpgradeKeeper.GetUpgradePlan(ctx)
}

// DeleteUpgrade clears any scheduled upgrade.
// Only the government can clear any scheduled upgrade.
func (k Keeper) DeleteUpgrade(ctx sdk.Context, address sdk.AccAddress) error {
	if !address.Equals(k.governmentKeeper.GetGovernmentAddress(ctx)) {
		return fmt.Errorf("only the government address can delete any upgrade")
	}

	k.UpgradeKeeper.ClearUpgradePlan(ctx)
	return nil
}
