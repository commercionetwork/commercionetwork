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