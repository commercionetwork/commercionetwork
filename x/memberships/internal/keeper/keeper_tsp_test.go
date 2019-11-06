package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AddTrustedServiceProvider_EmptyList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	k.AddTrustedServiceProvider(ctx, testTsp)

	var signers ctypes.Addresses
	store := ctx.KVStore(k.storeKey)
	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	k.cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Len(t, signers, 1)
	assert.Contains(t, signers, testTsp)
}

func TestKeeper_AddTrustedServiceProvider_ExistingList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	signers := ctypes.Addresses{testTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), k.cdc.MustMarshalBinaryBare(&signers))

	k.AddTrustedServiceProvider(ctx, testUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustedSignersStoreKey))
	k.cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Len(t, actual, 2)
	assert.Contains(t, actual, testTsp)
	assert.Contains(t, actual, testUser)
}

func TestKeeper_GetTrustedServiceProviders_EmptyList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	signers := k.GetTrustedServiceProviders(ctx)
	assert.Empty(t, signers)
}

func TestKeeper_GetTrustedServiceProviders_ExistingList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	signers := ctypes.Addresses{testTsp, testUser, TestUser2}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.TrustedSignersStoreKey), k.cdc.MustMarshalBinaryBare(&signers))

	actual := k.GetTrustedServiceProviders(ctx)
	assert.Len(t, actual, 3)
	assert.Contains(t, actual, testTsp)
	assert.Contains(t, actual, testUser)
	assert.Contains(t, actual, TestUser2)
}

func TestKeeper_IsTrustedServiceProvider_EmptyList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	assert.False(t, k.IsTrustedServiceProvider(ctx, testTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, testUser))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}

func TestKeeper_IsTrustedServiceProvider_ExistingList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	signers := ctypes.Addresses{testUser, testTsp}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.TrustedSignersStoreKey), k.cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, k.IsTrustedServiceProvider(ctx, testUser))
	assert.True(t, k.IsTrustedServiceProvider(ctx, testTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}
