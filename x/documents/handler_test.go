package documents

// import (
// 	"fmt"
// 	"strings"
// 	"testing"

// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
// 	"github.com/commercionetwork/commercionetwork/x/documents/keeper"
// 	"github.com/commercionetwork/commercionetwork/x/documents/types"
// 	"github.com/cosmos/cosmos-sdk/testutil/testdata"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/stretchr/testify/require"
// )

// func TestInvalidMsg(t *testing.T) {
// 	k := keeper.Keeper{}
// 	h := NewHandler(k)

// 	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
// 	require.Error(t, err)
// 	require.Nil(t, res)
// 	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, testdata.NewTestMsg())))
// }
