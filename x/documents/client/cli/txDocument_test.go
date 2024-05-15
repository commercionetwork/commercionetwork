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
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

const codeInsufficientFees = 13

func TestCmdShareDocument(t *testing.T) {

	cfg := network.DefaultConfig()

	bufGov, err := cfg.Codec.MarshalJSON(&governmentGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	buf, err := cfg.Codec.MarshalJSON(&documentsGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	txNet := network.New(t, cfg)
	txVal := txNet.Validators[0]
	ctx = txVal.ClientCtx

	for _, tt := range []struct {
		name    string
		fields  []string
		args    []string
		wantErr bool
		code    uint32
	}{
		{
			name: "valid",
			fields: []string{"cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr", "6a2f41a3-c54c-fce8-32d2-0324e1c32e22", "http://www.contentUri.com",
				"http://www.contentUri.com", "test", "http://www.contentUri.com", "48656c6c6f20476f7068657221234567", "md5"},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, txVal.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(txNet.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			code: codeInsufficientFees,
		},
		{
			name: "missing checksum algorithm",
			fields: []string{"cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr", "6a2f41a3-c54c-fce8-32d2-0324e1c32e22", "http://www.contentUri.com",
				"http://www.contentUri.com", "test", "http://www.contentUri.com", "48656c6c6f20476f7068657221234567"},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, txVal.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(txNet.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			wantErr: true,
		},
		{
			name: "invalid data for document",
			fields: []string{"cosmosABC", "6a2f41a3-c54c-fce8-32d2-0324e1c32e22", "http://www.contentUri.com",
				"http://www.contentUri.com", "test", "http://www.contentUri.com", "48656c6c6f20476f7068657221234567", "md5"},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, txVal.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(txNet.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commandArgs := append(tt.fields, tt.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShareDocument(), commandArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tt.code, resp.Code)
			}
		})
	}
}

func TestCmdSendDocumentReceipt(t *testing.T) {

	cfg := network.DefaultConfig()

	bufGov, err := cfg.Codec.MarshalJSON(&governmentGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	buf, err := cfg.Codec.MarshalJSON(&documentsGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	txNet := network.New(t, cfg)
	txVal := txNet.Validators[0]
	ctx = txVal.ClientCtx

	for _, tt := range []struct {
		name    string
		fields  []string
		args    []string
		wantErr bool
		code    uint32
	}{
		{
			name: "valid",
			fields: []string{"cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr", "8205AB5C69C6F167D974FDC17B51367430E7E274E5084CE9E45D949C82AB02E6", "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				"proof"},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, txVal.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(txNet.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			code: codeInsufficientFees,
		},
		{
			name: "invalid data for receipt",
			fields: []string{"cosmosABC", "8205AB5C69C6F167D974FDC17B51367430E7E274E5084CE9E45D949C82AB02E6", "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				"proof"},
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, txVal.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(txNet.Config.BondDenom, sdk.NewInt(10))).String()),
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commandArgs := append(tt.fields, tt.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSendDocumentReceipt(), commandArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tt.code, resp.Code)
			}
		})
	}
}
