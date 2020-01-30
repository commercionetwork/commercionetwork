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
	handler := NewHandler(k, bk)

	_, err := handler(ctx, msgIncrementsBRPool)
	require.NoError(t, err)
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, k, _, bk := SetupTestInput()

	handler := NewHandler(k, bk)

	_, err := handler(ctx, sdk.NewTestMsg())

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
