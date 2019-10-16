package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

// ---------------------
// --- Invites
// ---------------------

func Test_queryGetInvites_SpecificUser_Empty(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	path := []string{types.QueryGetInvites, TestUser.String()}

	var querier = NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
}

func Test_queryGetInvites_SpecificUser_Existing(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	invite := types.Invite{User: TestUser, Sender: TestUser2}
	k.SaveInvite(ctx, invite)

	path := []string{types.QueryGetInvites, TestUser.String()}

	var querier = NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 1, len(actual))
	assert.Contains(t, actual, invite)
}

func Test_queryGetInvites_Generic_Empty(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	path := []string{types.QueryGetInvites}

	var querier = NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
}

func Test_queryGetInvites_Generic_Existing(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	invite1 := types.Invite{User: TestUser, Sender: TestUser2}
	invite2 := types.Invite{User: TestUser2, Sender: TestUser2}
	k.SaveInvite(ctx, invite1)
	k.SaveInvite(ctx, invite2)

	path := []string{types.QueryGetInvites}

	var querier = NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, invite1)
	assert.Contains(t, actual, invite2)
}

// ---------------------
// --- Signers
// ---------------------

func Test_queryGetSigners_EmptyList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	path := []string{types.QueryGetTrustedServiceProviders}

	var querier = NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual []sdk.AccAddress
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
	assert.Equal(t, "[]", string(actualBz))
}

func Test_queryGetSigners_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	expected := []sdk.AccAddress{TestTsp, TestUser2}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&expected))

	path := []string{types.QueryGetTrustedServiceProviders}

	var querier = NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual []sdk.AccAddress
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser2)
}

// ---------------------
// --- Pool
// ---------------------

func Test_queryGetPoolFunds_EmptyPool(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	path := []string{types.QueryGetPoolFunds}

	var querier = NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual sdk.Coins
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetPoolFunds_ExistingPool(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&expected))

	path := []string{types.QueryGetPoolFunds}

	var querier = NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual sdk.Coins
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, expected, actual)
}

func TestQuerier_resolveIdentity_Existent(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()
	var querier = NewQuerier(k)

	membershipType := "bronze"
	_, _ = k.AssignMembership(ctx, TestUser, membershipType)

	path := []string{types.QueryGetMembership, TestUser.String()}
	bz, _ := querier(ctx, path, request)

	var actual MembershipResult
	cdc.MustUnmarshalJSON(bz, &actual)

	expected := MembershipResult{User: TestUser, MembershipType: membershipType}
	assert.Equal(t, expected, actual)
}

func TestQuerier_ResolveIdentity_NonExistent(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryGetMembership, "nunu"}
	_, err := querier(ctx, path, request)
	assert.Error(t, err)
}
