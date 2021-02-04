package government

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/government/keeper"
)

// GenesisState - government genesis state
type GenesisState struct {
	GovernmentAddress sdk.AccAddress `json:"government_address"`
	TumblerAddress    sdk.AccAddress `json:"tumbler_address"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets documents information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	errSetGov := keeper.SetGovernmentAddress(ctx, data.GovernmentAddress)

	errSetTumb := keeper.SetTumblerAddress(ctx, data.TumblerAddress)

	if errSetGov != nil {
		panic(errSetGov)
	}

	if errSetTumb != nil {
		panic(errSetTumb)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{
		GovernmentAddress: keeper.GetGovernmentAddress(ctx),
		TumblerAddress:    keeper.GetTumblerAddress(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.GovernmentAddress.Empty() {
		return errors.New("government address cannot be empty. Use the set-genesis-government-address command to set one")
	}

	if data.TumblerAddress.Empty() {
		return errors.New("tumbler address cannot be empty. Use the set-genesis-tumbler-address command to set one")
	}
	return nil
}
