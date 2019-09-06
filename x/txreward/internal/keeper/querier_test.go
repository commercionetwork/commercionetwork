package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var cdc, ctx, k = SetupTestInput()
var querier = NewQuerier(k)
var request abci.RequestQuery

func TestQuerier_getBlockRewardsPoolFunds(t *testing.T) {

	k.setBlockRewardsPool(ctx, TestBlockRewardsPool)

	path := []string{types.QueryBlockRewardsPoolFunds}
	actual, _ := querier(ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(cdc, &TestBlockRewardsPool)

	assert.Equal(t, expected, actual)
}

func TestQuerier_getBlockRewardsPoolFunders(t *testing.T) {
	k.setFunders(ctx, TestFunders)

	path := []string{types.QueryBlockRewardsPoolFunders}
	actual, _ := querier(ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(cdc, &TestFunders)

	assert.Equal(t, expected, actual)
}
