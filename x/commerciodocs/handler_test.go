package commerciodocs

/*
import (
	"commercio-network/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var handler = NewHandler(input.docsKeeper)

func Test_handleStoreDocument_documentHasAlreadyAnOwner(t *testing.T) {

	var address = "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0"
	var owner, _ = sdk.AccAddressFromBech32(address)

	docStore := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	docStore.Set([]byte(reference), owner)

	res := handler(input.ctx, msgStore)

	expected := sdk.ErrUnauthorized("The given account has no access to the document").Result()

	assert.Equal(t, expected, res)

}

func Test_handleStoreDocument_documentStoredCorrectly(t *testing.T) {
	docStore := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	docStore.Set([]byte(reference), owner)

	res := handler(input.ctx, msgStore)

	assert.Equal(t, sdk.Result{}, res)
}

func Test_handleShareDocument_documentHasAlreadyAnOwner(t *testing.T) {

	var address = "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0"
	var owner, _ = sdk.AccAddressFromBech32(address)

	docStore := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	docStore.Set([]byte(reference), owner)

	res := handler(input.ctx, msgShare)

	expected := sdk.ErrUnauthorized("The given account has no access to the document").Result()

	assert.Equal(t, expected, res)

}

func Test_handlerShareDocument_documentSharedCorrectly(t *testing.T) {

	docStore := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	docStore.Set([]byte(reference), owner)

	readersStore := input.ctx.KVStore(input.docsKeeper.readersStoreKey)
	var readers = []types.Did{msgShare.Sender}

	readersStore.Set([]byte(reference), input.docsKeeper.cdc.MustMarshalBinaryBare(&readers))

	res := handler(input.ctx, msgShare)

	assert.Equal(t, sdk.Result{}, res)
}

*/
