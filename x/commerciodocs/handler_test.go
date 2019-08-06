package commerciodocs

import (
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var msgStoreDoc = MsgStoreDoc{
	Identity:  keeper.TestOwnerIdentity,
	Reference: keeper.TestReference,
	Owner:     keeper.TestOwner,
	Metadata:  keeper.TestMetadata,
}

var msgShareDoc = MsgShareDoc{
	Owner:     keeper.TestOwner,
	Sender:    keeper.TestOwnerIdentity,
	Receiver:  keeper.TestRecipient,
	Reference: keeper.TestReference,
}

var testUtils = keeper.TestUtils
var handler = NewHandler(testUtils.DocsKeeper)

func TestValidMsg_StoreDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgStoreDoc)

	require.True(t, res.IsOK())
}

func TestValidMsg_ShareDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgShareDoc)

	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized commerciodocs message type"))
}

//TODO
// What to do with this tests? We cant access to keeper's private fields so IMO i will delete them
/*
func Test_handleStoreDocument_documentHasAlreadyAnOwner(t *testing.T) {

	var address = "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0"
	var owner, _ = sdk.AccAddressFromBech32(address)

	docStore := testUtils.Ctx.KVStore(testUtils.DocsKeeper.ownersStoreKey)
	docStore.Set([]byte(keeper.TestReference), owner)

	res := handler(testUtils.Ctx, MsgStoreDoc{})

	expected := sdk.ErrUnauthorized("The given account has no access to the document").Result()

	assert.Equal(t, expected, res)

}

func Test_handleStoreDocument_documentStoredCorrectly(t *testing.T) {
	docStore := testUtils.Ctx.KVStore(testUtils.DocsKeeper.ownersStoreKey)
	docStore.Set([]byte(keeper.TestReference), keeper.TestOwner)

	res := handler(testUtils.Ctx, MsgStoreDoc{})

	assert.Equal(t, sdk.Result{}, res)
}

func Test_handleShareDocument_documentHasAlreadyAnOwner(t *testing.T) {

	var address = "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0"
	var owner, _ = sdk.AccAddressFromBech32(address)

	docStore := testUtils.Ctx.KVStore(testUtils.DocsKeeper.ownersStoreKey)
	docStore.Set([]byte(keeper.TestReference), owner)

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
