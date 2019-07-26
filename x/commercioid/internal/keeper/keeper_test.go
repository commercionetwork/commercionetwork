package keeper

/*
import (
	"commercio-network/types"
	"commercio-network/x/commercioid"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_CreateIdentity(t *testing.T) {

	owner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0")
	ownerIdentity := types.Did("idid")

	identitiesStore := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.identitiesStoreKey)
	storeLen := len(identitiesStore.Get([]byte(ownerIdentity)))

	commercioid.input.idKeeper.CreateIdentity(commercioid.input.ctx, owner, ownerIdentity, commercioid.identityRef)

	afterOpLen := len(identitiesStore.Get([]byte(ownerIdentity)))

	assert.NotEqual(t, storeLen, afterOpLen)
}

func TestKeeper_GetDdoReferenceByDid(t *testing.T) {
	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.identitiesStoreKey)
	store.Set([]byte(commercioid.ownerIdentity), []byte(commercioid.identityRef))

	actual := commercioid.input.idKeeper.GetDdoReferenceByDid(commercioid.input.ctx, commercioid.ownerIdentity)

	assert.Equal(t, commercioid.identityRef, actual)
}

func TestKeeper_CanBeUsedBy_UserWithNoRegisteredIdentities(t *testing.T) {

	owner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca310")
	ownerIdentity := types.Did("idid")

	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(commercioid.identityRef))

	actual := commercioid.input.idKeeper.CanBeUsedBy(commercioid.input.ctx, owner, ownerIdentity)

	assert.False(t, actual)
}

func TestKeeper_CanBeUsedBy_UnregisteredIdentity(t *testing.T) {
	actual := commercioid.input.idKeeper.CanBeUsedBy(commercioid.input.ctx, commercioid.owner, commercioid.ownerIdentity)

	assert.True(t, actual)
}

func TestKeeper_CanBeUsedBy_OwnerOwnsTheGivenIdentity(t *testing.T) {
	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.identitiesStoreKey)
	store.Set([]byte(commercioid.ownerIdentity), []byte(commercioid.identityRef))

	var identities = []types.Did{commercioid.ownerIdentity}

	ownerStore := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.ownersStoresKey)
	ownerStore.Set([]byte(commercioid.owner), commercioid.input.idKeeper.Cdc.MustMarshalBinaryBare(&identities))

	actual := commercioid.input.idKeeper.CanBeUsedBy(commercioid.input.ctx, commercioid.owner, commercioid.ownerIdentity)

	assert.True(t, actual)
}

func TestKeeper_AddConnection(t *testing.T) {
	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.connectionsStoreKey)

	var connections = []types.Did{commercioid.ownerIdentity, commercioid.recipient}

	storeLen := len(store.Get(commercioid.input.cdc.MustMarshalBinaryBare(&connections)))

	commercioid.input.idKeeper.AddConnection(commercioid.input.ctx, commercioid.ownerIdentity, commercioid.recipient)

	afterOpLen := len(store.Get(commercioid.input.cdc.MustMarshalBinaryBare(&connections)))

	assert.Equal(t, storeLen, afterOpLen)
}

func TestKeeper_GetConnections(t *testing.T) {
	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.connectionsStoreKey)

	var connections = []types.Did{commercioid.ownerIdentity}

	expected := []types.Did{commercioid.ownerIdentity}

	store.Set([]byte(commercioid.ownerIdentity), commercioid.input.cdc.MustMarshalBinaryBare(&connections))

	actual := commercioid.input.idKeeper.GetConnections(commercioid.input.ctx, commercioid.ownerIdentity)

	assert.Equal(t, expected, actual)
}


*/
