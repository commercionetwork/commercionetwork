package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/testutil/network"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/client/cli"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func TestIncrementBlockRewardsPool(t *testing.T) {
	cfg := network.DefaultConfig()

	stateGov := govTypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[govTypes.ModuleName], &stateGov))
	stateGov.GovernmentAddress = "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae"
	bufGov, _ := cfg.Codec.MarshalJSON(&stateGov)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	net := network.New(t, cfg)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  bool
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				"1000" + types.BondDenom,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			err:  false,
			code: 13,
		},
		{
			desc: "invalid amount",
			args: []string{
				"" + types.BondDenom,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			err: true,
		},
		{
			desc: "zero amount",
			args: []string{
				"0" + types.BondDenom,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			err: true,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdIncrementBlockRewardsPool(), args)
			if tc.err {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}

func TestSetParams(t *testing.T) {
	cfg := network.DefaultConfig()

	stateGov := govTypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[govTypes.ModuleName], &stateGov))
	stateGov.GovernmentAddress = "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae"
	bufGov, _ := cfg.Codec.MarshalJSON(&stateGov)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	net := network.New(t, cfg)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  bool
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				"day", "0.5",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			err:  false,
			code: 13,
		},
		{
			desc: "invalid epoch identifier",
			args: []string{
				"year", "0.5",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			err: true,
		},
		{
			desc: "invalid earn rate",
			args: []string{
				"day", "",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			err: true,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSetParams(), args)
			if tc.err {
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}
