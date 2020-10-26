package keeper

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

var req abci.RequestQuery

func TestQuerier_queryGetCdp_foundCdp(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()

	k.SetPosition(ctx, testCdp)

	querier := NewQuerier(k)
	path := []string{types.QueryGetCdp, testCdpOwner.String(), strconv.FormatInt(testCdp.CreatedAt, 10)}
	actualBz, err := querier(ctx, path, req)

	var cdp types.Position
	k.cdc.MustUnmarshalJSON(actualBz, &cdp)
	require.Nil(t, err)
	require.Equal(t, testCdp, cdp)
}

func TestQuerier_queryGetCdp_notFound(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCdp, testCdpOwner.String(), strconv.FormatInt(testCdp.CreatedAt, 10)}
	_, err := querier(ctx, path, req)

	require.Error(t, err)
	expected := sdkErr.Wrap(sdkErr.ErrUnknownRequest, "couldn't find any cdp associated with the given address and timestamp")
	require.Equal(t, expected.Error(), err.Error())
}

func TestQuerier_queryGetCdps_found(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	k.SetPosition(ctx, testCdp)

	path := []string{types.QueryGetEtps, testCdpOwner.String(), strconv.FormatInt(testCdp.CreatedAt, 10)}
	actualBz, err := querier(ctx, path, req)
	require.Nil(t, err)

	var cdps []types.Position
	k.cdc.MustUnmarshalJSON(actualBz, &cdps)
	require.Equal(t, []types.Position{testCdp}, cdps)
}

func TestQuerier_queryGetCdps_notFound(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetEtps, testCdpOwner.String(), strconv.FormatInt(testCdp.CreatedAt, 10)}
	actualBz, err := querier(ctx, path, req)
	require.Nil(t, err)

	var cdps []types.Position
	k.cdc.MustUnmarshalJSON(actualBz, &cdps)
	require.Equal(t, []types.Position(nil), cdps)
}

func TestQuerier_queryCollateralRate(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	require.NoError(t, k.SetConversionRate(ctx, sdk.NewInt(2).ToDec()))
	querier := NewQuerier(k)
	actualBz, err := querier(ctx, []string{"collateral_rate"}, req)
	require.Nil(t, err)

	var rate sdk.Dec
	k.cdc.MustUnmarshalJSON(actualBz, &rate)
	require.Equal(t, sdk.NewDec(2), rate)
}
