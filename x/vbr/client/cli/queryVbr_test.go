package cli_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	"github.com/commercionetwork/commercionetwork/testutil/network"
// 	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
// 	"github.com/commercionetwork/commercionetwork/x/vbr/client/cli"
// 	"github.com/commercionetwork/commercionetwork/x/vbr/types"
// 	"github.com/cosmos/cosmos-sdk/client"
// 	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
// 	codec "github.com/cosmos/cosmos-sdk/codec"
// )

// var governmentGenesisState = govTypes.GenesisState{
// 	GovernmentAddress: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
// }

// var vbrGenesisState = types.DefaultGenesis()

// var ctx client.Context

// func TestQueries(t *testing.T) {

// 	cfg := network.DefaultConfig()

// 	bufGov, err := cfg.Codec.MarshalJSON(&governmentGenesisState)
// 	require.NoError(t, err)
// 	cfg.GenesisState[govTypes.ModuleName] = bufGov

// 	buf, err := cfg.Codec.MarshalJSON(vbrGenesisState)
// 	require.NoError(t, err)
// 	cfg.GenesisState[types.ModuleName] = buf

// 	net := network.New(t, cfg)
// 	val := net.Validators[0]
// 	ctx = val.ClientCtx

// 	t.Run("getParams", test_getParams)

// }

// func test_getParams(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		args    []string
// 		wantErr bool
// 	}{
// 		{
// 			name: "ok",
// 			args: []string{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetParams(), tt.args)
// 			if tt.wantErr {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 				var response types.QueryGetParamsResponse
// 				require.NoError(t, codec.JSONCodec.UnmarshalJSON(out.Bytes(), &response))
// 				require.Equal(t, vbrGenesisState.Params, response.Params)
// 			}
// 		})
// 	}
// }
