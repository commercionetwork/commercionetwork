package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_CreateIdentity(t *testing.T) {
	identitiesStore := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	storeData := len(identitiesStore.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String())))
	assert.Equal(t, 0, storeData)

	TestUtils.IdKeeper.SaveIdentity(TestUtils.Ctx, TestOwnerAddress, TestDidDocumentReference)

	afterOpLen := len(identitiesStore.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String())))
	assert.Equal(t, len(TestDidDocumentReference), afterOpLen)
}

func TestKeeper_EditIdentity(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set([]byte(types.IdentitiesStorePrefix+TestOwnerAddress.String()), []byte(TestDidDocumentReference))
	storeData := store.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String()))
	assert.Equal(t, []byte(TestDidDocumentReference), storeData)

	updatedIdentityRef := "ddo-reference-update"
	TestUtils.IdKeeper.SaveIdentity(TestUtils.Ctx, TestOwnerAddress, updatedIdentityRef)

	updatedLen := store.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String()))
	assert.Equal(t, []byte(updatedIdentityRef), updatedLen)
}

func TestKeeper_GetDidDocumentReferenceByDid(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set([]byte(types.IdentitiesStorePrefix+TestOwnerAddress.String()), []byte(TestDidDocumentReference))
	storeData := store.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String()))
	assert.Equal(t, []byte(TestDidDocumentReference), storeData)

	actual := TestUtils.IdKeeper.GetDidDocumentReferenceByDid(TestUtils.Ctx, TestOwnerAddress)
	assert.Equal(t, TestDidDocumentReference, actual)
}
