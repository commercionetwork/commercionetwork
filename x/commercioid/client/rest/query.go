package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf(
		"/%s/identities/{%s}", storeName, restName),
		resolveIdentityHandler(cliCtx, storeName)).
		Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/%s/identities/{%s}/connections", storeName, restName),
		getConnectionsHandler(cliCtx, storeName)).
		Methods("GET")
}

// ----------------------------------
// --- ResolveIdentity
// ----------------------------------

func resolveIdentityHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		route := fmt.Sprintf("custom/%s/identities/%s", storeName, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ----------------------------------
// --- GetConnections
// ----------------------------------

func getConnectionsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		route := fmt.Sprintf("custom/%s/connections/%s", storeName, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
