package app

import (
	"context"
	
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// App implements the common methods for a Cosmos SDK-based application
// specific blockchain.
type CosmosApp interface {
	// The assigned name of the app.
	Name() string

	// The application types codec.
	// NOTE: This shoult be sealed before being returned.
	LegacyAmino() *codec.LegacyAmino

	// Application updates every begin block.
	BeginBlock(context.Context) error

	// Application updates every end block.
	EndBlock(context.Context) error

	// Application update at chain (i.e app) initialization.
	InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain

	// Loads the app at a given height.
	LoadHeight(height int64) error

	// Exports the state of the application for a genesis file.
	ExportAppStateAndValidators(
		forZeroHeight bool, jailAllowedAddrs []string,
	) (types.ExportedApp, error)

	// All the registered module account addreses.
	ModuleAccountAddrs() map[string]bool
}
