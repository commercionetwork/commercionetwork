package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

func TestShareDocumentMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	creator := "A"
	for i := 0; i < 5; i++ {
		_, err := srv.CreateDocument(ctx, &types.MsgShareDocument{Sender: creator})
		require.NoError(t, err)
		//assert.Equal(t, i, int(resp.UUID))
	}
}
