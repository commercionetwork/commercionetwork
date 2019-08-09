package commerciodocs

import (
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var msgShareDocument = MsgShareDocument(keeper.TestingDocument)

var testUtils = keeper.TestUtils
var handler = NewHandler(testUtils.DocsKeeper)

func TestValidMsg_ShareDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgShareDocument)
	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized commerciodocs message type"))
}
