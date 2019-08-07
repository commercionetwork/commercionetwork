package commercioid

import (
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var msgSetId = MsgSetIdentity{
	Did:          keeper.TestOwnerIdentity,
	DDOReference: keeper.TestIdentityRef,
	Owner:        keeper.TestOwner,
}

var msgCreateConn = MsgCreateConnection{
	FirstUser:  keeper.TestOwnerIdentity,
	SecondUser: keeper.TestRecipient,
	Signer:     keeper.TestOwner,
}

var testUtils = keeper.TestUtils

var handler = NewHandler(testUtils.IdKeeper)

func TestValidMsg_StoreDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgSetId)

	require.True(t, res.IsOK())
}

func TestValidMsg_ShareDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgCreateConn)

	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized commercioid message type"))
}
