package cli_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/testutil/network"
	"github.com/commercionetwork/commercionetwork/x/documents/client/cli"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

var governmentGenesisState = govTypes.GenesisState{
	GovernmentAddress: "cosmos1wze8mn5nsgl9qrgazq6a92fvh7m5e6psjcx2du",
}
var documentsGenesisState = types.GenesisState{
	Documents: []*types.Document{&types.ValidDocument},
	Receipts: []*types.DocumentReceipt{
		&types.ValidDocumentReceiptRecipient1,
		&types.ValidDocumentReceiptRecipient2,
	},
}

var ctx client.Context

func TestTx(t *testing.T) {

	cfg := network.DefaultConfig()

	bufGov, err := cfg.Codec.MarshalJSON(&governmentGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[govTypes.ModuleName] = bufGov

	buf, err := cfg.Codec.MarshalJSON(&documentsGenesisState)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	net := network.New(t, cfg)
	val := net.Validators[0]
	ctx = val.ClientCtx

	t.Run("CmdShowDocument", testCmdShowDocument)
	t.Run("CmdSentDocuments", testCmdSentDocuments)

}

func testCmdShowDocument(t *testing.T) {

	for _, tt := range []struct {
		name     string
		args     []string
		expected *types.Document
		wantErr  bool
	}{
		{
			name:     "ok",
			args:     []string{types.ValidDocument.UUID},
			expected: &types.ValidDocument,
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
				var response types.QueryGetDocumentResponse
				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &response))
				require.Equal(t, tt.expected, response.Document)
			}
		})
	}
}

func testCmdSentDocuments(t *testing.T) {

	for _, tt := range []struct {
		name     string
		args     []string
		expected []*types.Document
		wantErr  bool
	}{
		{
			name:     "ok",
			args:     []string{types.ValidDocument.Sender},
			expected: documentsGenesisState.Documents,
		},
		{
			name:    "invalid address",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name: "no documents expected",
			args: []string{types.ValidDocumentReceiptRecipient1.Sender},
		},
		{
			name:    "no args",
			args:    []string{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSentDocuments(), tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var response types.QueryGetSentDocumentsResponse
				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &response))
				require.ElementsMatch(t, tt.expected, response.Document)
			}
		})
	}
}
