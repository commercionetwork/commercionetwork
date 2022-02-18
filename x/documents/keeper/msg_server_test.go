package keeper

import (
	"context"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context, Keeper) {
	keeper, ctx := setupKeeper(t)
	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx), *keeper
}

func TestShareDocument(t *testing.T) {
	srv, ctx, _ := setupMsgServer(t)
	for _, tc := range []struct {
		desc string
		msg  *types.MsgShareDocument
		err  bool
	}{
		{
			desc: "ok",
			msg: &types.MsgShareDocument{
				Sender:         types.ValidDocument.Sender,
				Recipients:     types.ValidDocument.Recipients,
				UUID:           types.AnotherValidDocument.UUID,
				Metadata:       types.ValidDocument.Metadata,
				ContentURI:     types.ValidDocument.ContentURI,
				Checksum:       types.ValidDocument.Checksum,
				EncryptionData: types.ValidDocument.EncryptionData,
				DoSign:         types.ValidDocument.DoSign,
			},
		},
		{
			desc: "invalid document",
			msg: &types.MsgShareDocument{
				Sender:         types.ValidDocument.Sender,
				Recipients:     types.ValidDocument.Recipients,
				UUID:           "",
				Metadata:       types.ValidDocument.Metadata,
				ContentURI:     types.ValidDocument.ContentURI,
				Checksum:       types.ValidDocument.Checksum,
				EncryptionData: types.ValidDocument.EncryptionData,
				DoSign:         types.ValidDocument.DoSign,
			},
			err: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := srv.ShareDocument(ctx, tc.msg)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.UUID, response.UUID)
			}
		})
	}
}

func TestSendDocumentReceipt(t *testing.T) {

	for _, tc := range []struct {
		desc            string
		storedDocuments []types.Document
		msg             *types.MsgSendDocumentReceipt
		err             bool
	}{
		{
			desc:            "ok",
			storedDocuments: []types.Document{types.ValidDocument},
			msg: &types.MsgSendDocumentReceipt{
				UUID:         types.ValidDocumentReceiptRecipient1.UUID,
				Sender:       types.ValidDocumentReceiptRecipient1.Sender,
				Recipient:    types.ValidDocumentReceiptRecipient1.Recipient,
				TxHash:       types.ValidDocumentReceiptRecipient1.TxHash,
				DocumentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
				Proof:        types.ValidDocumentReceiptRecipient1.Proof,
			},
		},
		{
			desc: "no corresponding document",
			msg: &types.MsgSendDocumentReceipt{
				UUID:         types.ValidDocumentReceiptRecipient1.UUID,
				Sender:       types.ValidDocumentReceiptRecipient1.Sender,
				Recipient:    types.ValidDocumentReceiptRecipient1.Recipient,
				TxHash:       types.ValidDocumentReceiptRecipient1.TxHash,
				DocumentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
				Proof:        types.ValidDocumentReceiptRecipient1.Proof,
			},
			err: true,
		},
		{
			desc:            "invalid document receipt",
			storedDocuments: []types.Document{types.ValidDocument},
			msg: &types.MsgSendDocumentReceipt{
				UUID:         "",
				Sender:       types.ValidDocumentReceiptRecipient1.Sender,
				Recipient:    types.ValidDocumentReceiptRecipient1.Recipient,
				TxHash:       types.ValidDocumentReceiptRecipient1.TxHash,
				DocumentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
				Proof:        types.ValidDocumentReceiptRecipient1.Proof,
			},
			err: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx, k := setupMsgServer(t)

			for _, document := range tc.storedDocuments {
				uctx := sdk.UnwrapSDKContext(ctx)
				err := k.SaveDocument(uctx, document)
				require.NoError(t, err)
			}

			response, err := srv.SendDocumentReceipt(ctx, tc.msg)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.UUID, response.UUID)
			}
		})
	}
}
