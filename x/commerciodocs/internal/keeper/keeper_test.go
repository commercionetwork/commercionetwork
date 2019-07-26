package keeper

/*
import (
	"commercio-network/types"
	"commercio-network/x/commerciodocs"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//testing add readers
func TestKeeper_addReaderForDocument(t *testing.T) {

	var readers = []types.Did{"reader", "reader2"}

	store := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.readersStoreKey)
	store.Set([]byte(commerciodocs.reference), commerciodocs.input.cdc.MustMarshalBinaryBare(&readers))

	currentLength := len(store.Get([]byte(commerciodocs.reference)))

	commerciodocs.input.docsKeeper.addReaderForDocument(commerciodocs.input.ctx, commerciodocs.ownerIdentity, commerciodocs.reference)

	afterOpLength := len(store.Get([]byte(commerciodocs.reference)))

	if afterOpLength < currentLength {
		t.Errorf("afterOpLength should be greater than currentLength")
	}
}

func TestKeeper_StoreDocument(t *testing.T) {

	ownerStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.ownersStoreKey)
	ownerStore.Set([]byte(commerciodocs.reference), []byte(commerciodocs.owner))

	metadataStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.metadataStoreKey)

	currentLength := len(metadataStore.Get([]byte(commerciodocs.reference)))

	commerciodocs.input.docsKeeper.StoreDocument(commerciodocs.input.ctx, commerciodocs.owner, commerciodocs.ownerIdentity, commerciodocs.reference, commerciodocs.metadata)

	afterOpLength := len(metadataStore.Get([]byte(commerciodocs.reference)))

	if afterOpLength < currentLength {
		t.Errorf("after operation length should be greater than current length")
	}
}

//Given reference has an owner
func TestKeeper_HasOwner_True(t *testing.T) {

	store := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.ownersStoreKey)
	store.Set([]byte(commerciodocs.reference), []byte(commerciodocs.owner))

	result := commerciodocs.input.docsKeeper.HasOwner(commerciodocs.input.ctx, commerciodocs.reference)

	assert.True(t, result)
}

//Given reference hasn't got an owner
func TestKeeper_HasOwner_False(t *testing.T) {

	reference := "reff"

	result := commerciodocs.input.docsKeeper.HasOwner(commerciodocs.input.ctx, reference)

	assert.False(t, result)
}

//Given owner is the owner of doc reference
func TestKeeper_IsOwner_True(t *testing.T) {

	store := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.ownersStoreKey)
	store.Set([]byte(commerciodocs.reference), []byte(commerciodocs.owner))

	res := commerciodocs.input.docsKeeper.IsOwner(commerciodocs.input.ctx, commerciodocs.owner, commerciodocs.reference)

	assert.True(t, res)
}

//Given owner isnt the owner of the doc reference
func TestKeeper_IsOwner_False(t *testing.T) {

	reference := "reff"
	res := commerciodocs.input.docsKeeper.IsOwner(commerciodocs.input.ctx, commerciodocs.owner, reference)

	assert.False(t, res)
}

func TestKeeper_GetMetadata_OfExistentDocument(t *testing.T) {

	metadataStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(commerciodocs.reference), []byte(commerciodocs.metadata))

	result := commerciodocs.input.docsKeeper.GetMetadata(commerciodocs.input.ctx, commerciodocs.reference)

	assert.Equal(t, commerciodocs.metadata, result)
}

func TestKeeper_GetMetadata_OfNonExistentDocument(t *testing.T) {

	reference := "reff"

	result := commerciodocs.input.docsKeeper.GetMetadata(commerciodocs.input.ctx, reference)

	assert.Equal(t, "", result)
}

func TestKeeper_CanReadDocument_True(t *testing.T) {

	readers := []types.Did{commerciodocs.ownerIdentity}

	readerStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(commerciodocs.reference), commerciodocs.input.cdc.MustMarshalBinaryBare(&readers))

	result := commerciodocs.input.docsKeeper.CanReadDocument(commerciodocs.input.ctx, commerciodocs.ownerIdentity, commerciodocs.reference)

	assert.True(t, result)
}

func TestKeeper_CanReadDocument_False(t *testing.T) {

	reference := "reff"

	result := commerciodocs.input.docsKeeper.CanReadDocument(commerciodocs.input.ctx, commerciodocs.ownerIdentity, reference)

	assert.False(t, result)
}

func TestKeeper_GetAuthorizedReaders(t *testing.T) {
	var readers = []types.Did{"reader", "reader2"}

	store := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.readersStoreKey)
	store.Set([]byte(commerciodocs.reference), commerciodocs.input.cdc.MustMarshalBinaryBare(&readers))

	res := commerciodocs.input.docsKeeper.GetAuthorizedReaders(commerciodocs.input.ctx, commerciodocs.reference)

	assert.Equal(t, readers, res)
}

func TestKeeper_ShareDocument_SenderAuthorizedToShare(t *testing.T) {

	var readers = []types.Did{commerciodocs.ownerIdentity}

	readerStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(commerciodocs.reference), commerciodocs.input.cdc.MustMarshalBinaryBare(&readers))

	result := commerciodocs.input.docsKeeper.ShareDocument(commerciodocs.input.ctx, commerciodocs.reference, commerciodocs.ownerIdentity, commerciodocs.recipient)

	assert.Nil(t, result)
}

func TestKeeper_ShareDocument_SenderUnauthorizedToShare(t *testing.T) {

	ownerIdentity := types.Did("notOwner")
	error := sdk.ErrUnauthorized(fmt.Sprintf("The sender with address %s doesnt have the rights on this document", ownerIdentity))

	result := commerciodocs.input.docsKeeper.ShareDocument(commerciodocs.input.ctx, commerciodocs.reference, ownerIdentity, commerciodocs.recipient)

	assert.NotNil(t, result)

	assert.Equal(t, error, result)
}

*/
