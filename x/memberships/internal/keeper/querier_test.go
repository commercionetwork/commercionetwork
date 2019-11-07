package keeper_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
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
	ctx, _, _, k := SetupTestInput()

	path := []string{types.QueryGetInvites, testUser.String()}

	var querier = keeper.NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
}

func Test_queryGetInvites_SpecificUser_Existing(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	invite := types.Invite{User: testUser, Sender: TestUser2}
	k.SaveInvite(ctx, invite)

	path := []string{types.QueryGetInvites, testUser.String()}

	var querier = keeper.NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 1, len(actual))
	assert.Contains(t, actual, invite)
}

func Test_queryGetInvites_Generic_Empty(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	path := []string{types.QueryGetInvites}

	var querier = keeper.NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
}

func Test_queryGetInvites_Generic_Existing(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	invite1 := types.Invite{User: testUser, Sender: TestUser2}
	invite2 := types.Invite{User: TestUser2, Sender: TestUser2}
	k.SaveInvite(ctx, invite1)
	k.SaveInvite(ctx, invite2)

	path := []string{types.QueryGetInvites}

	var querier = keeper.NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.Invite
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, invite1)
	assert.Contains(t, actual, invite2)
}

// ---------------------
// --- Signers
// ---------------------

func Test_queryGetSigners_EmptyList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	path := []string{types.QueryGetTrustedServiceProviders}

	var querier = keeper.NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual []sdk.AccAddress
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Empty(t, actual)
	assert.Equal(t, "[]", string(actualBz))
}

func Test_queryGetSigners_ExistingList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	expected := []sdk.AccAddress{testTsp, TestUser2}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), k.Cdc.MustMarshalBinaryBare(&expected))

	path := []string{types.QueryGetTrustedServiceProviders}

	var querier = keeper.NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual []sdk.AccAddress
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, testTsp)
	assert.Contains(t, actual, TestUser2)
}

// ---------------------
// --- Pool
// ---------------------

func Test_queryGetPoolFunds_EmptyPool(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	path := []string{types.QueryGetPoolFunds}

	var querier = keeper.NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual sdk.Coins
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetPoolFunds_ExistingPool(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)
	_ = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, expected)

	path := []string{types.QueryGetPoolFunds}

	var querier = keeper.NewQuerier(k)
	var request abci.RequestQuery
	actualBz, _ := querier(ctx, path, request)

	var actual sdk.Coins
	k.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, expected, actual)
}

func TestQuerier_resolveIdentity_Existent(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var querier = keeper.NewQuerier(k)

	membershipType := "bronze"
	_, _ = k.AssignMembership(ctx, testUser, membershipType)

	path := []string{types.QueryGetMembership, testUser.String()}
	bz, _ := querier(ctx, path, request)

	var actual keeper.MembershipResult
	k.Cdc.MustUnmarshalJSON(bz, &actual)

	expected := keeper.MembershipResult{User: testUser, MembershipType: membershipType}
	assert.Equal(t, expected, actual)
}

func TestQuerier_ResolveIdentity_NonExistent(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var querier = keeper.NewQuerier(k)

	path := []string{types.QueryGetMembership, "nunu"}
	_, err := querier(ctx, path, request)
	assert.Error(t, err)
}
