package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	_ = k.SetGovernmentAddress(ctx, TestAddress)

	path := []string{types.QueryGovernmentAddress}

	var actual sdk.AccAddress
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, TestAddress, actual)
}
