package v3_0_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	v220vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v2.2.0"
)

// Migrate convert v2.2.0 chain to v3.0.0
func Migrate(genVbr v220vbr.GenesisState) *types.GenesisState {
	return migrateVbr(genVbr)
}

func migrateVbr(genVbr v220vbr.GenesisState) *types.GenesisState {
	return &types.GenesisState{
		PoolAmount:        genVbr.PoolAmount,
		RewardRate:        sdk.NewDecWithPrec(1, 3),
		AutomaticWithdraw: true,
	}
}
