package keeper

import (
	"fmt"
	"testing"
	"time"

	commons "github.com/commercionetwork/commercionetwork/x/common/types"
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

	// creation
	ctx = ctx.WithBlockTime(time.Now())

	sdkCtx := sdk.WrapSDKContext(ctx)

	var context commons.Strings = []string{
		types.ContextDidV1,
		"https://w3id.org/security/suites/ed25519-2018/v1",
		"https://w3id.org/security/suites/x25519-2019/v1",
	}

	did := "did:com:13jckgxmj3v8jpqdeq8zxwcyhv7gc3dzmrqqger"
	// check only did:com

	msg := &types.MsgSetDidDocument{
		Context: context,
		ID:      addr.String(),
		VerificationMethod: []*types.VerificationMethod{
			{
				ID:                 fmt.Sprint(did, "#key-1"),
				Type:               "Ed25519VerificationKey2020",
				Controller:         did,
				PublicKeyMultibase: "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
			},
			{
				ID:                 fmt.Sprint(did, "#key-agreement-1"),
				Type:               "X25519KeyAgreementKey2019",
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
				ID:              "A",
				Type:            "agent",
				ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
			},
			{
				ID:              "B",
				Type:            "xdi",
				ServiceEndpoint: "https://commerc.io/xdi/serviceEndpoint/",
			},
		},
	}

	resp, err := srv.SetDidDocument(sdkCtx, msg)
	require.NoError(t, err)
	assert.Equal(t, addr.String(), resp.ID)

	// update
	ctx = sdk.UnwrapSDKContext(sdkCtx)
	ctx = ctx.WithBlockTime(time.Now().Add(time.Hour))

	sdkCtx = sdk.WrapSDKContext(ctx)

	resp, err = srv.SetDidDocument(sdkCtx, msg)
	require.NoError(t, err)
	assert.Equal(t, addr.String(), resp.ID)

	for _, d := range k.GetAllDidDocuments(ctx) {
		t.Log(d.Created, "VS", d.Updated)
		assert.NotEqual(t, d.Created, d.Updated)
	}

	assert.True(t, k.HasDidDocument(ctx, addr.String()))

}
