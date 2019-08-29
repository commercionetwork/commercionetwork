package docs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgShareDocument = MsgShareDocument(keeper.TestingDocument)
var msgDocumentReceipt = MsgSendDocumentReceipt(keeper.TestingDocumentReceipt)

var testUtils = keeper.TestUtils
var handler = NewHandler(testUtils.DocsKeeper)

func TestValidMsg_ShareDoc(t *testing.T) {
	res := handler(testUtils.Ctx, msgShareDocument)
	require.True(t, res.IsOK())
}

func TestValidMsg_DocReceipt(t *testing.T) {
	res := handler(testUtils.Ctx, msgDocumentReceipt)
	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
