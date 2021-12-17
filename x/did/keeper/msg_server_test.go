package keeper

import (
	"testing"
	"time"

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

	msg := &types.MsgSetDidDocument{
		Context: []string{
			"https://www.w3.org/ns/did/v1",
		},
		ID: addr.String(),
		VerificationMethod: []*types.VerificationMethod{
			{
				ID:                 "",
				Type:               "",
				Controller:         "",
				PublicKeyMultibase: "",
			},
		},
		Service: []*types.Service{
			{
				ID:              "cfbff1f9-8b30-4223-9648-5d4f7fc0a159",
				Type:            "agent",
				ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
			},
		},
		Authentication:       []*types.VerificationMethod{},
		AssertionMethod:      []*types.VerificationMethod{},
		CapabilityDelegation: []*types.VerificationMethod{},
		CapabilityInvocation: []*types.VerificationMethod{},
		KeyAgreement:         []*types.VerificationMethod{},
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
