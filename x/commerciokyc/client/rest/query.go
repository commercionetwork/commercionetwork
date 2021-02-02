package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

const (
	restParamAddress = "user"
	restParamTsp     = "tsp"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/commerciokyc/invites", getInvitesHandler(cliCtx)).Methods(http.MethodGet)

	r.HandleFunc("/commerciokyc/tsps", getTrustedServiceProvidersHandler(cliCtx)).Methods(http.MethodGet)

	r.HandleFunc("/commerciokyc/memberships", getMemberships(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(
		fmt.Sprintf("/commerciokyc/membership/{%s}", restParamAddress),
		getMembershipForAddr(cliCtx),
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/commerciokyc/funds",
		getGetPoolFunds(cliCtx)).
		Methods(http.MethodGet)

	r.HandleFunc(
		fmt.Sprintf("/commerciokyc/sold/{%s}", restParamTsp),
		getSoldForTsp(cliCtx),
	).Methods(http.MethodGet)

}

// ---------------------------------
// --- Commerciokyc
// ---------------------------------

// ----------------------------------
// --- Membership
// ----------------------------------

// @Summary Get Membership for given address
// @Description This endpoint returns the Membership
// @ID getMembershipForAddr
// @Produce json
// @Param address path string true "Address of the user"
// @Success 200 {object} x.JSONResult{result=types.Membership}
// @Failure 404
// @Router /commerciokyc/membership/{address} [get]
// @Tags x/commerciokyc
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

// @Summary Get All Memberships
// @Description This endpoint returns all the Memberships
// @ID getMemberships
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.Membership}
// @Failure 404
// @Router /commerciokyc/memberships [get]
// @Tags x/commerciokyc
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

// @Summary Get All Invites
// @Description This endpoint returns all the Invites
// @ID getInvitesHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.Invites}
// @Failure 404
// @Router /commerciokyc/invites [get]
// @Tags x/commerciokyc
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

// @Summary Get All Trusted Service Providers
// @Description This endpoint returns all the Trusted Service Providers
// @ID getTrustedServiceProvidersHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.AccAddress}
// @Failure 404
// @Router /commerciokyc/tsps [get]
// @Tags x/commerciokyc
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

// @Summary Get All Memberships sold by Trusted Service Provider
// @Description This endpoint returns all Memberships sold by a specific Trusted Service Provider
// @ID getSoldForTsp
// @Produce json
// @Param did path string true "Address of the tsp which to read the sold memberhip"
// @Success 200 {object} x.JSONResult{result=[]types.Membership}
// @Failure 404
// @Router /commerciokyc/sold/{address} [get]
// @Tags x/commerciokyc
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
// --- Funds
// ---------------------------------

// @Summary Get All Current pool funds
// @Description This endpoint returns current pool funds for accreditation block reward
// @ID getGetPoolFunds
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.Coin}
// @Failure 404
// @Router /commerciokyc/funds [get]
// @Tags x/commerciokyc
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
