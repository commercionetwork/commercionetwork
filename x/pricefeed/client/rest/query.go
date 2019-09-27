package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	restTokenName = "tokenName"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/pricefeed/prices"),
		getCurrentPricesHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("/pricefeed/prices/{%s}", restTokenName),
		getCurrentPriceHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("/pricefeed/oracles"),
		getOraclesHandler(cliCtx)).Methods("GET")
}

func getCurrentPriceHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenName := vars[restTokenName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCurrentPrice, tokenName)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getCurrentPricesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetCurrentPrices)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getOraclesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetOracles)

		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
