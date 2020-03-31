package keeper

import (
	"strconv"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

var req abci.RequestQuery

func TestQuerier_queryGetCdp_foundCdp(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	k.AddCdp(ctx, testCdp)

	querier := NewQuerier(k)
	path := []string{types.QueryGetCdp, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	actualBz, err := querier(ctx, path, req)

	var cdp types.Cdp
	k.cdc.MustUnmarshalJSON(actualBz, &cdp)
	require.Nil(t, err)
	require.Equal(t, testCdp, cdp)
}

func TestQuerier_queryGetCdp_notFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCdp, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	_, err := querier(ctx, path, req)

	require.Error(t, err)
	expected := sdkErr.Wrap(sdkErr.ErrUnknownRequest, "couldn't find any cdp associated with the given address and timestamp")
	require.Equal(t, expected.Error(), err.Error())
}

func TestQuerier_queryGetCdps_found(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	k.AddCdp(ctx, testCdp)

	path := []string{types.QueryGetCdps, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	actualBz, err := querier(ctx, path, req)
	require.Nil(t, err)

	var cdps types.Cdps
	k.cdc.MustUnmarshalJSON(actualBz, &cdps)
	require.Equal(t, types.Cdps{testCdp}, cdps)
}

func TestQuerier_queryGetCdps_notFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCdps, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	actualBz, err := querier(ctx, path, req)
	require.Nil(t, err)

	var cdps types.Cdps
	k.cdc.MustUnmarshalJSON(actualBz, &cdps)
	require.Equal(t, types.Cdps(nil), cdps)
}
