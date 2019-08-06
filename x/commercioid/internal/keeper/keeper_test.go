package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_CreateIdentity(t *testing.T) {

	identitiesStore := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.identitiesStoreKey)
	storeLen := len(identitiesStore.Get([]byte(TestOwnerIdentity)))

	TestUtils.IdKeeper.CreateIdentity(TestUtils.Ctx, TestOwner, TestOwnerIdentity, TestIdentityRef)

	afterOpLen := len(identitiesStore.Get([]byte(TestOwnerIdentity)))

	assert.NotEqual(t, storeLen, afterOpLen)
}

func TestKeeper_EditIdentity(t *testing.T) {

	updatedIdentityRef := "ddo-reference-update"

	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.identitiesStoreKey)
	store.Set([]byte(TestOwnerIdentity), []byte(TestIdentityRef))

	TestUtils.IdKeeper.EditIdentity(TestUtils.Ctx, TestOwner, TestOwnerIdentity, updatedIdentityRef)

	actual := store.Get([]byte(TestOwnerIdentity))

	assert.Equal(t, updatedIdentityRef, string(actual))

}

func TestKeeper_GetDdoReferenceByDid(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.identitiesStoreKey)
	store.Set([]byte(TestOwnerIdentity), []byte(TestIdentityRef))

	actual := TestUtils.IdKeeper.GetDdoReferenceByDid(TestUtils.Ctx, TestOwnerIdentity)

	assert.Equal(t, TestIdentityRef, actual)
}

func TestKeeper_CanBeUsedBy_UserWithNoRegisteredIdentities(t *testing.T) {

	owner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca310")
	ownerIdentity := types.Did("idid")

	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(TestIdentityRef))

	actual := TestUtils.IdKeeper.CanBeUsedBy(TestUtils.Ctx, owner, ownerIdentity)

	assert.False(t, actual)
}

func TestKeeper_CanBeUsedBy_UnregisteredIdentity(t *testing.T) {
	testOwner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca310")
	testOwnerIdentity := types.Did("did2")

	actual := TestUtils.IdKeeper.CanBeUsedBy(TestUtils.Ctx, testOwner, testOwnerIdentity)

	assert.True(t, actual)
}

func TestKeeper_CanBeUsedBy_OwnerOwnsTheGivenIdentity(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.identitiesStoreKey)
	store.Set([]byte(TestOwnerIdentity), []byte(TestIdentityRef))

	var identities = []types.Did{TestOwnerIdentity}

	ownerStore := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.ownersStoresKey)
	ownerStore.Set([]byte(TestOwner), TestUtils.IdKeeper.Cdc.MustMarshalBinaryBare(&identities))

	actual := TestUtils.IdKeeper.CanBeUsedBy(TestUtils.Ctx, TestOwner, TestOwnerIdentity)

	assert.True(t, actual)
}

func TestKeeper_AddConnection(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.connectionsStoreKey)

	var connections = []types.Did{TestOwnerIdentity, TestRecipient}

	storeLen := len(store.Get(TestUtils.Cdc.MustMarshalBinaryBare(&connections)))

	TestUtils.IdKeeper.AddConnection(TestUtils.Ctx, TestOwnerIdentity, TestRecipient)

	afterOpLen := len(store.Get(TestUtils.Cdc.MustMarshalBinaryBare(&connections)))

	assert.Equal(t, storeLen, afterOpLen)
}

func TestKeeper_GetConnections(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.connectionsStoreKey)

	var connections = []types.Did{TestOwnerIdentity}

	expected := []types.Did{TestOwnerIdentity}

	store.Set([]byte(TestOwnerIdentity), TestUtils.Cdc.MustMarshalBinaryBare(&connections))

	actual := TestUtils.IdKeeper.GetConnections(TestUtils.Ctx, TestOwnerIdentity)

	assert.Equal(t, expected, actual)
}
