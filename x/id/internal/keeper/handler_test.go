package keeper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidMsg_StoreDoc(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	handler := NewHandler(k)
	msgSetId := types.MsgSetIdentity(TestDidDocument)
	res := handler(ctx, msgSetId)

	assert.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	handler := NewHandler(k)
	res := handler(ctx, sdk.NewTestMsg())

	assert.False(t, res.IsOK())
	assert.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
