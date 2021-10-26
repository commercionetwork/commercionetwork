package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/client"
	restTypes "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/vbr/funds", getRetrieveBlockRewardsPoolFunds(cliCtx)).Methods("GET")
	r.HandleFunc("/vbr/reward_rate", getRewardRateHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/vbr/automatic_withdraw", getAutomaticWithdrawHandler(cliCtx)).Methods("GET")
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
func getRetrieveBlockRewardsPoolFunds(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryBlockRewardsPoolFunds)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w,
				http.StatusInternalServerError,
				fmt.Sprintf("Could not get total funds amount: \n %s", err),
			)
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get Reward rate
// @Description This endpoint returns the Reward rate, along with the height at which the resource was queried at
// @ID getRewardRateHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=types.Dec}
// @Failure 404
// @Router /vbr/reward_rate [get]
// @Tags x/vbr
func getRewardRateHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardRate)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get Automatic withdraw
// @Description This endpoint returns the Automatic withdraw flag, along with the height at which the resource was queried at
// @ID getAutomaticWithdrawHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=bool}
// @Failure 404
// @Router /vbr/automatic_withdraw [get]
// @Tags x/vbr
func getAutomaticWithdrawHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAutomaticWithdraw)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
