package docs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/**
THIS FILE IS USELESS BUT IT HELPS ME REMEMBER THIS FILE STRUCTURE IN ORDER TO USE IT IN FUTURE MODULES
*/

type GenesisState struct {
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return DefaultGenesisState()
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
