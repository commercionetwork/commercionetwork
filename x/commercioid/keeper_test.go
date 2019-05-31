package commercioid

import (
	"commercio-network/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_CreateIdentity(t *testing.T) {
	identitiesStore := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	storeLen := len(identitiesStore.Get([]byte(ownerIdentity)))

	input.idKeeper.CreateIdentity(input.ctx, owner, ownerIdentity, identityRef)

	afterOpLen := len(identitiesStore.Get([]byte(ownerIdentity)))

	assert.NotEqual(t, storeLen, afterOpLen)
}

func TestKeeper_GetDdoReferenceByDid(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(identityRef))

	actual := input.idKeeper.GetDdoReferenceByDid(input.ctx, ownerIdentity)

	assert.Equal(t, identityRef, actual)
}

func TestKeeper_CanBeUsedBy_UserWithNoRegisteredIdentities(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(identityRef))

	actual := input.idKeeper.CanBeUsedBy(input.ctx, owner, ownerIdentity)

	assert.False(t, actual)
}

func TestKeeper_CanBeUsedBy_UnregisteredIdentity(t *testing.T) {
	actual := input.idKeeper.CanBeUsedBy(input.ctx, owner, ownerIdentity)

	assert.True(t, actual)
}

func TestKeeper_CanBeUsedBy_OwnerOwnsTheGivenIdentity(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(identityRef))

	var identities = []types.Did{ownerIdentity}

	ownerStore := input.ctx.KVStore(input.idKeeper.ownersStoresKey)
	ownerStore.Set([]byte(owner), input.idKeeper.cdc.MustMarshalBinaryBare(&identities))

	actual := input.idKeeper.CanBeUsedBy(input.ctx, owner, ownerIdentity)

	assert.True(t, actual)
}

func TestKeeper_AddConnection(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.connectionsStoreKey)

	var connections = []types.Did{ownerIdentity, recipient}

	storeLen := len(store.Get(input.cdc.MustMarshalBinaryBare(&connections)))

	input.idKeeper.AddConnection(input.ctx, ownerIdentity, recipient)

	afterOpLen := len(store.Get(input.cdc.MustMarshalBinaryBare(&connections)))

	assert.Equal(t, storeLen, afterOpLen)
}

func TestKeeper_GetConnections(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.connectionsStoreKey)

	var connections []types.Did

	expected := []types.Did{ownerIdentity}

	store.Set([]byte(ownerIdentity), input.cdc.MustMarshalBinaryBare(&connections))

	actual := input.idKeeper.GetConnections(input.ctx, ownerIdentity)

	assert.Equal(t, expected, actual)
}
