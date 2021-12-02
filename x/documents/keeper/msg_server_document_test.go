package keeper

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

func TestShareDocumentMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	for i := 0; i < 5; i++ {
		_, err := srv.CreateDocument(ctx, &types.MsgShareDocument{
			UUID: uuid.NewV4().String(),
		})
		require.NoError(t, err)
		//assert.Equal(t, i, int(resp.UUID))
	}
}
