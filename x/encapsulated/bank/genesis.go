package custombank

import (
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// GenesisState - bank genesis state
type GenesisState struct {
	bank.GenesisState
	BlockedAccounts ctypes.Addresses `json:"blocked_accounts"`
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
