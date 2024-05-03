package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	//restTypes "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	restParamAddress = "user"
	restParamTsp     = "tsp"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/commerciokyc/funds", getPoolFunds(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/commercionetwork/commerciokyc/invite/{%s}", restParamAddress),
		getInviteHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciokyc/invites", getInvitesHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciokyc/tsps", getTrustedServiceProvidersHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/commerciokyc/memberships", getMemberships(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/commerciokyc/membership/{%s}", restParamAddress),
		getMembershipForAddr(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/commerciokyc/sold/{%s}", restParamTsp),
		getSoldForTsp(cliCtx)).Methods("GET")

}

func getPoolFunds(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetPoolFunds)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getMembershipForAddr(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[restParamAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetMembership, address)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getMemberships(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetMemberships)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// ----------------------------------
// --- Invites
// ----------------------------------

func getInvitesHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetInvites)

		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getInviteHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetInvite)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// ---------------------------------
// --- Trusted Service Providers
// ---------------------------------

func getTrustedServiceProvidersHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetTrustedServiceProviders)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getSoldForTsp(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tsp := vars[restParamTsp]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetTspMemberships, tsp)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
