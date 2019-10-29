package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AddTrustedServiceProvider_EmptyList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	k.AddTrustedServiceProvider(ctx, TestTsp)

	var signers ctypes.Addresses
	store := ctx.KVStore(k.StoreKey)
	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Len(t, signers, 1)
	assert.Contains(t, signers, TestTsp)
}

func TestKeeper_AddTrustedServiceProvider_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	signers := ctypes.Addresses{TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	k.AddTrustedServiceProvider(ctx, TestUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Len(t, actual, 2)
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
}

func TestKeeper_GetTrustedServiceProviders_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	signers := k.GetTrustedServiceProviders(ctx)
	assert.Empty(t, signers)
}

func TestKeeper_GetTrustedServiceProviders_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	signers := ctypes.Addresses{TestTsp, TestUser, TestUser2}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	actual := k.GetTrustedServiceProviders(ctx)
	assert.Len(t, actual, 3)
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
	assert.Contains(t, actual, TestUser2)
}

func TestKeeper_IsTrustedServiceProvider_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}

func TestKeeper_IsTrustedServiceProvider_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	signers := ctypes.Addresses{TestUser, TestTsp}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, k.IsTrustedServiceProvider(ctx, TestUser))
	assert.True(t, k.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}
