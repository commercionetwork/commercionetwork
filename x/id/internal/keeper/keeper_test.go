package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_CreateIdentity(t *testing.T) {
	identitiesStore := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	storeData := len(identitiesStore.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String())))
	assert.Equal(t, 0, storeData)

	TestUtils.IdKeeper.SaveIdentity(TestUtils.Ctx, TestOwnerAddress, TestDidDocumentUri)

	afterOpLen := len(identitiesStore.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String())))
	assert.Equal(t, len(TestDidDocumentUri), afterOpLen)
}

func TestKeeper_EditIdentity(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set([]byte(types.IdentitiesStorePrefix+TestOwnerAddress.String()), []byte(TestDidDocumentUri))
	storeData := store.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String()))
	assert.Equal(t, []byte(TestDidDocumentUri), storeData)

	updatedIdentityRef := "ddo-reference-update"
	TestUtils.IdKeeper.SaveIdentity(TestUtils.Ctx, TestOwnerAddress, updatedIdentityRef)

	updatedLen := store.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String()))
	assert.Equal(t, []byte(updatedIdentityRef), updatedLen)
}

func TestKeeper_GetDidDocumentUriByDid(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set([]byte(types.IdentitiesStorePrefix+TestOwnerAddress.String()), []byte(TestDidDocumentUri))
	storeData := store.Get([]byte(types.IdentitiesStorePrefix + TestOwnerAddress.String()))
	assert.Equal(t, []byte(TestDidDocumentUri), storeData)

	actual := TestUtils.IdKeeper.GetDidDocumentUriByDid(TestUtils.Ctx, TestOwnerAddress)
	assert.Equal(t, TestDidDocumentUri, actual)
}
