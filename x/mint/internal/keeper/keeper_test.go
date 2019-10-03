package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_SetCreditsDenom(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	denom := "test"
	k.SetCreditsDenom(ctx, denom)
	store := ctx.KVStore(k.StoreKey)
	denomBz := store.Get([]byte(types.CreditsDenomStoreKey))
	assert.Equal(t, denom, string(denomBz))
}

func TestKeeper_GetCreditsDenom(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	denom := "test"
	k.SetCreditsDenom(ctx, denom)
	actual := k.GetCreditsDenom(ctx)
	assert.Equal(t, denom, actual)
}

func TestKeeper_GetCDPkey(t *testing.T) {
	_, _, _, k := SetupTestInput()
	expected := types.CDPStoreKey + TestOwner.String()
	actual := k.GetCDPkey(TestOwner)
	assert.Equal(t, []byte(expected), actual)
}

func TestKeeper_GetUsersSet(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set(k.GetCDPkey(TestOwner), k.Cdc.MustMarshalBinaryBare(TestCdp))
	users := k.GetUsersSet(ctx)
	assert.Equal(t, TestOwner, users[0])
}

func TestKeeper_AddCDP_notAlreadyAdded(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var cdps types.CDPs
	k.AddCDP(ctx, TestCdp)

	store := ctx.KVStore(k.StoreKey)
	actualBz := store.Get(k.GetCDPkey(TestOwner))
	k.Cdc.MustUnmarshalBinaryBare(actualBz, &cdps)

	assert.Equal(t, TestCdp, cdps[0])
}

func TestKeeper_AddCDP_AlreadyAdded(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var cdps types.CDPs
	k.AddCDP(ctx, TestCdp)
	k.AddCDP(ctx, TestCdp)

	store := ctx.KVStore(k.StoreKey)
	actualBz := store.Get(k.GetCDPkey(TestOwner))
	k.Cdc.MustUnmarshalBinaryBare(actualBz, &cdps)

	assert.Len(t, cdps, 1)
}

func TestKeeper_GetCDPs(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var expected = types.CDPs{TestCdp}
	k.AddCDP(ctx, TestCdp)

	actual := k.GetCDPs(ctx, TestOwner)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetCDPs_empty(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var expected = types.CDPs(nil)

	actual := k.GetCDPs(ctx, TestOwner)
	assert.Equal(t, expected, actual)
}

func TestKeeper_DeleteCDP_deleted(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	k.AddCDP(ctx, TestCdp)
	actual := k.DeleteCDP(ctx, TestCdp)
	assert.True(t, actual)
}

func TestKeeper_DeleteCDP_notDeleted(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	actual := k.DeleteCDP(ctx, TestCdp)
	assert.False(t, actual)
}

func TestKeeper_OpenCDP_InvalidDepositedAmount(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	invalidReq := types.CDPRequest{
		Signer:          TestOwner,
		DepositedAmount: sdk.NewCoins(sdk.NewCoin("testcoin", sdk.NewInt(0))),
		Timestamp:       "",
	}
	err := k.OpenCDP(ctx, invalidReq)
	assert.Error(t, err)
	expected := sdk.ErrInvalidCoins(invalidReq.DepositedAmount.String())
	assert.Equal(t, expected, err)
}
