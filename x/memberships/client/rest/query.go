package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	restParamAddress = "user"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/memberships/{%s}", restParamAddress),
		getMembershipHandler(cliCtx)).
		Methods("GET")
}

func getMembershipHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
