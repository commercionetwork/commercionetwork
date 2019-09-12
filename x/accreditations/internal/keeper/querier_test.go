package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ---------------------
// --- Accrediters
// ---------------------

func Test_queryGetAccrediter_NonExistent(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete(TestUser)

	path := []string{types.QueryGetAccrediter, TestUser.String()}

	var querier = NewQuerier(accreditationKeeper)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual AccrediterResponse
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, TestUser, actual.User)
	assert.True(t, actual.Accrediter.Empty())
}

func Test_queryGetAccrediter_Existent(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	accreditation := types.Accreditation{User: TestUser, Accrediter: TestAccrediter}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	path := []string{types.QueryGetAccrediter, TestUser.String()}

	var querier = NewQuerier(accreditationKeeper)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual AccrediterResponse
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, TestUser, actual.User)
	assert.Equal(t, TestAccrediter, actual.Accrediter)
}

// ---------------------
// --- Signers
// ---------------------

func Test_queryGetSigners_EmptyList(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete([]byte(types.TrustedSignersStoreKey))

	path := []string{types.QueryGetSigners}

	var querier = NewQuerier(accreditationKeeper)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual []sdk.AccAddress
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
	assert.Equal(t, "[]", string(actualBz))
}

func Test_queryGetSigners_ExistingList(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	expected := []sdk.AccAddress{TestSigner, TestAccrediter}

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&expected))

	path := []string{types.QueryGetSigners}

	var querier = NewQuerier(accreditationKeeper)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual []sdk.AccAddress
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestSigner)
	assert.Contains(t, actual, TestAccrediter)
}

// ---------------------
// --- Pool
// ---------------------

func Test_queryGetPoolFunds_EmptyPool(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete([]byte(types.LiquidityPoolKey))

	path := []string{types.QueryGetPoolFunds}

	var querier = NewQuerier(accreditationKeeper)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual sdk.Coins
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetPoolFunds_ExistingPool(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolKey), cdc.MustMarshalBinaryBare(&expected))

	path := []string{types.QueryGetPoolFunds}

	var querier = NewQuerier(accreditationKeeper)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual sdk.Coins
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, expected, actual)
}
