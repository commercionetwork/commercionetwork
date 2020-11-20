package keeper

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

var msgIncrementsBRPool = types.MsgIncrementBlockRewardsPool{
	Funder: TestFunder,
	Amount: TestAmount,
}

func TestValidMsg_IncrementBRPool(t *testing.T) {
	_, ctx, k, _, bk, _ := SetupTestInput(false)

	_ = bk.SetCoins(ctx, TestFunder, TestAmount)
	handler := NewHandler(k)

	_, err := handler(ctx, msgIncrementsBRPool)
	require.NoError(t, err)

	macc := k.VbrAccount(ctx)

	initialPool, _ := TestBlockRewardsPool.TruncateDecimal()
	expectedTotalAmount := initialPool.Add(TestAmount...)

	require.Equal(t, expectedTotalAmount, macc.GetCoins())
}

var expectedRate = sdk.NewDecWithPrec(1, 2)
var msgSetRewardRate = types.MsgSetRewardRate{
	Government: TestDelegator,
	RewardRate: expectedRate,
}

func Test_MsgSetRewardRate(t *testing.T) {
	_, ctx, k, _, _, _ := SetupTestInput(false)
	_ = k.govKeeper.SetGovernmentAddress(ctx, TestDelegator)

	handler := NewHandler(k)
	_, err := handler(ctx, msgSetRewardRate)
	require.NoError(t, err)

	actual := k.GetRewardRate(ctx)

	require.Equal(t, expectedRate, actual)

}

var msgSetAutomaticWithdraw = types.MsgSetAutomaticWithdraw{
	Government:        TestDelegator,
	AutomaticWithdraw: false,
}

func Test_MsgSetAutomaticWithdraw(t *testing.T) {
	_, ctx, k, _, _, _ := SetupTestInput(false)
	_ = k.govKeeper.SetGovernmentAddress(ctx, TestDelegator)

	handler := NewHandler(k)
	_, err := handler(ctx, msgSetAutomaticWithdraw)
	require.NoError(t, err)

	actual := k.GetAutomaticWithdraw(ctx)

	require.False(t, actual)

}

func TestInvalidMsg(t *testing.T) {
	_, ctx, k, _, _, _ := SetupTestInput(false)

	handler := NewHandler(k)

	_, err := handler(ctx, sdk.NewTestMsg())

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
