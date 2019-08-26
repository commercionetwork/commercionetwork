package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

const (
	restName = "identity"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, moduleName string) {
	registerQueryRoutes(cliCtx, r, moduleName)
	registerTxRoutes(cliCtx, r, moduleName)
}
