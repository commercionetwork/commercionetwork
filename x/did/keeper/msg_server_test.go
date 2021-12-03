package keeper

import (
	"context"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	keeper, ctx := setupKeeper(t)

	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx)
}

func TestSetIdentityMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	_, _, addr := testdata.KeyTestPubAddr()
	resp, err := srv.SetDidDocument(ctx, &types.MsgSetDidDocument{DidDocument: &types.DidDocument{ID: addr.String()}})
	require.NoError(t, err)
	assert.Equal(t, addr.String(), resp.ID)
	t.FailNow()
}
