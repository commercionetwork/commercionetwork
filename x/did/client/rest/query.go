package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	restTypes "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	identityParam = "identity"
	proofParam    = "proof"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router, querierRoute string) {
	r.HandleFunc(fmt.Sprintf(
		"/identities/{%s}", identityParam),
		resolveIdentityHandler(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(fmt.Sprintf(
		"/powerUpRequest/{%s}", proofParam),
		resolvePowerUpRequestHandler(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(
		"/approvedPowerUpRequests",
		resolveApprovedPowerUpRequests(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(
		"/rejectedPowerUpRequests",
		resolveRejectedPowerUpRequests(cliCtx, querierRoute)).
		Methods("GET")

	r.HandleFunc(
		"/pendingPowerUpRequests",
		resolvePendingPowerUpRequests(cliCtx, querierRoute)).
		Methods("GET")
}

// @Summary Get a user Did Document
// @Description This endpoint returns a user Did Document, along with the height at which the resource was queried at
// @ID id_resolveIdentityHandler
// @Produce json
// @Param did path string true "Address of the user for which to read the Did Document"
// @Success 200 {object} x.JSONResult{result=keeper.ResolveIdentityResponse}
// @Failure 404
// @Router /identities/{did} [get]
// @Tags x/id
func resolveIdentityHandler(cliCtx client.Context, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[identityParam]

		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveDid, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get a user Did power up request
// @Description This endpoint returns a user Did power up request, along with the height at which the resource was queried at
// @ID id_resolvePowerUpRequest
// @Produce json
// @Param id path string true "Request id"
// @Success 200 {object} x.JSONResult{result=types.DidPowerUpRequest}
// @Failure 404
// @Router /powerUpRequest/{id} [get]
// @Tags x/id
func resolvePowerUpRequestHandler(cliCtx client.Context, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[proofParam]

		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolvePowerUpRequest, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get the user Did power up approved requests
// @Description This endpoint returns the user Did power up approved requests, along with the height at which the resource was queried at
// @ID id_resolveApprovedPowerUpRequests
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.DidPowerUpRequest}
// @Failure 404
// @Router /approvedPowerUpRequests [get]
// @Tags x/id
func resolveApprovedPowerUpRequests(cliCtx client.Context, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryGetApprovedPowerUpRequest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get the user Did power up rejected power up requests
// @Description This endpoint returns the user Did power up rejected requests, along with the height at which the resource was queried at
// @ID id_resolveRejectedPowerUpRequests
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.DidPowerUpRequest}
// @Failure 404
// @Router /rejectedPowerUpRequests [get]
// @Tags x/id
func resolveRejectedPowerUpRequests(cliCtx client.Context, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryGetRejectedPowerUpRequest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get the user Did power up pending requests
// @Description This endpoint returns the user Did power up pending requests, along with the height at which the resource was queried at
// @ID id_resolvePendingPowerUpRequests
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.DidPowerUpRequest}
// @Failure 404
// @Router /pendingPowerUpRequests [get]
// @Tags x/id
func resolvePendingPowerUpRequests(cliCtx client.Context, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryGetPendingPowerUpRequest)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
