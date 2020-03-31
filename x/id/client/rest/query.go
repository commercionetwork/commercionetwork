package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	identityParam = "identity"
	proofParam    = "proof"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, querierRoute string) {
	r.HandleFunc(fmt.Sprintf(
		"/identities/{%s}", identityParam),
		resolveIdentityHandler(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(fmt.Sprintf(
		"/depositRequests/{%s}", proofParam),
		resolveDepositRequestHandler(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(fmt.Sprintf(
		"/powerUpRequest/{%s}", proofParam),
		resolvePowerUpRequestHandler(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(fmt.Sprintf(
		"/approvedPowerUpRequests"),
		resolveApprovedPowerUpRequests(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(fmt.Sprintf(
		"/rejectedPowerUpRequests"),
		resolveRejectedPowerUpRequests(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(fmt.Sprintf(
		"/pendingPowerUpRequests"),
		resolvePendingPowerUpRequests(cliCtx, querierRoute)).
		Methods("GET")
}

func resolveIdentityHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[identityParam]

		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveDid, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func resolveDepositRequestHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[proofParam]

		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveDepositRequest, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func resolvePowerUpRequestHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[proofParam]

		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolvePowerUpRequest, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func resolveApprovedPowerUpRequests(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryGetApprovedPowerUpRequest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func resolveRejectedPowerUpRequests(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryGetRejectedPowerUpRequest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
func resolvePendingPowerUpRequests(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryGetPendingPowerUpRequest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
