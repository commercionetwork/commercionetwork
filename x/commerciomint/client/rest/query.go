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
		fmt.Sprintf("/commerciomint/{%s}/etp", restuser),
		getEtpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("/commerciomint/{%s}/etpsOwner", restuser),
		getEtpsByOwnerHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/etps", getAllEtpsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciomint/params", getParamsHandler(cliCtx)).Methods("GET")
}

// ----------------------------------
// --- Commerciomint
// ----------------------------------

func getEtpsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[restuser]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, id, types.QueryGetEtpRest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getEtpsByOwnerHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerAddr := vars[restuser]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, ownerAddr, types.QueryGetEtpsByOwnerRest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
func getAllEtpsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetallEtpsRest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
func getParamsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetParamsRest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
