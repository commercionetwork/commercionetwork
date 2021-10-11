package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/id/types"
)

func TestSetIdentityMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	creator := "A"
	for i := 0; i < 5; i++ {
		_, err := srv.SetIdentity(ctx, &types.MsgSetIdentity{ID: creator})
		require.NoError(t, err)
		//assert.Equal(t, i, int(resp.UUID))
	}
}
