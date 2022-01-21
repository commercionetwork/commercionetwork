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
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router, querierRoute string) {
	r.HandleFunc(fmt.Sprintf(
		"/identities/{%s}", identityParam),
		resolveIdentityHandler(cliCtx, querierRoute)).
		Methods("GET")
}

// @Summary Get a user DID document
// @Description This endpoint returns a user DID document, along with the height at which the resource was queried at
// @ID id_resolveIdentityHandler
// @Produce json
// @Param did path string true "Address of the user for which to read the DID document"
// @Success 200 {object} x.JSONResult{result=keeper.ResolveIdentityResponse}
// @Failure 404
// @Router /identities/{did} [get]
// @Tags x/did
func resolveIdentityHandler(cliCtx client.Context, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[identityParam]

		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveIdentity, paramType)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
