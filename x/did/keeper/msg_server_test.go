package keeper

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/did/types"
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

	// create
	dateString := types.ValidIdentity.Metadata.Created
	createdTimestamp, err := time.Parse(types.ComplaintW3CTime, dateString)
	require.NoError(t, err)
	ctx = ctx.WithBlockTime(createdTimestamp.UTC())

	sdkCtx := sdk.WrapSDKContext(ctx)

	msg := types.MsgSetIdentity{
		DidDocument: types.ValidIdentity.DidDocument,
	}

	did := msg.DidDocument.ID

	_, err = k.GetLastIdentityOfAddress(ctx, did)
	assert.Error(t, err)

	resp, err := srv.UpdateIdentity(sdkCtx, &msg)
	require.NoError(t, err)
	assert.Equal(t, &types.MsgSetIdentityResponse{}, resp)

	// try to update the identity with the same DDO as the previous one
	_, err = srv.UpdateIdentity(sdkCtx, &msg)
	require.Error(t, err)

	firstIdentity, err := k.GetLastIdentityOfAddress(ctx, did)
	assert.NoError(t, err)
	require.Equal(t, msg.DidDocument, firstIdentity.DidDocument)
	expectedFirstMetadata := types.Metadata{
		Created: dateString,
		Updated: dateString,
	}
	require.Equal(t, &expectedFirstMetadata, firstIdentity.Metadata)

	// update
	ctx = sdk.UnwrapSDKContext(sdkCtx)
	updatedTimestamp := createdTimestamp.Add(time.Hour)
	ctx = ctx.WithBlockTime(updatedTimestamp)

	sdkCtx = sdk.WrapSDKContext(ctx)

	newMsg := msg
	newMsg.DidDocument.AssertionMethod = []string{"#key-1"}

	resp, err = srv.UpdateIdentity(sdkCtx, &newMsg)
	require.NoError(t, err)
	assert.Equal(t, &types.MsgSetIdentityResponse{}, resp)

	identityUpdated, err := k.GetLastIdentityOfAddress(ctx, did)
	assert.NoError(t, err)
	require.Equal(t, newMsg.DidDocument, identityUpdated.DidDocument)

	require.Equal(t, firstIdentity.Metadata.Created, identityUpdated.Metadata.Created)
	require.NotEqual(t, firstIdentity.Metadata.Updated, identityUpdated.Metadata.Updated)

}
