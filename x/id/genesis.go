package id

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/id/keeper"
	"github.com/commercionetwork/commercionetwork/x/id/types"
)

// GenesisState - id genesis state
type GenesisState struct {
	DidDocuments    []types.DidDocument       `json:"did_documents"`
	PowerUpRequests []types.DidPowerUpRequest `json:"power_up_requests"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets ids information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	for _, didDocument := range data.DidDocuments {
		if err := keeper.SaveDidDocument(ctx, didDocument); err != nil {
			panic(err)
		}
	}

	for _, powerUp := range data.PowerUpRequests {
		if err := keeper.StorePowerUpRequest(ctx, powerUp); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	identities := keeper.GetDidDocuments(ctx)
	requests := keeper.GetPowerUpRequests(ctx)

	return GenesisState{
		DidDocuments:    identities,
		PowerUpRequests: requests,
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
