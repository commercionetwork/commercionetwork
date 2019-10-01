package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

const (
	identityParam = "identity"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, querierRoute string) {
	r.HandleFunc(fmt.Sprintf(
		"/identities/{%s}", identityParam),
		resolveIdentityHandler(cliCtx, querierRoute)).
		Methods("GET")
}

func resolveIdentityHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[identityParam]

		route := fmt.Sprintf("custom/%s/identities/%s", querierRoute, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
