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
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{}, sdk.DecCoins{{Denom: testDenom1, Amount: one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: two}, {Denom: testDenom2, Amount: one}}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom2, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	assert.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	assert.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: "xxx", Amount: one}, {Denom: "yyy", Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: one}, {Denom: "ccc", Amount: one}, {Denom: "yyy", Amount: one}, {Denom: "zzz", Amount: one}}))
}
