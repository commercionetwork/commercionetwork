package commercioid

import (
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var msgSetId = MsgSetIdentity{
	DidDocumentReference: keeper.TestDidDocumentReference,
	Owner:                keeper.TestOwnerAddress,
}

var testUtils = keeper.TestUtils
var handler = NewHandler(testUtils.IdKeeper)

func TestValidMsg_StoreDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgSetId)
	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized commercioid message type"))
}
