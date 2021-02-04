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

// ----------------------------------
// --- Vbr
// ----------------------------------

// @Summary Get All Current VBR pool funds
// @Description This endpoint returns current pool funds for validator block reward
// @ID getRetrieveBlockRewardsPoolFunds
// @Produce json
// @Success 200 {object} x.JSONResult{result=types.DecCoins}
// @Failure 404
// @Router /vbr/funds [get]
// @Tags x/vbr
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
