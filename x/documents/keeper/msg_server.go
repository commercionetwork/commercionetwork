package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/documents/types"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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
	timestamp := ctx.BlockTime()
	doc := types.Document{
		ContentURI: msg.ContentURI,
		UUID: msg.UUID,
		Metadata: msg.Metadata,
		Checksum: msg.Checksum,
		Sender:  msg.Sender,
		Recipients: msg.Recipients,
		EncryptionData: msg.EncryptionData,
		DoSign: msg.DoSign,
		Timestamp: &timestamp,
	}

	if err := k.SaveDocument(ctx, doc); err != nil {
		return nil, err
	}

	ctypes.EmitCommonEvents(ctx, msg.Sender)
	return &types.MsgShareDocumentResponse{UUID: msg.UUID}, nil
}

func (k msgServer) SendDocumentReceipt(goCtx context.Context, msg *types.MsgSendDocumentReceipt) (*types.MsgSendDocumentReceiptResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	timestamp := ctx.BlockTime()
	docReceipts := types.DocumentReceipt{
		UUID: msg.UUID,
		Sender: msg.Sender,
		Recipient: msg.Recipient,
		TxHash: msg.TxHash,
		DocumentUUID: msg.DocumentUUID,
		Proof: msg.Proof,
		Timestamp: &timestamp,
	}
	if err := k.SaveReceipt(ctx, docReceipts); err != nil {
		return nil, err
	}

	ctypes.EmitCommonEvents(ctx, msg.Sender)
	return &types.MsgSendDocumentReceiptResponse{UUID: msg.UUID}, nil
}
