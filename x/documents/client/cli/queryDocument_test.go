package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/testutil/network"
	"github.com/commercionetwork/commercionetwork/x/documents/client/cli"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

func TestCmdShowDocument(t *testing.T) {
	cfg := network.DefaultConfig()

	stateGov := govTypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[govTypes.ModuleName], &stateGov))
	stateGov.GovernmentAddress = "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae"
	bufGov, _ := cfg.Codec.MarshalJSON(&stateGov)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	state := types.GenesisState{
		Documents: []*types.Document{&types.ValidDocument},
		Receipts:  []*types.DocumentReceipt{},
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	net := network.New(t, cfg)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tt := range []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "ok",
			args: []string{types.ValidDocument.UUID},
		},
		{
			name:    "document not in store",
			args:    []string{types.AnotherValidDocument.UUID},
			wantErr: true,
		},
		{
			name:    "no args",
			args:    []string{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDocument(), tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetDocumentResponse
				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, &types.ValidDocument, resp.Document)
			}
		})
	}
}
