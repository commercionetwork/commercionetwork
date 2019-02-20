package rest

import (
	"commercio-network/types"
	"fmt"
	"net/http"

	"commercio-network/x/commerciodocs"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gorilla/mux"
)

const (
	restName = "document"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/documents", storeName),
		storeDocumentHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/documents/{%s}", storeName, restName),
		getDocumentMetadataHandler(cdc, cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/documents/{%s}/sharing", storeName, restName),
		shareDocumentHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/documents/{%s}/readers", storeName, restName),
		getDocumentReadersHandler(cdc, cliCtx, storeName)).Methods("GET")
}

// ----------------------------------
// --- StoreDocument
// ----------------------------------

type storeDocumentReq struct {
	BaseReq   utils.BaseReq `json:"base_req"`
	Owner     string        `json:"owner"`
	Identity  string        `json:"identity"`
	Reference string        `json:"reference"`
	Metadata  string        `json:"metadata"`
}

func storeDocumentHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req storeDocumentReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := commerciodocs.NewMsgStoreDocument(addr, types.Did(req.Identity), req.Reference, req.Metadata)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}

// ----------------------------------
// --- GetDocumentMetadata
// ----------------------------------

func getDocumentMetadataHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/metadata/%s", storeName, paramType), nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// ----------------------------------
// --- ShareDocument
// ----------------------------------

type shareDocumentReq struct {
	BaseReq  utils.BaseReq `json:"base_req"`
	Owner    string        `json:"owner"`
	Sender   string        `json:"sender"`
	Receiver string        `json:"receiver"`
}

func shareDocumentHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		var req shareDocumentReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := commerciodocs.NewMsgShareDocument(addr, paramType, types.Did(req.Sender), types.Did(req.Receiver))
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}

// ----------------------------------
// --- GetDocumentReaders
// ----------------------------------

func getDocumentReadersHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/readers/%s", storeName, paramType), nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
