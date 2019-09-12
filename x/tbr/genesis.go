package tbr

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct{}

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
