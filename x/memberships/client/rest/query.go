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
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/invites", getInvitesHandler(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("/invites/{%s}", restParamAddress),
		getInvitesForUserHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc("/tsps", getTrustedServiceProvidersHandler(cliCtx)).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/membership/{%s}", restParamAddress),
		getMembershipForAddr(cliCtx),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/accreditations-funds",
		getGetPoolFunds(cliCtx)).
		Methods("GET")
}

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

func getInvitesForUserHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[restParamAddress]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetInvites, address)
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
