package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, Keeper, sdk.Context) {
	keeper, ctx := setupKeeper(t)

	return NewMsgServerImpl(*keeper), *keeper, ctx
}

func Test_SetDidDocument(t *testing.T) {
	srv, k, ctx := setupMsgServer(t)
	_, _, addr := testdata.KeyTestPubAddr()

	sdkCtx := sdk.WrapSDKContext(ctx)

	resp, err := srv.SetDidDocument(sdkCtx, &types.MsgSetDidDocument{ID: addr.String()})
	require.NoError(t, err)
	assert.Equal(t, addr.String(), resp.ID)

	resp, err = srv.SetDidDocument(sdkCtx, &types.MsgSetDidDocument{ID: addr.String()})
	require.NoError(t, err)
	assert.Equal(t, addr.String(), resp.ID)

	for _, d := range k.GetAllDidDocuments(ctx) {
		assert.True(t, d.Created != d.Updated)
	}

	assert.True(t, k.HasDidDocument(ctx, addr.String()))

}
