package did

import (
	"github.com/commercionetwork/commercionetwork/x/did/keeper"
	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, elem := range genState.DidDocuments {
		k.UpdateDidDocument(ctx, *elem)
	}

	// this line is used by starport scaffolding # genesis/module/init

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	didDocuments := k.GetAllDidDocuments(ctx)
	for _, elem := range didDocuments {
		elem := elem
		genesis.DidDocuments = append(genesis.DidDocuments, &elem)
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
