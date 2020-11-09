package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

const (
	restParamAddress = "user"
	restParamTsp     = "tsp"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/invites", getInvitesHandler(cliCtx)).Methods(http.MethodGet)

	r.HandleFunc("/tsps", getTrustedServiceProvidersHandler(cliCtx)).Methods(http.MethodGet)

	r.HandleFunc("/memberships", getMemberships(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(
		fmt.Sprintf("/membership/{%s}", restParamAddress),
		getMembershipForAddr(cliCtx),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/funds",
		getGetPoolFunds(cliCtx)).
		Methods(http.MethodGet)

	r.HandleFunc(
		fmt.Sprintf("/sold/{%s}", restParamTsp),
		getSoldForTsp(cliCtx),
	).Methods(http.MethodGet)

}

// ----------------------------------
// --- Memberships
// ----------------------------------

func getMembershipForAddr(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[restParamAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetMembership, address)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getMemberships(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetMemberships)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ----------------------------------
// --- Invites
// ----------------------------------

func getInvitesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetInvites)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ---------------------------------
// --- Trusted Service Providers
// ---------------------------------

func getTrustedServiceProvidersHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetTrustedServiceProviders)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getSoldForTsp(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tsp := vars[restParamTsp]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetTspMemberships, tsp)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ---------------------------------
// --- Abr
// ---------------------------------

func getGetPoolFunds(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetPoolFunds)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
