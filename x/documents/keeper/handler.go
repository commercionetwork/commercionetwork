package keeper

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgShareDocument:
			return handleMsgShareDocument(ctx, keeper, msg)
		case types.MsgSendDocumentReceipt:
			return handleMsgSendDocumentReceipt(ctx, keeper, msg)
		case types.MsgAddSupportedMetadataSchema:
			return handleMsgAddSupportedMetadataSchema(ctx, keeper, msg)
		case types.MsgAddTrustedMetadataSchemaProposer:
			return handleMsgAddTrustedMetadataSchemaProposer(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgShareDocument(ctx sdk.Context, keeper Keeper, msg types.MsgShareDocument) (*sdk.Result, error) {

	// The metadata schema type is being specified
	if len(strings.TrimSpace(msg.Metadata.SchemaType)) != 0 {

		// Check its validity
		if !keeper.IsMetadataSchemeTypeSupported(ctx, msg.Metadata.SchemaType) {
			errMsg := fmt.Sprintf("Unsupported metadata schema: %s", msg.Metadata.SchemaType)
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}

		// Delete the custom data
		msg.Metadata.Schema = nil
	}

	// Share the document
	if err := keeper.SaveDocument(ctx, types.Document(msg)); err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Document successfully shared"}, nil
}

func handleMsgSendDocumentReceipt(ctx sdk.Context, keeper Keeper, msg types.MsgSendDocumentReceipt) (*sdk.Result, error) {
	if err := keeper.SaveReceipt(ctx, types.DocumentReceipt(msg)); err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Receipt Document successfully sent"}, nil
}

func handleMsgAddSupportedMetadataSchema(ctx sdk.Context, keeper Keeper, msg types.MsgAddSupportedMetadataSchema) (*sdk.Result, error) {

	// Make sure the signer is valid
	if !keeper.IsTrustedSchemaProposer(ctx, msg.Signer) {
		errMsg := fmt.Sprintf("Signer is not a trusted one: %s", msg.Signer.String())
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
	}

	// Add the schema
	keeper.AddSupportedMetadataScheme(ctx, msg.Schema)

	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Metadata schema successfully added"}, nil
}

func handleMsgAddTrustedMetadataSchemaProposer(ctx sdk.Context, keeper Keeper, msg types.MsgAddTrustedMetadataSchemaProposer) (*sdk.Result, error) {
	// Authenticate the signer
	governmentAddr := keeper.GovernmentKeeper.GetGovernmentAddress(ctx)
	if !msg.Signer.Equals(governmentAddr) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Only the government can add a trusted metadata schema proposer")
	}

	// Add the trusted schema proposer
	keeper.AddTrustedSchemaProposer(ctx, msg.Proposer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Trusted Metadata Schema Proposer successfully added"}, nil
}
