package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateDocument(goCtx context.Context, msg *types.MsgShareDocument) (*types.MsgShareDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var document = types.Document{
		Sender:         msg.Sender,
		Recipients:     msg.Recipients,
		UUID:           msg.UUID,
		Metadata:       msg.Metadata,
		ContentURI:     msg.ContentURI,
		Checksum:       msg.Checksum,
		EncryptionData: msg.EncryptionData,
		DoSign:         msg.DoSign,
	}
	/*
	id := k.AppendDocument(
		ctx,
		document,
	)*/
	err := k.SaveDocument(ctx, document)
	if err != nil {
		return nil, err
	}

	return &types.MsgShareDocumentResponse{UUID: document.UUID}, nil
}

func (k msgServer) SendDocument(goCtx context.Context, msg *types.MsgSendDocumentReceipt) (*types.MsgSendDocumentReceiptResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.SaveReceipt(ctx, types.DocumentReceipt(*msg)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))
	return &types.MsgSendDocumentReceiptResponse{UUID: msg.UUID}, nil
}
