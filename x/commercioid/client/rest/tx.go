package rest

import (
	"commercio-network/types"
	types2 "commercio-network/x/commercioid/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf(
		"/%s/identities", storeName),
		postUpsertIdentityHandlerFn(cliCtx)).
		Methods("PUT")
	r.HandleFunc(fmt.Sprintf(
		"/%s/connections", storeName),
		postCreateConnectionHandlerFn(cliCtx)).
		Methods("POST")
}

type (

	// ----------------------------------
	// --- Upsert identity
	// ----------------------------------
	upsertIdentityReq struct {
		BaseReq      rest.BaseReq `json:"base_req"`
		Owner        string       `json:"owner"`
		Did          string       `json:"did"`
		DdoReference string       `json:"ddo_reference"`
	}

	// ----------------------------------
	// --- CreateConnection
	// ----------------------------------
	createConnectionReq struct {
		BaseReq    rest.BaseReq `json:"base_req"`
		Owner      string       `json:"owner"`
		FirstUser  string       `json:"first_user"`
		SecondUser string       `json:"second_user"`
	}
)

func postUpsertIdentityHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req upsertIdentityReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types2.NewMsgSetIdentity(types.Did(req.Did), req.DdoReference, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postCreateConnectionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createConnectionReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types2.NewMsgCreateConnection(types.Did(req.FirstUser), types.Did(req.SecondUser), addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
