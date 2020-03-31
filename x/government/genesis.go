package government

import (
	"errors"

	"github.com/commercionetwork/commercionetwork/x/government/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	GovernmentAddress sdk.AccAddress `json:"government_address"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	err := keeper.SetGovernmentAddress(ctx, data.GovernmentAddress)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{
		GovernmentAddress: keeper.GetGovernmentAddress(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.GovernmentAddress.Empty() {
		return errors.New("government address cannot be empty. Use the set-genesis-government-address command to set one")
	}
	return nil
}
