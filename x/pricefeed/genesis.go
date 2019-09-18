package mint

import (
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Oracles ctypes.Addresses
	//todo think about this
	prices sdk.Int
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{Oracles: ctypes.Addresses{}}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, genState GenesisState) {
	for _, oracle := range genState.Oracles {
		keeper.AddOracle(ctx, oracle)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {

	oracles := keeper.GetOracles(ctx)

	return GenesisState{
		Oracles: oracles,
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state GenesisState) error {
	return nil
}
