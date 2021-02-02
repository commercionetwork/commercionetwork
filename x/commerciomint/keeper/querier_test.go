package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

var req abci.RequestQuery

func TestQuerier_queryGetEtps(t *testing.T) {
	ctx, _, _, _, k := SetupTestInput()

	k.SetPosition(ctx, testEtp)

	querier := NewQuerier(k)
	path := []string{types.QueryGetEtps, testEtpOwner.String()}
	actualBz, err := querier(ctx, path, req)

	var etps []types.Position
	k.cdc.MustUnmarshalJSON(actualBz, &etps)
	require.Nil(t, err)
	require.True(t, etps[0].Equals(testEtp))
}

func TestQuerier_queryConversionRate(t *testing.T) {
	ctx, _, _, _, k := SetupTestInput()
	require.NoError(t, k.SetConversionRate(ctx, sdk.NewDec(2)))
	querier := NewQuerier(k)
	actualBz, err := querier(ctx, []string{types.QueryConversionRate}, req)
	require.Nil(t, err)

	var rate sdk.Dec
	k.cdc.MustUnmarshalJSON(actualBz, &rate)
	require.Equal(t, sdk.NewDec(2), rate)
}
