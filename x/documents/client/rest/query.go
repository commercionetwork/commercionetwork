package rest

import (
	"fmt"
	"net/http"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	restTypes "github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)
const (
	addressRestParameterName = "user"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/docs/{%s}/sent", addressRestParameterName),
		getSentDocumentsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/docs/{%s}/received", addressRestParameterName),
		getReceivedDocumentsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/receipts/{%s}/sent", addressRestParameterName),
		getSentDocumentsReceiptsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/receipts/{%s}/received", addressRestParameterName),
		getReceivedDocumentsReceiptsHandler(cliCtx)).
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
// @Tags x/documents
func getSentDocumentsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute,  address, types.QuerySentDocuments)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
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
// @Tags x/documents
func getReceivedDocumentsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, address, types.QueryReceivedDocuments)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
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
// @Tags x/documents
func getSentDocumentsReceiptsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, address, types.QuerySentReceipts)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
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
// @Tags x/documents
func getReceivedDocumentsReceiptsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[addressRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, address, types.QueryReceivedReceipts)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}