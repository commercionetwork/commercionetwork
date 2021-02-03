package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

const (
	restOwnerAddress = "ownerAddress"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/commerciomint/etps/{%s}", restOwnerAddress),
		getEtpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/etps", getConversionRateHandler(cliCtx)).Methods("GET")
}

// ----------------------------------
// --- Commerciomint
// ----------------------------------

// @Summary Get all the Exchange Trade Positions for user
// @Description This endpoint returns the Exchange Trade Position, along with the blocktime at which the resource was queried at
// @ID getEtpsHandler
// @Produce json
// @Param address path string true "Address of the user"
// @Success 200 {object} x.JSONResult{result=[]types.Position}
// @Failure 404
// @Router /commerciomint/etps/{address} [get]
// @Tags x/commerciomint
func getEtpsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restOwnerAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetEtps, ownerAddr)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get Conversion rate
// @Description This endpoint returns the Conversion rate, along with the height at which the resource was queried at
// @ID getConversionRateHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=types.Dec}
// @Failure 404
// @Router /commerciomint/etps [get]
// @Tags x/commerciomint
func getConversionRateHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryConversionRate)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
