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

var testUtils = keeper.TestUtils

var handler = NewHandler(testUtils.TBRKeeper)

func TestValidMsg_IncrementBRPool(t *testing.T) {
	res := handler(testUtils.Ctx, msgIncrementsBRPool)
	require.True(t, res.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	res := handler(testUtils.Ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
