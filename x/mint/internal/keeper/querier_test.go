package keeper

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var req abci.RequestQuery

func TestQuerier_queryGetCdp_foundCdp(t *testing.T) {
	cdc, ctx, _, _, k := SetupTestInput()

	k.AddCdp(ctx, TestCdp)

	parsedTimeStamp := TestCdp.Timestamp.Format(time.RFC3339)

	querier := NewQuerier(k)
	path := []string{types.QueryGetCdp, TestOwner.String(), parsedTimeStamp}
	actualBz, err := querier(ctx, path, req)

	var cdp types.Cdp
	cdc.MustUnmarshalJSON(actualBz, &cdp)
	assert.Nil(t, err)
	assert.Equal(t, TestCdp, cdp)
}

func TestQuerier_queryGetCdp_notFound(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	parsedTimeStamp := TestCdp.Timestamp.Format(time.RFC3339)

	path := []string{types.QueryGetCdp, TestOwner.String(), parsedTimeStamp}
	_, err := querier(ctx, path, req)

	assert.Error(t, err)
	expected := sdk.ErrUnknownRequest("couldn't find any cdp associated with the given address and timestamp")
	assert.Equal(t, expected, err)
}

func TestQuerier_queryGetCdps_found(t *testing.T) {
	cdc, ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	k.AddCdp(ctx, TestCdp)

	path := []string{types.QueryGetCdps, TestOwner.String(), TestCdp.Timestamp.String()}
	actualBz, err := querier(ctx, path, req)
	assert.Nil(t, err)

	var cdps types.Cdps
	cdc.MustUnmarshalJSON(actualBz, &cdps)
	assert.Equal(t, types.Cdps{TestCdp}, cdps)
}

func TestQuerier_queryGetCdps_notFound(t *testing.T) {
	cdc, ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCdps, TestOwner.String(), TestCdp.Timestamp.String()}
	actualBz, err := querier(ctx, path, req)
	assert.Nil(t, err)

	var cdps types.Cdps
	cdc.MustUnmarshalJSON(actualBz, &cdps)
	assert.Equal(t, types.Cdps(nil), cdps)
}
