package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/government/address",
		getGovernmentAddr(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		"/government/tumbler",
		getTumblerAddr(cliCtx)).
		Methods("GET")
}

// @Summary Get the government address
// @Description This endpoint returns the address that the government has currently, along with the height at which the resource was queried at
// @ID government_address
// @Produce json
// @Success 200 {object} rest.getGovernmentAddrResponse
// @Failure 404
// @Router /government/address [get]
// @Tags x/government
func getGovernmentAddr(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGovernmentAddress)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

type getGovernmentAddrResponse struct {
	Height int `json:"height" example:"1234"`
	Result struct {
		Government_address string `json:"government_address" example:"did:com:1pxelxpwdjqsz23kpz87lmc2qgmkjd35x7uq9zv"`
	} `json:"result"`
}

var _ getGovernmentAddrResponse

// @Summary Get the tumbler address
// @Description This endpoint returns the address that the tumbler has currently, along with the height at which the resource was queried at
// @ID government_tumbler
// @Produce json
// @Success 200 {object} rest.getTumblerAddrResponse
// @Failure 404
// @Router /government/tumbler [get]
// @Tags x/government
func getTumblerAddr(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTumblerAddress)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

type getTumblerAddrResponse struct {
	Height int `json:"height" example:"1234"`
	Result struct {
		Tumbler_address string `json:"tumbler_address" example:"did:com:1cqq6qveuqkqxek4up92asqkawlqsxh9k5xnnll"`
	} `json:"result"`
}

var _ getTumblerAddrResponse
