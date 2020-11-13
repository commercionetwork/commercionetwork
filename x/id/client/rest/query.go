package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commercionetwork/commercionetwork/x/id/types"
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

// @Summary Get a user Did power up request
// @Description This endpoint returns a user Did power up request, along with the height at which the resource was queried at
// @ID id_resolvePowerUpRequest
// @Produce json
// @Success 200 {object} rest.resolvePowerUpRequestResponse
// @Failure 404
// @Router /powerUpRequest/{id} [get]
// @Tags x/id
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

type resolvePowerUpRequestResponse struct {
	Height string                               `json:"height" example:"1234"`
	Result resolveRejectedPowerUpRequestsResult `json:"result"`
}

var _ resolvePowerUpRequestResponse

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

// @Summary Get the user Did power up rejected requests
// @Description This endpoint returns the user Did power up rejected requests, along with the height at which the resource was queried at
// @ID id_resolveRejectedPowerUpRequests
// @Produce json
// @Success 200 {object} rest.resolveRejectedPowerUpRequestsResponse
// @Failure 404
// @Router /rejectedPowerUpRequests [get]
// @Tags x/id
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

type resolveRejectedPowerUpRequestsResponse struct {
	Height string                                 `json:"height" example:"1234"`
	Result []resolveRejectedPowerUpRequestsResult `json:"result"`
}

var _ resolveRejectedPowerUpRequestsResponse

type resolveRejectedPowerUpRequestsResult struct {
	Status    resolveRejectedPowerUpRequestsStatus
	Claimant  string         `json:"claimant" example:"did:com:150jp3tx96frukqg6v870etf02q0cp7em78wu48"`
	Amount    []amountResult `json:"amount"`
	Proof     string         `json:"proof" example:"S5hyg4slMxm9fK8PTNDs8tHmQcBfWXG0vqrNHLXY5K1qUz3QwZYjR9nzJoNDJh18aPsXper7rNBbyZPOm5K//x8Bqm2EJkdnHd7woa5eFqpziGaHxqvgPaLGspH47tnVilARTeF23L2NVHWcEWuo9U5cWg52l1lOixOG+DehT3vC9KjLqg0YqBoL2u0LTLqQMON4UUjC8JwzT/RMs30OYGsWuLc9s48RtJCQJZ+yAg3U6jZn3OokGwWWjYxF9tAsMR48KilHsPigsa9WPnaAyCMSJ05hOqjBxWiSHYiH1nAefFqHtNFXhJF3LRUCJ2xnSHxJC5Ndj4HFzUjyK4aiV1mtRlRcsqmXU80HEk7IzI74HYpW74F8LzXNsh8Pbl7HXoIzEiOHB5XStFnrxkIL3sYAJGH/pGbX3SxeyfoZhY4ikEyqX3OB7Pat2yHh/63XSPThRVpD7g0gy5N2aKBz3vrHCPhe3QQTzWmKlJOcg1FE5ZtSUEHdVQbm1GD9zP6KZDfbekh9+xU0EFczW9JF/we61LTvMF1KoxaBpL46O/J6ROEOQsb03hLEMadBKxZ+XaqAHiQWKu6G5YH2opNTGKcvSyNfDInOvAygUOfzLgTCWp7JOU09hWBKW1ya2yJNJMZ6q9giEAlqS/qqYy4gAqZKjt7nF0siOb3Vz6zEaXdhCcqrfnNN6n/kFXWz24yAucW+/EHt+hsygEVUZQ=="`
	Id        string         `json:"id" example:"d423c645-fd50-4841-8138-192ee8e23dde"`
	Proof_key string         `json:"proof_key" example:"L0QIWxtHeWeUQhmfWqB2n+MZXFqEYctltilM0j69tBd1drUoUSz/vUkaPadQAdKqtQOD43Py7/JZt5IFyx7iDdphzJEX7bqq+B6nC2DQUeISEiXwtDmJYMp20/N23DY2T7L/Z/dzbxRZDWoUhtr9fRPeJL8NHtPqU9YZw2f1tgMk2t/ZMKtBhYzO5BnF8Crmshjw6b6KA3fK+j7YrmF8fVpVFCdz5jd7cprf5RIqwVjt4w1cYZWeKvGLWeGVX3oiCB67EzXZVUCsD03evr90GDY9qGLfUaWJdBkNjByDotLY0OhrKpcZ+O0IZyZv1+YKx7ZDoPAsEJqpqw4M9bGQRg=="`
}

type amountResult struct {
	Denom  string `json:"denom" example:"ucommercio"`
	Amount int    `json:"amount" example:"100000000000000"`
}

type resolveRejectedPowerUpRequestsStatus struct {
	Type_   string `json:"type" example:"rejected"`
	Message string `json:"message" example:"insufficient fund"`
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
