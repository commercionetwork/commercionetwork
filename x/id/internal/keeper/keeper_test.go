package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_CreateIdentity(t *testing.T) {
	identitiesStore := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	storeData := len(identitiesStore.Get(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress)))
	assert.Equal(t, 0, storeData)

	TestUtils.IdKeeper.SaveIdentity(TestUtils.Ctx, TestOwnerAddress, TestDidDocumentUri)

	afterOpLen := len(identitiesStore.Get(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress)))
	assert.Equal(t, len(TestDidDocumentUri), afterOpLen)
}

func TestKeeper_EditIdentity(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress), []byte(TestDidDocumentUri))
	storeData := store.Get(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress))
	assert.Equal(t, []byte(TestDidDocumentUri), storeData)

	updatedIdentityRef := "ddo-reference-update"
	TestUtils.IdKeeper.SaveIdentity(TestUtils.Ctx, TestOwnerAddress, updatedIdentityRef)

	updatedLen := store.Get(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress))
	assert.Equal(t, []byte(updatedIdentityRef), updatedLen)
}

func TestKeeper_GetDidDocumentUriByDid(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress), []byte(TestDidDocumentUri))
	storeData := store.Get(TestUtils.IdKeeper.getIdentiyStoreKey(TestOwnerAddress))
	assert.Equal(t, []byte(TestDidDocumentUri), storeData)

	actual := TestUtils.IdKeeper.GetDidDocumentUriByDid(TestUtils.Ctx, TestOwnerAddress)
	assert.Equal(t, TestDidDocumentUri, actual)
}

// -------------------------
// --- Genesis utils
// -------------------------

func TestKeeper_GetIdentities(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	first, err := sdk.AccAddressFromBech32("cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second, err := sdk.AccAddressFromBech32("cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third, err := sdk.AccAddressFromBech32("cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth, err := sdk.AccAddressFromBech32("cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth, err := sdk.AccAddressFromBech32("cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")
	assert.Nil(t, err)

	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(first), []byte("first"))
	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(second), []byte("second"))
	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(third), []byte("third"))
	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(fourth), []byte("fourth"))
	store.Set(TestUtils.IdKeeper.getIdentiyStoreKey(fifth), []byte("fifth"))

	actual, err := TestUtils.IdKeeper.GetIdentities(TestUtils.Ctx)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(actual))
	assert.Contains(t, actual, types.Identity{Owner: first, DidDocument: "first"})
	assert.Contains(t, actual, types.Identity{Owner: second, DidDocument: "second"})
	assert.Contains(t, actual, types.Identity{Owner: third, DidDocument: "third"})
	assert.Contains(t, actual, types.Identity{Owner: fourth, DidDocument: "fourth"})
	assert.Contains(t, actual, types.Identity{Owner: fifth, DidDocument: "fifth"})
}

func TestKeeper_SetIdentities(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	first, err := sdk.AccAddressFromBech32("cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second, err := sdk.AccAddressFromBech32("cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third, err := sdk.AccAddressFromBech32("cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth, err := sdk.AccAddressFromBech32("cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth, err := sdk.AccAddressFromBech32("cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")
	assert.Nil(t, err)

	identities := []types.Identity{
		{Owner: first, DidDocument: "first"},
		{Owner: second, DidDocument: "second"},
		{Owner: third, DidDocument: "third"},
		{Owner: fourth, DidDocument: "fourth"},
		{Owner: fifth, DidDocument: "fifth"},
	}
	TestUtils.IdKeeper.SetIdentities(TestUtils.Ctx, identities)

	assert.True(t, store.Has(TestUtils.IdKeeper.getIdentiyStoreKey(first)))
	assert.True(t, store.Has(TestUtils.IdKeeper.getIdentiyStoreKey(second)))
	assert.True(t, store.Has(TestUtils.IdKeeper.getIdentiyStoreKey(third)))
	assert.True(t, store.Has(TestUtils.IdKeeper.getIdentiyStoreKey(fourth)))
	assert.True(t, store.Has(TestUtils.IdKeeper.getIdentiyStoreKey(fifth)))
}
