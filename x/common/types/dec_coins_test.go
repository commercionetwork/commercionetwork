package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/common/types"
)

func TestIsAllGTE(t *testing.T) {

	one := sdk.NewDec(1)
	two := sdk.NewDec(2)

	testDenom1 := "atom"
	testDenom2 := "muon"

	require.True(t, types.IsAllGTE(sdk.DecCoins{}, sdk.DecCoins{}))
	require.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{}))
	require.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: one}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: two}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{}, sdk.DecCoins{{Denom: testDenom1, Amount: one}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: two}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: one}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: two}, {Denom: testDenom2, Amount: one}}))
	require.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}}))
	require.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom2, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	require.True(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: one}}, sdk.DecCoins{{Denom: testDenom1, Amount: one}, {Denom: testDenom2, Amount: two}}))
	require.False(t, types.IsAllGTE(sdk.DecCoins{{Denom: "xxx", Amount: one}, {Denom: "yyy", Amount: one}}, sdk.DecCoins{{Denom: testDenom2, Amount: one}, {Denom: "ccc", Amount: one}, {Denom: "yyy", Amount: one}, {Denom: "zzz", Amount: one}}))
}
