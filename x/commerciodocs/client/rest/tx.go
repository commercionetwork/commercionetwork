package rest

import (
	"fmt"
	types2 "github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(
		fmt.Sprintf("/%s/documents", storeName),
		storeDocumentHandler(cliCtx)).Methods("POST")
	r.HandleFunc(
		fmt.Sprintf("/%s/documents/{%s}/sharing", storeName, restName),
		shareDocumentHandler(cliCtx)).Methods("POST")
}

type (

	// ----------------------------------
	// --- StoreDocument
	// ----------------------------------
	storeDocumentReq struct {
		BaseReq   rest.BaseReq `json:"base_req"`
		Owner     string       `json:"owner"`
		Identity  string       `json:"identity"`
		Reference string       `json:"reference"`
		Metadata  string       `json:"metadata"`
	}

	// ----------------------------------
	// --- ShareDocument
	// ----------------------------------
	shareDocumentReq struct {
		BaseReq  rest.BaseReq `json:"base_req"`
		Owner    string       `json:"owner"`
		Sender   string       `json:"sender"`
		Receiver string       `json:"receiver"`
	}
)

func storeDocumentHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeDocumentReq
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
		msg := types.NewMsgStoreDocument(addr, req.Identity, req.Reference, req.Metadata)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func shareDocumentHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		var req shareDocumentReq
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
		msg := types.NewMsgShareDocument(addr, paramType, req.Sender, req.Receiver)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
