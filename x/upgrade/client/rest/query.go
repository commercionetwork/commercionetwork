package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	restTypes "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)


func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/commercionetwork/upgrade/currentUpgrade",getCurrentUpgrade(cliCtx)).Methods("GET")
}

func getCurrentUpgrade(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrent)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
