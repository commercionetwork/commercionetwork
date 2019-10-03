package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var req abci.RequestQuery

func TestQuerier_queryGetCDP_foundCDP(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	k.AddCDP(ctx, TestCdp)

	path := []string{types.QueryGetCDP, TestOwner.String(), TestTimestamp}

	var cdp types.CDP
	actualBz, err := querier(ctx, path, req)
	k.Cdc.MustUnmarshalJSON(actualBz, &cdp)
	assert.Nil(t, err)
	assert.Equal(t, TestCdp, cdp)
}

func TestQuerier_queryGetCDP_notFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{types.QueryGetCDP, TestOwner.String(), TestTimestamp}
	_, err := querier(ctx, path, req)
	assert.Error(t, err)
	expected := sdk.ErrUnknownRequest("couldn't find any cdp associated with the given address and timestamp")
	assert.Equal(t, expected, err)
}

func TestQuerier_queryGetCDPs_found(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	k.AddCDP(ctx, TestCdp)

	var cdps types.CDPs
	path := []string{types.QueryGetCDPs, TestOwner.String(), TestTimestamp}
	actualBz, err := querier(ctx, path, req)
	k.Cdc.MustUnmarshalJSON(actualBz, &cdps)
	assert.Nil(t, err)
	assert.Equal(t, types.CDPs{TestCdp}, cdps)
}

func TestQuerier_queryGetCDPs_notFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	querier := NewQuerier(k)

	var cdps types.CDPs
	path := []string{types.QueryGetCDPs, TestOwner.String(), TestTimestamp}
	actualBz, err := querier(ctx, path, req)
	k.Cdc.MustUnmarshalJSON(actualBz, &cdps)
	assert.Nil(t, err)
	assert.Equal(t, types.CDPs(nil), cdps)

}
