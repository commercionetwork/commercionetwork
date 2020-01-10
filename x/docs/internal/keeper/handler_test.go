package keeper

import (
	"fmt"
	"strings"
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// -----------------------------
// --- handleMsgShareDocument
// -----------------------------

func Test_handleMsgShareDocument_CustomMetadataSpecs(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgShareDocument := types.NewMsgShareDocument(TestingDocument)

	res := handler(ctx, msgShareDocument)
	assert.True(t, res.IsOK())
}

func Test_handleMsgShareDocument_MetadataSchemeType_Supported(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgShareDocument := types.MsgShareDocument(types.Document{
		Sender:     TestingSender,
		Recipients: ctypes.Addresses{TestingRecipient},
		UUID:       TestingDocument.UUID,
		Metadata: types.DocumentMetadata{
			ContentURI: TestingDocument.Metadata.ContentURI,
			SchemaType: "metadata-schema",
		},
		ContentURI: TestingDocument.ContentURI,
		Checksum:   TestingDocument.Checksum,
	})
	supportedSchema := types.MetadataSchema{
		Type:      "metadata-schema",
		SchemaURI: "https://example.com/schema",
		Version:   "1.0.0",
	}
	k.AddSupportedMetadataScheme(ctx, supportedSchema)

	res := handler(ctx, msgShareDocument)
	assert.True(t, res.IsOK())
}

func Test_handleMsgShareDocument_MetadataSchemeType_NotSupported(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	msgShareDocument := types.MsgShareDocument(types.Document{
		Sender:     TestingSender,
		Recipients: ctypes.Addresses{TestingRecipient},
		UUID:       TestingDocument.UUID,
		Metadata: types.DocumentMetadata{
			ContentURI: TestingDocument.Metadata.ContentURI,
			SchemaType: "non-existent-schema-type",
		},
		ContentURI: TestingDocument.ContentURI,
		Checksum:   TestingDocument.Checksum,
	})

	res := handler(ctx, msgShareDocument)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

// -----------------------------------
// --- handleMsgSendDocumentReceipt
// -----------------------------------

func Test_handleMsgSendDocumentReceipt(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	msgSendDocumentReceipt := types.MsgSendDocumentReceipt(TestingDocumentReceipt)

	res := handler(ctx, msgSendDocumentReceipt)
	assert.True(t, res.IsOK())
}

// -------------------------------------------
// --- handleMsgAddSupportedMetadataSchema
// -------------------------------------------

func Test_handleMsgAddSupportedMetadataSchema_NotTrustedSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgAddSupportedMetadataSchema := types.MsgAddSupportedMetadataSchema{
		Signer: TestingSender,
		Schema: types.MetadataSchema{
			Type:      "schema-type",
			SchemaURI: "https://example.com/schema",
			Version:   "1.0.0",
		},
	}
	res := handler(ctx, msgAddSupportedMetadataSchema)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddSupportedMetadataSchema_TrustedSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	k.AddTrustedSchemaProposer(ctx, TestingSender)

	msgAddSupportedMetadataSchema := types.MsgAddSupportedMetadataSchema{
		Signer: TestingSender,
		Schema: types.MetadataSchema{
			Type:      "schema-type",
			SchemaURI: "https://example.com/schema",
			Version:   "1.0.0",
		},
	}
	res := handler(ctx, msgAddSupportedMetadataSchema)
	assert.True(t, res.IsOK())
}

// ------------------------------------------------
// --- handleMsgAddTrustedMetadataSchemaProposer
// ------------------------------------------------

func Test_handleMsgAddTrustedMetadataSchemaProposer_MissingGovernment(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: TestingSender,
		Signer:   TestingRecipient,
	}
	res := handler(ctx, msgAddTrustedMetadataSchemaProposer)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddTrustedMetadataSchemaProposer_IncorrectSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	_ = k.GovernmentKeeper.SetGovernmentAddress(ctx, TestingRecipient)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: TestingSender,
		Signer:   TestingSender,
	}
	res := handler(ctx, msgAddTrustedMetadataSchemaProposer)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddTrustedMetadataSchemaProposer_CorrectSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	_ = k.GovernmentKeeper.SetGovernmentAddress(ctx, TestingRecipient)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: TestingSender,
		Signer:   TestingRecipient,
	}
	res := handler(ctx, msgAddTrustedMetadataSchemaProposer)
	assert.True(t, res.IsOK())
}

// -------------------
// --- Default case
// -------------------

func Test_invalidMsg(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	res := handler(ctx, sdk.NewTestMsg())
	assert.False(t, res.IsOK())
	assert.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
