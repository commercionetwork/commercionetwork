package rest

import (
	"fmt"
	"net/http"

	"commercio-network/x/commercioauth"
	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/register", storeName),
		storeDocumentHandler(cdc, cliCtx)).Methods("POST")
}

// ----------------------------------
// --- StoreDocument
// ----------------------------------

type registerAccountReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Signer   string       `json:"signer"`
	Address  string       `json:"address"`
	KeyType  string       `json:"key_type"`
	KeyValue string       `json:"key_value"`
}

func storeDocumentHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerAccountReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Signer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := commercioauth.NewMsgCreateAccount(addr, req.Address, req.KeyType, req.KeyValue)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		//clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq,[]sdk.Msg{msg})

	}
}
