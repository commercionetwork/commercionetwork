package commercioid

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var handler = NewHandler(input.idKeeper)

func Test_handleMsgCreateIdentity(t *testing.T) {

	actual := handler(input.ctx, msgSetId)

	assert.Equal(t, sdk.Result{}, actual)
}

func Test_handleMsgCreateIdentity_incorrectSigner(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(identityRef))

	var msgSetId = MsgSetIdentity{
		DID:          ownerIdentity,
		DDOReference: identityRef,
		Owner:        sdk.AccAddress{},
	}

	expected := sdk.ErrUnauthorized("Incorrect Signer").Result()

	actual := handler(input.ctx, msgSetId)

	assert.Equal(t, expected, actual)
}

func Test_handleMsgCreateConnection(t *testing.T) {

	actual := handler(input.ctx, msgCreateConn)

	assert.Equal(t, sdk.Result{}, actual)
}

func Test_handleMsgCreateConnection_SignerIsntTheOwnerOfIdentity(t *testing.T) {
	store := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(identityRef))

	var msgCreateConn = MsgCreateConnection{
		FirstUser:  ownerIdentity,
		SecondUser: recipient,
		Signer:     sdk.AccAddress{},
	}

	expected := sdk.ErrUnauthorized("The signer must own either the first or the second DID").Result()

	actual := handler(input.ctx, msgCreateConn)

	assert.Equal(t, expected, actual)
}
