package docs

import (
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgShareDocument:
			return handleMsgShareDocument(ctx, keeper, msg)
		case MsgSendDocumentReceipt:
			return handleMsgSendDocumentReceipt(ctx, keeper, msg)
		case MsgAddSupportedMetadataSchema:
			return handleMsgAddSupportedMetadataSchema(ctx, keeper, msg)
		case MsgAddTrustedMetadataSchemaProposer:
			return handleMsgAddTrustedMetadataSchemaProposer(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgShareDocument(ctx sdk.Context, keeper Keeper, msg MsgShareDocument) sdk.Result {

	// The metadata schema type is being specified
	if len(strings.TrimSpace(msg.Metadata.SchemaType)) != 0 {

		// Check its validity
		if !keeper.IsMetadataSchemeTypeSupported(ctx, msg.Metadata.SchemaType) {
			errMsg := fmt.Sprintf("Unsupported metadata schema: %s", msg.Metadata.SchemaType)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		// Delete the custom data
		msg.Metadata.Schema = nil
	}

	// Share the document
	if err := keeper.SaveDocument(ctx, Document(msg)); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleMsgSendDocumentReceipt(ctx sdk.Context, keeper Keeper, msg MsgSendDocumentReceipt) sdk.Result {
	if err := keeper.SaveReceipt(ctx, types.DocumentReceipt(msg)); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleMsgAddSupportedMetadataSchema(ctx sdk.Context, keeper Keeper, msg MsgAddSupportedMetadataSchema) sdk.Result {

	// Make sure the signer is valid
	if !keeper.IsTrustedSchemaProposer(ctx, msg.Signer) {
		errMsg := fmt.Sprintf("Signer is not a trusted one: %s", msg.Signer.String())
		return sdk.ErrInvalidAddress(errMsg).Result()
	}

	// Add the schema
	keeper.AddSupportedMetadataScheme(ctx, msg.Schema)

	return sdk.Result{}
}

func handleMsgAddTrustedMetadataSchemaProposer(ctx sdk.Context, keeper Keeper, msg MsgAddTrustedMetadataSchemaProposer) sdk.Result {
	// Authenticate the signer
	governmentAddr := keeper.GovernmentKeeper.GetGovernmentAddress(ctx)
	if !msg.Signer.Equals(governmentAddr) {
		errMsg := fmt.Sprintf("Only the government can add a trusted metadata schema proposer")
		return sdk.ErrInvalidAddress(errMsg).Result()
	}

	// Add the trusted schema proposer
	keeper.AddTrustedSchemaProposer(ctx, msg.Proposer)
	return sdk.Result{}
}
