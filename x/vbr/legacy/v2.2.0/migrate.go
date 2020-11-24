package v2_2_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v212vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v2.1.2"
)

// Migrate convert v2.1.2 chain to v2.2.0
func Migrate(genVbr v212vbr.GenesisState) GenesisState {
	return migrateVbr(genVbr)
}

func migrateVbr(genVbr v212vbr.GenesisState) GenesisState {
	return GenesisState{
		PoolAmount:        genVbr.PoolAmount,
		RewardRate:        sdk.NewDecWithPrec(1, 3),
		AutomaticWithdraw: true,
	}
}
