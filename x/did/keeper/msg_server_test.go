package keeper

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	configTestPrefixes()
}

func configTestPrefixes() {
	AccountAddressPrefix := "did:com:"
	AccountPubKeyPrefix := AccountAddressPrefix + "pub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.Seal()
}

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
	validBase64RsaVerificationKey2018 := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB"
	validBase64RsaSignature2018 := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHgkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScADG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUczhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7ZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0O2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfKHQIDAQAB"

	verificationMethodVer := types.VerificationMethod{
		ID:                 msg.DidDocument.ID + "#keys-1",
		Type:               types.RsaVerificationKey2018,
		Controller:         msg.DidDocument.ID,
		PublicKeyMultibase: string(types.MultibaseCodeBase64NoPadding) + validBase64RsaVerificationKey2018,
	}
	verificationMethodSign := types.VerificationMethod{
		ID:                 msg.DidDocument.ID + "#keys-2",
		Type:               types.RsaSignature2018,
		Controller:         msg.DidDocument.ID,
		PublicKeyMultibase: string(types.MultibaseCodeBase64NoPadding) + validBase64RsaSignature2018,
	}

	msg.DidDocument.VerificationMethod = []*types.VerificationMethod{
		&verificationMethodVer,
		&verificationMethodSign,
	}

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

	newMsg.DidDocument.AssertionMethod = []string{"#keys-1"}

	resp, err = srv.UpdateIdentity(sdkCtx, &newMsg)
	require.NoError(t, err)
	assert.Equal(t, &types.MsgSetIdentityResponse{}, resp)

	identityUpdated, err := k.GetLastIdentityOfAddress(ctx, did)
	assert.NoError(t, err)
	require.Equal(t, newMsg.DidDocument, identityUpdated.DidDocument)

	require.Equal(t, firstIdentity.Metadata.Created, identityUpdated.Metadata.Created)
	require.NotEqual(t, firstIdentity.Metadata.Updated, identityUpdated.Metadata.Updated)

}
