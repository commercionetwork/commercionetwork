package keeper

import (
	"fmt"
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

// TODO use valid DID document content from the 'did.keeper.types' testing package
func Test_SetDidDocument(t *testing.T) {
	srv, k, ctx := setupMsgServer(t)

	// creation
	ctx = ctx.WithBlockTime(time.Now())

	sdkCtx := sdk.WrapSDKContext(ctx)

	did := "did:com:13jckgxmj3v8jpqdeq8zxwcyhv7gc3dzmrqqger"

	msg := types.MsgSetDidDocument{
		Context: []string{
			types.ContextDidV1,
			"https://w3id.org/security/suites/ed25519-2018/v1",
			"https://w3id.org/security/suites/x25519-2019/v1",
		},
		ID: did,
		VerificationMethod: []*types.VerificationMethod{
			{
				ID:                 fmt.Sprint(did, "#key-1"),
				Type:               "RsaSignature2018",
				Controller:         did,
				PublicKeyMultibase: "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
			},
			{
				ID:                 fmt.Sprint(did, "#key-agreement-1"),
				Type:               "RsaVerificationKey2018",
				Controller:         did,
				PublicKeyMultibase: "FK2c4QudVyaodvX9LARDsbihkVBvWxe8oiJAiYQ2JpdC",
			},
		},
		Authentication: []string{
			fmt.Sprint(did, "#key-1"),
		},
		AssertionMethod: []string{
			fmt.Sprint(did, "#key-1"),
		},
		KeyAgreement: []string{
			fmt.Sprint(did, "#key-agreement-1"),
		},
		CapabilityInvocation: nil,
		CapabilityDelegation: nil,
		Service: []*types.Service{
			{
				ID:              "https://bar.example.com",
				Type:            "agent",
				ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
			},
			{
				ID:              "https://foo.example.com",
				Type:            "xdi",
				ServiceEndpoint: "https://commerc.io/xdi/serviceEndpoint/",
			},
		},
	}

	_, err := k.GetDidDocumentOfAddress(ctx, did)
	assert.Error(t, err)

	resp, err := srv.SetDidDocument(sdkCtx, &msg)
	require.NoError(t, err)
	assert.Equal(t, did, resp.ID)

	var ddo types.DidDocument

	ddo, err = k.GetDidDocumentOfAddress(ctx, did)
	assert.NoError(t, err)
	requireEqualMsgSetDidDocumentWithDidDocument(t, msg, ddo)

	// update
	ctx = sdk.UnwrapSDKContext(sdkCtx)
	ctx = ctx.WithBlockTime(time.Now().Add(time.Hour))

	sdkCtx = sdk.WrapSDKContext(ctx)

	newMsg := msg
	newMsg.AssertionMethod = []string{"#key-1"}

	resp, err = srv.SetDidDocument(sdkCtx, &newMsg)
	require.NoError(t, err)
	assert.Equal(t, did, resp.ID)

	assert.True(t, k.HasDidDocument(ctx, did))

	ddoUpdated, err := k.GetDidDocumentOfAddress(ctx, did)
	assert.NoError(t, err)
	requireEqualMsgSetDidDocumentWithDidDocument(t, newMsg, ddoUpdated)

	require.Equal(t, ddo.Created, ddoUpdated.Created)
	require.NotEqual(t, ddo.Updated, ddoUpdated.Updated)

}

func requireEqualMsgSetDidDocumentWithDidDocument(t *testing.T, msg types.MsgSetDidDocument, ddo types.DidDocument) {
	require.Equal(t, msg.ID, ddo.ID)
	require.Equal(t, msg.Context, ddo.Context)
	require.Equal(t, msg.AssertionMethod, ddo.AssertionMethod)
	require.Equal(t, msg.Authentication, ddo.Authentication)
	require.Equal(t, msg.CapabilityDelegation, ddo.CapabilityDelegation)
	require.Equal(t, msg.CapabilityInvocation, ddo.CapabilityInvocation)
	require.Equal(t, msg.KeyAgreement, ddo.KeyAgreement)
	require.Equal(t, msg.Service, ddo.Service)
	require.Equal(t, msg.VerificationMethod, ddo.VerificationMethod)
}
