package keeper

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/stretchr/testify/require"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/documents/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// -----------------------------
// --- handleMsgShareDocument
// -----------------------------

func Test_handleMsgShareDocument_CustomMetadataSpecs(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)

	msgShareDocument := types.NewMsgShareDocument(TestingDocument)

	_, err := handler(ctx, msgShareDocument)
	require.NoError(t, err)
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

	_, err := handler(ctx, msgShareDocument)
	require.NoError(t, err)
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

	_, err := handler(ctx, msgShareDocument)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
}

// -----------------------------------
// --- handleMsgSendDocumentReceipt
// -----------------------------------

func Test_handleMsgSendDocumentReceipt(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	require.NoError(t, k.SaveDocument(ctx, TestingDocument))

	tdr := TestingDocumentReceipt
	tdr.DocumentUUID = TestingDocument.UUID

	msgSendDocumentReceipt := types.MsgSendDocumentReceipt(tdr)

	_, err := handler(ctx, msgSendDocumentReceipt)
	require.NoError(t, err)
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
	_, err := handler(ctx, msgAddSupportedMetadataSchema)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
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
	_, err := handler(ctx, msgAddSupportedMetadataSchema)
	require.NoError(t, err)
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
	_, err := handler(ctx, msgAddTrustedMetadataSchemaProposer)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
}

func Test_handleMsgAddTrustedMetadataSchemaProposer_IncorrectSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	_ = k.GovernmentKeeper.SetGovernmentAddress(ctx, TestingRecipient)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: TestingSender,
		Signer:   TestingSender,
	}
	_, err := handler(ctx, msgAddTrustedMetadataSchemaProposer)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
}

func Test_handleMsgAddTrustedMetadataSchemaProposer_CorrectSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	_ = k.GovernmentKeeper.SetGovernmentAddress(ctx, TestingRecipient)

	msgAddTrustedMetadataSchemaProposer := types.MsgAddTrustedMetadataSchemaProposer{
		Proposer: TestingSender,
		Signer:   TestingRecipient,
	}
	_, err := handler(ctx, msgAddTrustedMetadataSchemaProposer)
	require.NoError(t, err)
}

// -------------------
// --- Default case
// -------------------

func Test_invalidMsg(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var handler = NewHandler(k)
	_, err := handler(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
