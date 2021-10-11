package documents

import (
	"github.com/commercionetwork/commercionetwork/x/documents/keeper"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets documents information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	for _, doc := range data.Documents {
		if err := keeper.SaveDocument(ctx, *doc); err != nil {
			panic(err)
		}
	}

	for _, receipt := range data.Receipts {
		if err := keeper.SaveReceipt(ctx, *receipt); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Documents: exportDocuments(ctx, keeper),
		Receipts:  exportReceipts(ctx, keeper),
	}
}
