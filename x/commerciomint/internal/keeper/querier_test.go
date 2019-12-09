package keeper

import (
	"strconv"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
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
	assert.Nil(t, err)
	assert.Equal(t, testCdp, cdp)
}

func TestQuerier_queryGetCdp_notFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCdp, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	_, err := querier(ctx, path, req)

	assert.Error(t, err)
	expected := sdk.ErrUnknownRequest("couldn't find any cdp associated with the given address and timestamp")
	assert.Equal(t, expected, err)
}

func TestQuerier_queryGetCdps_found(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	k.AddCdp(ctx, testCdp)

	path := []string{types.QueryGetCdps, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	actualBz, err := querier(ctx, path, req)
	assert.Nil(t, err)

	var cdps types.Cdps
	k.cdc.MustUnmarshalJSON(actualBz, &cdps)
	assert.Equal(t, types.Cdps{testCdp}, cdps)
}

func TestQuerier_queryGetCdps_notFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCdps, testCdpOwner.String(), strconv.FormatInt(testCdp.Timestamp, 10)}
	actualBz, err := querier(ctx, path, req)
	assert.Nil(t, err)

	var cdps types.Cdps
	k.cdc.MustUnmarshalJSON(actualBz, &cdps)
	assert.Equal(t, types.Cdps(nil), cdps)
}
