package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/vbr/funds",
		getRetrieveBlockRewardsPoolFunds(cliCtx))
}

func getRetrieveBlockRewardsPoolFunds(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryBlockRewardsPoolFunds)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w,
				http.StatusInternalServerError,
				fmt.Sprintf("Could not get total funds amount: \n %s", err),
			)
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
