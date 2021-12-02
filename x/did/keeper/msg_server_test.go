package keeper

import (
	"context"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	keeper, ctx := setupKeeper(t)
	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx)
}

// func TestSetIdentityMsgServerCreate(t *testing.T) {
// 	srv, ctx := setupMsgServer(t)
// 	creator := "A"
// 	for i := 0; i < 5; i++ {
// 		_, err := srv.SetDid(ctx, &types.MsgSetDid{ID: creator})
// 		require.NoError(t, err)
// 		//assert.Equal(t, i, int(resp.UUID))
// 	}
// }
