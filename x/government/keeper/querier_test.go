package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

var request abci.RequestQuery

// ----------------------------------
// --- Documents
// ----------------------------------

func TestQuerier_queryGetGovernmentAddress(t *testing.T) {
	cdc, ctx, k := SetupTestInput(true)
	var querier = NewQuerier(k)

	want := QueryGovernmentResponse{
		GovernmentAddress: governmentTestAddress.String(),
	}

	path := []string{types.QueryGovernmentAddress}

	var actual QueryGovernmentResponse
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, want, actual)
}

func TestQuerier_queryGetTumblerAddress(t *testing.T) {
	cdc, ctx, k := SetupTestInput(true)
	var querier = NewQuerier(k)

	want := QueryTumblerResponse{
		TumblerAddress: tumblerTestAddress.String(),
	}

	path := []string{types.QueryTumblerAddress}

	var actual QueryTumblerResponse
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, want, actual)
}
