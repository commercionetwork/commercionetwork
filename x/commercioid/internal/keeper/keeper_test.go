package keeper

import (
	"commercio-network/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func TestKeeper_CreateIdentity(t *testing.T) {

	owner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0")
	ownerIdentity := types.Did("idid")

	identitiesStore := testUtils.ctx.KVStore(testUtils.idKeeper.identitiesStoreKey)
	storeLen := len(identitiesStore.Get([]byte(ownerIdentity)))

	testUtils.idKeeper.CreateIdentity(testUtils.ctx, owner, ownerIdentity, testIdentityRef)

	afterOpLen := len(identitiesStore.Get([]byte(ownerIdentity)))

	assert.NotEqual(t, storeLen, afterOpLen)
}

func TestKeeper_GetDdoReferenceByDid(t *testing.T) {
	store := testUtils.ctx.KVStore(testUtils.idKeeper.identitiesStoreKey)
	store.Set([]byte(testOwnerIdentity), []byte(testIdentityRef))

	actual := testUtils.idKeeper.GetDdoReferenceByDid(testUtils.ctx, testOwnerIdentity)

	assert.Equal(t, testIdentityRef, actual)
}

func TestKeeper_CanBeUsedBy_UserWithNoRegisteredIdentities(t *testing.T) {

	owner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca310")
	ownerIdentity := types.Did("idid")

	store := testUtils.ctx.KVStore(testUtils.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(testIdentityRef))

	actual := testUtils.idKeeper.CanBeUsedBy(testUtils.ctx, owner, ownerIdentity)

	assert.False(t, actual)
}

func TestKeeper_CanBeUsedBy_UnregisteredIdentity(t *testing.T) {
	testOwner, _ := sdk.AccAddressFromBech32("cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca310")
	testOwnerIdentity := types.Did("idid")

	actual := testUtils.idKeeper.CanBeUsedBy(testUtils.ctx, testOwner, testOwnerIdentity)

	assert.True(t, actual)
}

func TestKeeper_CanBeUsedBy_OwnerOwnsTheGivenIdentity(t *testing.T) {
	store := testUtils.ctx.KVStore(testUtils.idKeeper.identitiesStoreKey)
	store.Set([]byte(testOwnerIdentity), []byte(testIdentityRef))

	var identities = []types.Did{testOwnerIdentity}

	ownerStore := testUtils.ctx.KVStore(testUtils.idKeeper.ownersStoresKey)
	ownerStore.Set([]byte(testOwner), testUtils.idKeeper.Cdc.MustMarshalBinaryBare(&identities))

	actual := testUtils.idKeeper.CanBeUsedBy(testUtils.ctx, testOwner, testOwnerIdentity)

	assert.True(t, actual)
}

func TestKeeper_AddConnection(t *testing.T) {
	store := testUtils.ctx.KVStore(testUtils.idKeeper.connectionsStoreKey)

	var connections = []types.Did{testOwnerIdentity, testRecipient}

	storeLen := len(store.Get(testUtils.cdc.MustMarshalBinaryBare(&connections)))

	testUtils.idKeeper.AddConnection(testUtils.ctx, testOwnerIdentity, testRecipient)

	afterOpLen := len(store.Get(testUtils.cdc.MustMarshalBinaryBare(&connections)))

	assert.Equal(t, storeLen, afterOpLen)
}

func TestKeeper_GetConnections(t *testing.T) {
	store := testUtils.ctx.KVStore(testUtils.idKeeper.connectionsStoreKey)

	var connections = []types.Did{testOwnerIdentity}

	expected := []types.Did{testOwnerIdentity}

	store.Set([]byte(testOwnerIdentity), testUtils.cdc.MustMarshalBinaryBare(&connections))

	actual := testUtils.idKeeper.GetConnections(testUtils.ctx, testOwnerIdentity)

	assert.Equal(t, expected, actual)
}

type testInput struct {
	cdc        *codec.Codec
	ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.BaseKeeper
	idKeeper   Keeper
}

//This function create an enviroment to test modules
func setupTestInput() testInput {

	db := db.NewMemDB()
	cdc := testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// CommercioID
	keyIDIdentities := sdk.NewKVStoreKey("id_identities")
	keyIDOwners := sdk.NewKVStoreKey("id_owners")
	keyIDConnections := sdk.NewKVStoreKey("id_connections")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyIDIdentities, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDOwners, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDConnections, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	idk := NewKeeper(keyIDIdentities, keyIDOwners, keyIDConnections, cdc)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{
		cdc:        cdc,
		ctx:        ctx,
		accKeeper:  ak,
		bankKeeper: bk,
		idKeeper:   idk,
	}

}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}
