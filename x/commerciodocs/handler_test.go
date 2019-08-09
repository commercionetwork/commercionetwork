package commerciodocs

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestConnectionAddress, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var msgShareDocument = MsgShareDocument{
	Document: types.Document{
		Sender:     TestOwnerAddress,
		Recipient:  TestConnectionAddress,
		ContentUri: "https://example.com/document",
		Metadata: types.DocumentMetadata{
			ContentUri: "",
			Schema: types.DocumentMetadataSchema{
				Uri:     "https://example.com/document/metadata/schema",
				Version: "1.0.0",
			},
			Proof: "73666c68676c7366676c7366676c6a6873666c6a6768",
		},
		Checksum: types.DocumentChecksum{
			Value:     "93dfcaf3d923ec47edb8580667473987",
			Algorithm: "md5",
		},
	},
}

var testUtils = keeper.TestUtils
var handler = NewHandler(testUtils.DocsKeeper)

func TestValidMsg_ShareDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgShareDocument)
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
