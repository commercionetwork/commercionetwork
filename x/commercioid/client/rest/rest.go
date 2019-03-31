package rest

import (
	"commercio-network/types"
	"fmt"
	"net/http"

	"commercio-network/x/commercioid"
	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

const (
	restName = "identity"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/identities", storeName),
		upsertIdentityHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/connections", storeName),
		createConnectionHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/identities/{%s}", storeName, restName),
		resolveIdentityHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/identities/{%s}/connections", storeName, restName),
		getConnectionsHandler(cdc, cliCtx, storeName)).Methods("GET")
}

// ----------------------------------
// --- Upsert identity
// ----------------------------------

type upsertIdentityReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	Owner        string       `json:"owner"`
	Did          string       `json:"did"`
	DdoReference string       `json:"ddo_reference"`
}

func upsertIdentityHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req upsertIdentityReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
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
		msg := commercioid.NewMsgSetIdentity(types.Did(req.Did), req.DdoReference, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		//clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq,[]sdk.Msg{msg})
	}
}

// ----------------------------------
// --- CreateConnection
// ----------------------------------

type createConnectionReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Owner      string       `json:"owner"`
	FirstUser  string       `json:"first_user"`
	SecondUser string       `json:"second_user"`
}

func createConnectionHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createConnectionReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
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
		msg := commercioid.NewMsgCreateConnection(types.Did(req.FirstUser), types.Did(req.SecondUser), addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.CompleteAndBroadcastTxREST(w, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
		//clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq,[]sdk.Msg{msg})
	}
}

// ----------------------------------
// --- ResolveIdentity
// ----------------------------------

func resolveIdentityHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/identities/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// ----------------------------------
// --- GetConnections
// ----------------------------------

func getConnectionsHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/connections/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
