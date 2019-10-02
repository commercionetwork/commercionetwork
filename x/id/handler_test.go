package id

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgSetId = MsgSetIdentity{
	Owner:       TestOwnerAddress,
	DidDocument: TestDidDocument,
}

func TestValidMsg_StoreDoc(t *testing.T) {
	_, ctx, govK, k := TestSetup()

	var handler = NewHandler(k, govK)
	res := handler(ctx, msgSetId)

	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, govK, k := TestSetup()

	var handler = NewHandler(k, govK)
	res := handler(ctx, sdk.NewTestMsg())

	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
