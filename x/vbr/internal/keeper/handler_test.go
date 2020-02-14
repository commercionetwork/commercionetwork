package keeper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgIncrementsBRPool = types.MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: TestAmount,
}

func TestValidMsg_IncrementBRPool(t *testing.T) {
	_, ctx, k, _, bk := SetupTestInput()

	_ = bk.SetCoins(ctx, TestFunder, TestAmount)
	handler := NewHandler(k)

	res := handler(ctx, msgIncrementsBRPool)
	require.True(t, res.IsOK())

	macc := k.VbrAccount(ctx)

	initialPool, _ := TestBlockRewardsPool.TruncateDecimal()
	expectedTotalAmount := initialPool.Add(TestAmount)

	require.Equal(t, expectedTotalAmount, macc.GetCoins())
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	handler := NewHandler(k)

	res := handler(ctx, sdk.NewTestMsg())

	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
