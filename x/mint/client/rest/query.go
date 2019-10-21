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
	restOwnerAddress = "ownerAddress"
	restTimestamp    = "timestamp"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/mint/cdps/{%s}", restOwnerAddress),
		getCdpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("mint/cdps/{%s}/{%s}", restOwnerAddress, restTimestamp),
		getCdpHandler(cliCtx)).Methods("GET")
}

func getCdpHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restOwnerAddress]
		timestamp := vars[restTimestamp]

		route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryGetCdp, ownerAddr, timestamp)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getCdpsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restOwnerAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCdps, ownerAddr)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
