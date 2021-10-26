package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/client"
	restTypes "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	restOwnerAddress = "ownerAddress"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/commerciomint/etps/{%s}", restOwnerAddress),
		getEtpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/conversion_rate", getConversionRateHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/freeze_period", getFreezePeriodHandler(cliCtx)).Methods("GET")
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
func getEtpsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restOwnerAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetEtps, ownerAddr)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get Conversion rate
// @Description This endpoint returns the Conversion rate, along with the height at which the resource was queried at
// @ID getConversionRateHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=types.Dec}
// @Failure 404
// @Router /commerciomint/conversion_rate [get]
// @Tags x/commerciomint
func getConversionRateHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryConversionRateRest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get Freeze period
// @Description This endpoint returns the Freeze period, along with the height at which the resource was queried at
// @ID getFreezePeriodHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=time.Duration}
// @Failure 404
// @Router /commerciomint/freeze_period [get]
// @Tags x/commerciomint
func getFreezePeriodHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFreezePeriodRest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
