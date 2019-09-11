package txreward

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgIncrementsBRPool = MsgIncrementsBlockRewardsPool{
	Funder: keeper.TestFunder,
	Amount: keeper.TestAmount,
}

func TestValidMsg_IncrementBRPool(t *testing.T) {
	_, ctx, k, _, _ := TestSetup()

	handler := NewHandler(k)

	res := handler(ctx, msgIncrementsBRPool)
	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, k, _, _ := TestSetup()

	handler := NewHandler(k)

	res := handler(ctx, sdk.NewTestMsg())

	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
