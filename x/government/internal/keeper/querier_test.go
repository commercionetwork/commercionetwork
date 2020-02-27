package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

// ----------------------------------
// --- Documents
// ----------------------------------

func TestQuerier_queryGetGovernmentAddress(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	err := k.SetGovernmentAddress(ctx, TestAddress)

	require.NoError(t, err)

	want := QueryGovernmentResponse{
		GovernmentAddress: TestAddress.String(),
	}

	path := []string{types.QueryGovernmentAddress}

	var actual QueryGovernmentResponse
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, want, actual)
}

func TestQuerier_queryGetTumblerAddress(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	err := k.SetTumblerAddress(ctx, TestAddress)

	require.NoError(t, err)

	want := QueryTumblerResponse{
		TumblerAddress: TestAddress.String(),
	}

	path := []string{types.QueryTumblerAddress}

	var actual QueryTumblerResponse
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, want, actual)
}
