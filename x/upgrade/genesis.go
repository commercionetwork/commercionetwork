package upgrade

import (
	"github.com/commercionetwork/commercionetwork/x/upgrade/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis validates the upgrade genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {

}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data GenesisState) {
	return NewGenesisState()
}
