package keeper

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commercioid"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	db2 "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

//testing add readers
func TestKeeper_addReaderForDocument(t *testing.T) {

	var readers = []types.Did{"reader", "reader2"}

	store := testUtils.ctx.KVStore(testUtils.docsKeeper.readersStoreKey)
	store.Set([]byte(testReference), testUtils.cdc.MustMarshalBinaryBare(&readers))

	currentLength := len(store.Get([]byte(testReference)))

	testUtils.docsKeeper.addReaderForDocument(testUtils.ctx, testOwnerIdentity, testReference)

	afterOpLength := len(store.Get([]byte(testReference)))

	if afterOpLength < currentLength {
		t.Errorf("afterOpLength should be greater than currentLength")
	}
}

func TestKeeper_StoreDocument(t *testing.T) {

	ownerStore := testUtils.ctx.KVStore(testUtils.docsKeeper.ownersStoreKey)
	ownerStore.Set([]byte(testReference), []byte(testOwner))

	metadataStore := testUtils.ctx.KVStore(testUtils.docsKeeper.metadataStoreKey)

	currentLength := len(metadataStore.Get([]byte(testReference)))

	testUtils.docsKeeper.StoreDocument(testUtils.ctx, testOwner, testOwnerIdentity, testReference, testMetadata)

	afterOpLength := len(metadataStore.Get([]byte(testReference)))

	if afterOpLength < currentLength {
		t.Errorf("after operation length should be greater than current length")
	}
}

//Given testReference has an testOwner
func TestKeeper_HasOwner_True(t *testing.T) {

	store := testUtils.ctx.KVStore(testUtils.docsKeeper.ownersStoreKey)
	store.Set([]byte(testReference), []byte(testOwner))

	result := testUtils.docsKeeper.HasOwner(testUtils.ctx, testReference)

	assert.True(t, result)
}

//Given testReference hasn't got an testOwner
func TestKeeper_HasOwner_False(t *testing.T) {

	reference := "reff"

	result := testUtils.docsKeeper.HasOwner(testUtils.ctx, reference)

	assert.False(t, result)
}

//Given testOwner is the testOwner of doc testReference
func TestKeeper_IsOwner_True(t *testing.T) {

	store := testUtils.ctx.KVStore(testUtils.docsKeeper.ownersStoreKey)
	store.Set([]byte(testReference), []byte(testOwner))

	res := testUtils.docsKeeper.IsOwner(testUtils.ctx, testOwner, testReference)

	assert.True(t, res)
}

//Given testOwner isnt the testOwner of the doc testReference
func TestKeeper_IsOwner_False(t *testing.T) {

	reference := "reff"
	res := testUtils.docsKeeper.IsOwner(testUtils.ctx, testOwner, reference)

	assert.False(t, res)
}

func TestKeeper_GetMetadata_OfExistentDocument(t *testing.T) {

	metadataStore := testUtils.ctx.KVStore(testUtils.docsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(testReference), []byte(testMetadata))

	result := testUtils.docsKeeper.GetMetadata(testUtils.ctx, testReference)

	assert.Equal(t, testMetadata, result)
}

func TestKeeper_GetMetadata_OfNonExistentDocument(t *testing.T) {

	reference := "reff"

	result := testUtils.docsKeeper.GetMetadata(testUtils.ctx, reference)

	assert.Equal(t, "", result)
}

func TestKeeper_CanReadDocument_True(t *testing.T) {

	readers := []types.Did{testOwnerIdentity}

	readerStore := testUtils.ctx.KVStore(testUtils.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(testReference), testUtils.cdc.MustMarshalBinaryBare(&readers))

	result := testUtils.docsKeeper.CanReadDocument(testUtils.ctx, testOwnerIdentity, testReference)

	assert.True(t, result)
}

func TestKeeper_CanReadDocument_False(t *testing.T) {

	reference := "reff"

	result := testUtils.docsKeeper.CanReadDocument(testUtils.ctx, testOwnerIdentity, reference)

	assert.False(t, result)
}

func TestKeeper_GetAuthorizedReaders(t *testing.T) {
	var readers = []types.Did{"reader", "reader2"}

	store := testUtils.ctx.KVStore(testUtils.docsKeeper.readersStoreKey)
	store.Set([]byte(testReference), testUtils.cdc.MustMarshalBinaryBare(&readers))

	res := testUtils.docsKeeper.GetAuthorizedReaders(testUtils.ctx, testReference)

	assert.Equal(t, readers, res)
}

func TestKeeper_ShareDocument_SenderAuthorizedToShare(t *testing.T) {

	var readers = []types.Did{testOwnerIdentity}

	readerStore := testUtils.ctx.KVStore(testUtils.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(testReference), testUtils.cdc.MustMarshalBinaryBare(&readers))

	result := testUtils.docsKeeper.ShareDocument(testUtils.ctx, testReference, testOwnerIdentity, testRecipient)

	assert.Nil(t, result)
}

func TestKeeper_ShareDocument_SenderUnauthorizedToShare(t *testing.T) {

	ownerIdentity := types.Did("notOwner")
	error := sdk.ErrUnauthorized(fmt.Sprintf("The sender with testAddress %s doesnt have the rights on this document", ownerIdentity))

	result := testUtils.docsKeeper.ShareDocument(testUtils.ctx, testReference, ownerIdentity, testRecipient)

	assert.NotNil(t, result)

	assert.Equal(t, error, result)
}

type testInput struct {
	cdc        *codec.Codec
	ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.BaseKeeper
	docsKeeper Keeper
}

//This function create an enviroment to test modules
func setupTestInput() testInput {

	db := db2.NewMemDB()
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

	// CommercioDOCS
	keyDOCSOwners := sdk.NewKVStoreKey("docs_owners")
	keyDOCSMetadata := sdk.NewKVStoreKey("docs_metadata")
	keyDOCSSharing := sdk.NewKVStoreKey("docs_sharing")
	keyDOCSReaders := sdk.NewKVStoreKey("docs_readers")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyDOCSReaders, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDOCSOwners, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDOCSMetadata, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDOCSSharing, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDIdentities, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDOwners, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDConnections, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	idk := commercioid.NewKeeper(keyIDIdentities, keyIDOwners, keyIDConnections, cdc)
	dck := NewKeeper(idk, keyDOCSOwners, keyDOCSMetadata, keyDOCSSharing, keyDOCSReaders, cdc)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{
		cdc:        cdc,
		ctx:        ctx,
		accKeeper:  ak,
		bankKeeper: bk,
		docsKeeper: dck,
	}

}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}
