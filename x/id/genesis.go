package id

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// GenesisState - id genesis state
type GenesisState struct {
	DidDocuments           []DidDocument       `json:"did_documents"`
	DepositRequests        []DidDepositRequest `json:"deposit_requests"`
	PowerUpRequests        []DidPowerUpRequest `json:"power_up_requests"`
	DepositPool            sdk.Coins           `json:"deposit_pool"`
	HandledPowerUpRequests []string            `json:"handled_power_up_requests"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets ids information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper supply.Keeper, data GenesisState) {
	moduleAcc := keeper.GetIdModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ModuleName))
	}

	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(data.DepositPool); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	for _, didDocument := range data.DidDocuments {
		if err := keeper.SaveDidDocument(ctx, didDocument); err != nil {
			panic(err)
		}
	}

	for _, deposit := range data.DepositRequests {
		if err := keeper.StoreDidDepositRequest(ctx, deposit); err != nil {
			panic(err)
		}
	}

	for _, powerUp := range data.PowerUpRequests {
		if err := keeper.StorePowerUpRequest(ctx, powerUp); err != nil {
			panic(err)
		}
	}

	if data.HandledPowerUpRequests != nil {
		keeper.SetHandledPowerUpRequestsReferences(ctx, data.HandledPowerUpRequests)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	identities, err := keeper.GetDidDocuments(ctx)
	if err != nil {
		panic(err)
	}

	return GenesisState{
		DidDocuments:           identities,
		DepositRequests:        keeper.GetDepositRequests(ctx),
		PowerUpRequests:        keeper.GetPowerUpRequests(ctx),
		DepositPool:            keeper.GetPoolAmount(ctx),
		HandledPowerUpRequests: keeper.GetHandledPowerUpRequestsReferences(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
