package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/documents/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ShareDocument(goCtx context.Context, msg *types.MsgShareDocument) (*types.MsgShareDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.SaveDocument(ctx, types.Document(*msg)); err != nil {
		return nil, err
	}

	// TODO: add event

	return &types.MsgShareDocumentResponse{UUID: msg.UUID}, nil
}

func (k msgServer) SendDocumentReceipt(goCtx context.Context, msg *types.MsgSendDocumentReceipt) (*types.MsgSendDocumentReceiptResponse, error) {
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
