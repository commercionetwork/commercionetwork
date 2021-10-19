package government

import (
	"github.com/commercionetwork/commercionetwork/x/government/keeper"
	"github.com/commercionetwork/commercionetwork/x/government/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {

	// this line is used by starport scaffolding # genesis/module/init
	govAddr, err := sdk.AccAddressFromBech32(genState.GovernmentAddress)
	if err != nil {
		panic(err)

	}

	errSetGov := k.SetGovernmentAddress(ctx, govAddr)
	if errSetGov != nil {
		panic(errSetGov)
	}

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	genesis.GovernmentAddress = k.GetGovernmentAddress(ctx).String()
	return genesis
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ types.GenesisState) error {
	return nil
}
