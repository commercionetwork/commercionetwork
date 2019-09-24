package docs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------
// --- handleMsgShareDocument
// -----------------------------

func Test_handleMsgShareDocument_CustomMetadataSpecs(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgShareDocument := NewMsgShareDocument(TestingSender, types.Addresses{TestingRecipient}, TestingDocument)

	res := handler(ctx, msgShareDocument)
	require.True(t, res.IsOK())
}

func Test_handleMsgShareDocument_MetadataSchemeType_Supported(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgShareDocument := MsgShareDocument{
		Sender:     TestingSender,
		Recipients: types.Addresses{TestingRecipient},
		Document: Document{
			Uuid: TestingDocument.Uuid,
			Metadata: DocumentMetadata{
				ContentUri: TestingDocument.Metadata.ContentUri,
				SchemaType: "metadata-schema",
			},
			ContentUri: TestingDocument.ContentUri,
			Checksum:   TestingDocument.Checksum,
		},
	}
	supportedSchema := MetadataSchema{
		Type:      "metadata-schema",
		SchemaUri: "https://example.com/schema",
		Version:   "1.0.0",
	}
	k.AddSupportedMetadataScheme(ctx, supportedSchema)

	res := handler(ctx, msgShareDocument)
	assert.True(t, res.IsOK())
}

func Test_handleMsgShareDocument_MetadataSchemeType_NotSupported(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	msgShareDocument := MsgShareDocument{
		Sender:     TestingSender,
		Recipients: types.Addresses{TestingRecipient},
		Document: Document{
			Uuid: TestingDocument.Uuid,
			Metadata: DocumentMetadata{
				ContentUri: TestingDocument.Metadata.ContentUri,
				SchemaType: "non-existent-schema-type",
			},
			ContentUri: TestingDocument.ContentUri,
			Checksum:   TestingDocument.Checksum,
		},
	}

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
	msgSendDocumentReceipt := MsgSendDocumentReceipt(TestingDocumentReceipt)

	res := handler(ctx, msgSendDocumentReceipt)
	assert.True(t, res.IsOK())
}

// -------------------------------------------
// --- handleMsgAddSupportedMetadataSchema
// -------------------------------------------

func Test_handleMsgAddSupportedMetadataSchema_NotTrustedSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgAddSupportedMetadataSchema := MsgAddSupportedMetadataSchema{
		Signer: TestingSender,
		Schema: MetadataSchema{
			Type:      "schema-type",
			SchemaUri: "https://example.com/schema",
			Version:   "1.0.0",
		},
	}
	res := handler(ctx, msgAddSupportedMetadataSchema)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleMsgAddSupportedMetadataSchema_TrustedSigner(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	signers := []sdk.AccAddress{TestingSender}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(MetadataSchemaProposersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	msgAddSupportedMetadataSchema := MsgAddSupportedMetadataSchema{
		Signer: TestingSender,
		Schema: MetadataSchema{
			Type:      "schema-type",
			SchemaUri: "https://example.com/schema",
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

	msgAddTrustedMetadataSchemaProposer := MsgAddTrustedMetadataSchemaProposer{
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

	msgAddTrustedMetadataSchemaProposer := MsgAddTrustedMetadataSchemaProposer{
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

	msgAddTrustedMetadataSchemaProposer := MsgAddTrustedMetadataSchemaProposer{
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
	assert.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
