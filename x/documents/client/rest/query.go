package rest

import (
	"fmt"
	"net/http"

	//restTypes "github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

const (
	userRestParameterName = "user"
	UUIDRestParameterName = "UUID"
)

func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/docs/{%s}/sent", userRestParameterName),
		getSentDocumentsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/docs/{%s}/received", userRestParameterName),
		getReceivedDocumentsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/receipts/{%s}/sent", userRestParameterName),
		getSentDocumentsReceiptsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/receipts/{%s}/received", userRestParameterName),
		getReceivedDocumentsReceiptsHandler(cliCtx)).
		Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/commercionetwork/docs/{%s}/receipts", UUIDRestParameterName),
		getDocumentsReceiptsHandler(cliCtx)).
		Methods("GET")
}

// ----------------------------------
// --- Documents
// ----------------------------------

func getSentDocumentsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[userRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, address, types.QuerySentDocuments)
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
		address := vars[userRestParameterName]

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
		address := vars[userRestParameterName]

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
		address := vars[userRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, address, types.QueryReceivedReceipts)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}

func getDocumentsReceiptsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		documentUUID := vars[UUIDRestParameterName]

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, documentUUID, types.QueryDocumentReceipts)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			restTypes.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		restTypes.PostProcessResponse(w, cliCtx, res)
	}
}
