package custombank

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

// GenesisState - bank genesis state
type GenesisState struct {
	bank.GenesisState
	BlockedAccounts ctypes.Addresses `json:"blocked_accounts"`
}

// MarshalJSON implements the json.Marshaler interface. We do this because Amino
// does not respect the JSON stdlib embedding semantics.
func (g GenesisState) MarshalJSON() ([]byte, error) {
	type state GenesisState
	return json.Marshal(state(g))
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		GenesisState:    bank.DefaultGenesisState(),
		BlockedAccounts: ctypes.Addresses{},
	}
}

// InitGenesis sets bank information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, genState GenesisState) {
	bank.InitGenesis(ctx, keeper, genState.GenesisState)
	keeper.SetBlockedAddresses(ctx, genState.BlockedAccounts)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		GenesisState:    bank.ExportGenesis(ctx, keeper),
		BlockedAccounts: keeper.GetBlockedAddresses(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
