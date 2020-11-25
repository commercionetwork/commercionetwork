package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commercionetwork/commercionetwork/x/docs/types"
)

const (
	addressRestParameterName = "user"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/docs/{%s}/sent", addressRestParameterName),
		getSentDocumentsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/docs/{%s}/received", addressRestParameterName),
		getReceivedDocumentsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/receipts/{%s}/sent", addressRestParameterName),
		getSentDocumentsReceiptsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/receipts/{%s}/received", addressRestParameterName),
		getReceivedDocumentsReceiptsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		"/docs/metadataSchemes",
		getSupportedMetadataSchemesHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		"/docs/metadataSchemes/proposers",
		getTrustedMetadataSchemesProposersHandler(cliCtx)).
		Methods("GET")
}

// ----------------------------------
// --- Documents
// ----------------------------------

// @Summary Get the sent documents
// @Description This endpoint returns the sent documents, along with the height at which the resource was queried at
// @ID getSentDocumentsHandler
// @Produce json
// @Param address path string true "Address of the user"
// @Success 200 {object} x.JSONResult{result=[]types.Document}
// @Failure 404
// @Router /docs/{address}/sent [get]
// @Tags x/docs
func getSentDocumentsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentDocuments, address)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get the received documents
// @Description This endpoint returns the received documents, along with the height at which the resource was queried at
// @ID getReceivedDocumentsHandler
// @Produce json
// @Param address path string true "Address of the user"
// @Success 200 {object} x.JSONResult{result=[]types.Document}
// @Failure 404
// @Router /docs/{address}/received [get]
// @Tags x/docs
func getReceivedDocumentsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryReceivedDocuments, address)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ---------------------------------
// --- Document receipts
// ---------------------------------

// @Summary Get the sent receipts
// @Description This endpoint returns the sent receipts, along with the height at which the resource was queried at
// @ID getSentDocumentsReceiptsHandler
// @Produce json
// @Param address path string true "Address of the user"
// @Success 200 {object} x.JSONResult{result=[]types.DocumentReceipt}
// @Failure 404
// @Router /receipts/{address}/sent [get]
// @Tags x/docs
func getSentDocumentsReceiptsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentReceipts, address)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// @Summary Get the received receipts
// @Description This endpoint returns the received receipts, along with the height at which the resource was queried at
// @ID getReceivedDocumentsReceiptsHandler
// @Produce json
// @Param address path string true "Address of the user"
// @Success 200 {object} x.JSONResult{result=[]types.DocumentReceipt}
// @Failure 404
// @Router /receipts/{address}/received [get]
// @Tags x/docs
func getReceivedDocumentsReceiptsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryReceivedReceipts, address)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ----------------------------------
// --- Document metadata schemes
// ----------------------------------

// @Summary Get the metadata schemes
// @Description This endpoint returns the supported metadata schemes, along with the height at which the resource was queried at
// @ID getSupportedMetadataSchemesHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]types.MetadataSchema}
// @Failure 404
// @Router /docs/metadataSchemes [get]
// @Tags x/docs
func getSupportedMetadataSchemesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySupportedMetadataSchemes)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// -----------------------------------------
// --- Document metadata schemes proposers
// -----------------------------------------

// @Summary Get the metadata proposers
// @Description This endpoint returns the trusted metadata proposers, along with the height at which the resource was queried at
// @ID getSupportedMetadataSchemesHandler
// @Produce json
// @Success 200 {object} x.JSONResult{result=[]string}
// @Failure 404
// @Router /docs/metadataSchemes/proposers [get]
// @Tags x/docs
func getTrustedMetadataSchemesProposersHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTrustedMetadataProposers)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
