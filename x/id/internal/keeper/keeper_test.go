package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_CreateIdentity(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)

	k.SaveIdentity(ctx, TestOwnerAddress, TestDidDocument)

	var stored types.DidDocument
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, TestDidDocument, stored)
}

func TestKeeper_EditIdentity(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	updatedDidDocument := types.DidDocument{Uri: "ddo-reference-update", ContentHash: TestDidDocument.ContentHash}
	k.SaveIdentity(ctx, TestOwnerAddress, updatedDidDocument)

	var stored types.DidDocument
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, updatedDidDocument, stored)
}

func TestKeeper_GetDidDocumentByOwner_ExistingDidDocument(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, found := k.GetDidDocumentByOwner(ctx, TestOwnerAddress)

	assert.True(t, found)
	assert.Equal(t, TestDidDocument, actual)
}

// -------------------------
// --- Genesis utils
// -------------------------

func TestKeeper_GetIdentities(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	first, _ := sdk.AccAddressFromBech32("cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second, _ := sdk.AccAddressFromBech32("cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third, _ := sdk.AccAddressFromBech32("cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth, _ := sdk.AccAddressFromBech32("cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth, _ := sdk.AccAddressFromBech32("cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

	store.Set(k.getIdentityStoreKey(first), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(second), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(third), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(fourth), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(fifth), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, err := k.GetIdentities(ctx)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(actual))
	assert.Contains(t, actual, types.Identity{Owner: first, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: second, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: third, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: fourth, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: fifth, DidDocument: TestDidDocument})
}

func TestKeeper_SetIdentities(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	first, _ := sdk.AccAddressFromBech32("cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second, _ := sdk.AccAddressFromBech32("cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third, _ := sdk.AccAddressFromBech32("cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth, _ := sdk.AccAddressFromBech32("cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth, _ := sdk.AccAddressFromBech32("cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

	identities := []types.Identity{
		{Owner: first, DidDocument: TestDidDocument},
		{Owner: second, DidDocument: TestDidDocument},
		{Owner: third, DidDocument: TestDidDocument},
		{Owner: fourth, DidDocument: TestDidDocument},
		{Owner: fifth, DidDocument: TestDidDocument},
	}
	k.SetIdentities(ctx, identities)

	assert.True(t, store.Has(k.getIdentityStoreKey(first)))
	assert.True(t, store.Has(k.getIdentityStoreKey(second)))
	assert.True(t, store.Has(k.getIdentityStoreKey(third)))
	assert.True(t, store.Has(k.getIdentityStoreKey(fourth)))
	assert.True(t, store.Has(k.getIdentityStoreKey(fifth)))
}
