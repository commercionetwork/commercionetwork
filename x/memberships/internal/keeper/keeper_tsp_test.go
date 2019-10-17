package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AddTrustedServiceProvider_EmptyList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	k.AddTrustedServiceProvider(ctx, TestTsp)

	var signers ctypes.Addresses
	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Equal(t, 1, len(signers))
	assert.Contains(t, signers, TestTsp)
}

func TestKeeper_AddTrustedServiceProvider_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	k.AddTrustedServiceProvider(ctx, TestUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
}

func TestKeeper_GetTrustedServiceProviders_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)

	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := k.GetTrustedServiceProviders(ctx)

	assert.Empty(t, signers)
}

func TestKeeper_GetTrustedServiceProviders_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestTsp, TestUser, TestUser2}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	actual := k.GetTrustedServiceProviders(ctx)
	assert.Equal(t, 3, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
	assert.Contains(t, actual, TestUser2)
}

func TestKeeper_IsTrustedServiceProvider_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	assert.False(t, k.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}

func TestKeeper_IsTrustedServiceProvider_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestUser, TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, k.IsTrustedServiceProvider(ctx, TestUser))
	assert.True(t, k.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}
