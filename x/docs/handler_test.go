package docs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testUtils = keeper.TestUtils
var handler = NewHandler(testUtils.DocsKeeper)

// -----------------------------
// --- handleMsgShareDocument
// -----------------------------

func Test_handleMsgShareDocument_CustomMetadataSpecs(t *testing.T) {
	msgShareDocument := MsgShareDocument(keeper.TestingDocument)

	res := handler(testUtils.Ctx, msgShareDocument)
	require.True(t, res.IsOK())
}

func Test_handleMsgShareDocument_MetadataSchemeType_Supported(t *testing.T) {
	msgShareDocument := MsgShareDocument{
		Sender:    keeper.TestingDocument.Sender,
		Recipient: keeper.TestingDocument.Recipient,
		Uuid:      keeper.TestingDocument.Uuid,
		Metadata: types.DocumentMetadata{
			ContentUri: keeper.TestingDocument.Metadata.ContentUri,
			SchemaType: "metadata-schema",
			Proof:      keeper.TestingDocument.Metadata.Proof,
		},
		ContentUri: keeper.TestingDocument.ContentUri,
		Checksum:   keeper.TestingDocument.Checksum,
	}
	supportedSchema := types.MetadataSchema{
		Type:      "metadata-schema",
		SchemaUri: "https://example.com/schema",
		Version:   "1.0.0",
	}
	testUtils.DocsKeeper.AddSupportedMetadataScheme(testUtils.Ctx, supportedSchema)

	res := handler(testUtils.Ctx, msgShareDocument)
	assert.True(t, res.IsOK())
}

func Test_handleMsgShareDocument_MetadataSchemeType_NotSupported(t *testing.T) {
	msgShareDocument := MsgShareDocument{
		Sender:    keeper.TestingDocument.Sender,
		Recipient: keeper.TestingDocument.Recipient,
		Uuid:      keeper.TestingDocument.Uuid,
		Metadata: types.DocumentMetadata{
			ContentUri: keeper.TestingDocument.Metadata.ContentUri,
			SchemaType: "non-existent-schema-type",
			Proof:      keeper.TestingDocument.Metadata.Proof,
		},
		ContentUri: keeper.TestingDocument.ContentUri,
		Checksum:   keeper.TestingDocument.Checksum,
	}

	res := handler(testUtils.Ctx, msgShareDocument)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

// -----------------------------------
// --- handleMsgSendDocumentReceipt
// -----------------------------------

func Test_handleMsgSendDocumentReceipt(t *testing.T) {
	msgSendDocumentReceipt := MsgSendDocumentReceipt(keeper.TestingDocumentReceipt)

	res := handler(testUtils.Ctx, msgSendDocumentReceipt)
	assert.True(t, res.IsOK())
}

// -------------------------------------------
// --- handleMsgAddSupportedMetadataSchema
// -------------------------------------------

func Test_handleMsgAddSupportedMetadataSchema_NotTrustedSigner(t *testing.T) {
	store := testUtils.Ctx.KVStore(testUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.MetadataSchemaProposersStoreKey))

	msgAddSupportedMetadataSchema := MsgAddSupportedMetadataSchema{
		Signer: keeper.TestingSender,
		Schema: types.MetadataSchema{
			Type:      "schema-type",
			SchemaUri: "https://example.com/schema",
			Version:   "1.0.0",
		},
	}
	res := handler(testUtils.Ctx, msgAddSupportedMetadataSchema)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddSupportedMetadataSchema_TrustedSigner(t *testing.T) {
	signers := []sdk.AccAddress{keeper.TestingSender}

	store := testUtils.Ctx.KVStore(testUtils.DocsKeeper.StoreKey)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), testUtils.Cdc.MustMarshalBinaryBare(&signers))

	msgAddSupportedMetadataSchema := MsgAddSupportedMetadataSchema{
		Signer: keeper.TestingSender,
		Schema: types.MetadataSchema{
			Type:      "schema-type",
			SchemaUri: "https://example.com/schema",
			Version:   "1.0.0",
		},
	}
	res := handler(testUtils.Ctx, msgAddSupportedMetadataSchema)
	assert.True(t, res.IsOK())
}

// ------------------------------------------------
// --- handleMsgAddTrustedMetadataSchemaProposer
// ------------------------------------------------

func Test_handleMsgAddTrustedMetadataSchemaProposer_MissingGovernment(t *testing.T) {
	store := testUtils.Ctx.KVStore(testUtils.DocsKeeper.GovernmentKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: keeper.TestingSender,
		Signer:   keeper.TestingRecipient,
	}
	res := handler(testUtils.Ctx, msgAddTrustedMetadataSchemaProposer)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddTrustedMetadataSchemaProposer_IncorrectSigner(t *testing.T) {
	testUtils.DocsKeeper.GovernmentKeeper.SetGovernmentAddress(testUtils.Ctx, keeper.TestingRecipient)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: keeper.TestingSender,
		Signer:   keeper.TestingSender,
	}
	res := handler(testUtils.Ctx, msgAddTrustedMetadataSchemaProposer)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddTrustedMetadataSchemaProposer_CorrectSigner(t *testing.T) {
	testUtils.DocsKeeper.GovernmentKeeper.SetGovernmentAddress(testUtils.Ctx, keeper.TestingRecipient)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: keeper.TestingSender,
		Signer:   keeper.TestingRecipient,
	}
	res := handler(testUtils.Ctx, msgAddTrustedMetadataSchemaProposer)
	assert.True(t, res.IsOK())
}

// -------------------
// --- Default case
// -------------------

func Test_invalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	assert.False(t, res.IsOK())
	assert.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
