package memberships

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState state at genesis
type GenesisState struct {
	Minters Minters `json:"minters"` // List of users allowed to sign a membership giving tx
}

// InitGenesis sets membership information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, genState GenesisState) {
	for _, minter := range genState.Minters {
		keeper.AddTrustedMinter(ctx, minter)
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minters: Minters{},
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Minters: keeper.GetTrustedMinters(ctx),
	}
}
