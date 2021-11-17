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
	restuser = "user"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/commerciomint/etp/{%s}", restuser),
		getEtpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("/commerciomint/owner/{%s}", restuser), 
		getEtpsByOwnerHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/etps", getAllEtpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/conversion_rate", getConversionRateHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/freeze_period", getFreezePeriodHandler(cliCtx)).Methods("GET")
}

// ----------------------------------
// --- Commerciomint
// ----------------------------------

func getEtpsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[restuser]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetEtp, id)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getEtpsByOwnerHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restuser]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetEtpsByOwner, ownerAddr)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
func getAllEtpsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetallEtps)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

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
