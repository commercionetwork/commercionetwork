package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	restOwnerAddress = "restOwnerAddress"
	restTimestamp    = "restTimestamp"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/mint/CDPs/{%s}", restOwnerAddress),
		getCDPsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("mint/CDP/{%s}/{%s}", restOwnerAddress, restTimestamp),
		getCDPHandler(cliCtx)).Methods("GET")
}

func getCDPHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restOwnerAddress]
		timestamp := vars[restTimestamp]

		route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryGetCDP, ownerAddr, timestamp)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getCDPsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restOwnerAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCDPs, ownerAddr)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
