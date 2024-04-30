package rest

// import (
// 	"fmt"
// 	"net/http"

// 	//restTypes "github.com/cosmos/cosmos-sdk/types/rest"

// 	"github.com/commercionetwork/commercionetwork/x/did/types"
// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/gorilla/mux"
// )

// const (
// 	identityParam = "identity"
// )

// func RegisterRoutes(cliCtx client.Context, r *mux.Router, querierRoute string) {
// 	r.HandleFunc(fmt.Sprintf(
// 		"/identities/{%s}", identityParam),
// 		resolveIdentityHandler(cliCtx, querierRoute)).
// 		Methods("GET")
// }

// func resolveIdentityHandler(cliCtx client.Context, querierRoute string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		paramType := vars[identityParam]

// 		route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveIdentity, paramType)
// 		res, _, err := cliCtx.QueryWithData(route, nil)
// 		if err != nil {
// 			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
// 			return
// 		}

// 		restTypes.PostProcessResponse(w, cliCtx, res)
// 	}
// }
