package keeper

import (
	"context"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
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
			desc: "share a document",
			msg: &types.MsgShareDocument{
				UUID:       uuid.NewV4().String(),
				Sender:     testingSender.String(),
				Recipients: []string{testingRecipient.String()},
			},
			err: false,
		},
		{
			desc: "invalid document uuid",
			msg: &types.MsgShareDocument{
				UUID:       "",
				Sender:     testingSender.String(),
				Recipients: []string{testingRecipient.String()},
			},
			err: true,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := srv.ShareDocument(ctx, tc.msg)
			if tc.err {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.UUID, response.UUID)
			}
		})
	}
}

func TestSendDocument(t *testing.T) {
	srv, ctx, k := setupMsgServer(t)
	docs := createNDocument(&k, sdk.UnwrapSDKContext(ctx), 1)

	for _, tc := range []struct {
		desc string
		msg  *types.MsgSendDocumentReceipt
		err  bool
	}{
		{
			desc: "send a document receipt",
			msg: &types.MsgSendDocumentReceipt{
				UUID:         uuid.NewV4().String(),
				Sender:       docs[0].Recipients[0],
				Recipient:    docs[0].Sender,
				DocumentUUID: docs[0].UUID,
			},
			err: false,
		},
		{
			desc: "invalid document receipt uuid",
			msg: &types.MsgSendDocumentReceipt{
				UUID:         "",
				Sender:       docs[0].Recipients[0],
				Recipient:    docs[0].Sender,
				DocumentUUID: docs[0].UUID,
			},
			err: true,
		},
		{
			desc: "invalid document uuid",
			msg: &types.MsgSendDocumentReceipt{
				UUID:         uuid.NewV4().String(),
				Sender:       docs[0].Recipients[0],
				Recipient:    docs[0].Sender,
				DocumentUUID: "",
			},
			err: true,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := srv.SendDocumentReceipt(ctx, tc.msg)
			if tc.err {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.UUID, response.UUID)
			}
		})
	}
}
