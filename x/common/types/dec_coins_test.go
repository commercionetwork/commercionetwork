package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestIsAllGTE(t *testing.T) {

	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	testDenom1 := "atom"
	testDenom2 := "muon"

	assert.True(t, types.IsAllGTE(sdk.DecCoins{}, sdk.DecCoins{}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}}, sdk.DecCoins{}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}, {testDenom2, one}}, sdk.DecCoins{{testDenom2, one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}, {testDenom2, one}}, sdk.DecCoins{{testDenom2, two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{}, sdk.DecCoins{{testDenom1, one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}}, sdk.DecCoins{{testDenom1, two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}}, sdk.DecCoins{{testDenom2, one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}, {testDenom2, two}}, sdk.DecCoins{{testDenom1, two}, {testDenom2, one}}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}}, sdk.DecCoins{{testDenom1, one}}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, two}}, sdk.DecCoins{{testDenom1, one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}}, sdk.DecCoins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom2, two}}, sdk.DecCoins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom2, one}}, sdk.DecCoins{{testDenom1, one}, {testDenom2, two}}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}, {testDenom2, two}}, sdk.DecCoins{{testDenom1, one}, {testDenom2, one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{testDenom1, one}, {testDenom2, one}}, sdk.DecCoins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{"xxx", one}, {"yyy", one}}, sdk.DecCoins{{testDenom2, one}, {"ccc", one}, {"yyy", one}, {"zzz", one}}))
}
