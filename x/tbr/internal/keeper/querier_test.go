package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestQuerier_getBlockRewardsPoolFunds(t *testing.T) {
	var cdc, ctx, k, _, _ = SetupTestInput()
	var querier = NewQuerier(k)

	k.SetTotalRewardPool(ctx, TestBlockRewardsPool)

	path := []string{types.QueryBlockRewardsPoolFunds}
	actual, _ := querier(ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(cdc, &TestBlockRewardsPool)

	assert.Equal(t, expected, actual)
}
