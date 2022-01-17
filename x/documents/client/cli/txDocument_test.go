package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/testutil/network"
	"github.com/commercionetwork/commercionetwork/x/documents/client/cli"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

func TestShareDocument(t *testing.T) {
	cfg := network.DefaultConfig()

	stateGov := govTypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[govTypes.ModuleName], &stateGov))
	stateGov.GovernmentAddress = "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae"
	bufGov, _ := cfg.Codec.MarshalJSON(&stateGov)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	net := network.New(t, cfg)
	val := net.Validators[0]
	ctx := val.ClientCtx

	fields := []string{"cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr", "6a2f41a3-c54c-fce8-32d2-0324e1c32e22", "http://www.contentUri.com",
		"http://www.contentUri.com", "test", "http://www.contentUri.com", "48656c6c6f20476f7068657221234567", "md5"}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			code: 13,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShareDocument(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}
