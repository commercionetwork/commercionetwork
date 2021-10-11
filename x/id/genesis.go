package id

import (
	"github.com/commercionetwork/commercionetwork/x/id/keeper"
	"github.com/commercionetwork/commercionetwork/x/id/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, elem := range genState.DocumentList {
		k.AppendId(ctx, *elem)
	}

	// this line is used by starport scaffolding # genesis/module/init

	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	didDocumenttList := k.GetAllDidDocument(ctx)
	for _, elem := range didDocumenttList {
		elem := elem
		genesis.DocumentList = append(genesis.DocumentList, &elem)
	}

	// this line is used by starport scaffolding # genesis/module/export

	genesis.PortId = k.GetPort(ctx)

	return genesis
}
