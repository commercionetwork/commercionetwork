package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	/*"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"*/
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

func TestDocumentQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDocument(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDocumentRequest
		response *types.QueryGetDocumentResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetDocumentRequest{UUID: msgs[0].UUID},
			response: &types.QueryGetDocumentResponse{Document: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetDocumentRequest{UUID: msgs[1].UUID},
			response: &types.QueryGetDocumentResponse{Document: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetDocumentRequest{UUID: fmt.Sprint(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Document(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

/*
func TestDocumentQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDocument(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryGetSentDocumentsRequest {
		return &types.QueryGetSentDocumentsRequest{
			Address: string(TestingSender),
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SentDocuments(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			/*for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.Document[j-i])
			}*
			for _, respDocument := range resp.Document {
				assert.Contains(t, msgs, *respDocument)
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SentDocuments(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			/*for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.Document[j-i])
			}*
			for _, respDocument := range resp.Document {
				assert.Contains(t, msgs, *respDocument)
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.SentDocuments(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.SentDocuments(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
*/
