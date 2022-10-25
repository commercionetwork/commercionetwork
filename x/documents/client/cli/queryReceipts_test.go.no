package cli_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/client/cli"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
)

func testCmdSentReceipts(t *testing.T) {

	for _, tt := range []struct {
		name     string
		args     []string
		flags    []string
		expected []*types.DocumentReceipt
		wantErr  bool
	}{
		{
			name: "ok",
			args: []string{types.ValidDocumentReceiptRecipient1.Sender},
			expected: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
			},
		},
		{
			name:    "invalid address",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name: "no receipts expected",
			args: []string{types.ValidDocument.Sender},
		},
		// unknown flag: --page
		// missing AddPaginationFlagsToCmd in CmdSentDocuments
		// {
		// 	name: "invalid pagination flags",
		// 	args: []string{types.ValidDocumentReceiptRecipient1.Sender},
		// 	flags: []string{
		// 		fmt.Sprintf("--%s=2", flags.FlagPage),
		// 		fmt.Sprintf("--%s=true", flags.FlagOffset),
		// 	},
		// },
		{
			name:    "no args",
			args:    []string{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commandArgs := append(tt.args, tt.flags...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSentReceipts(), commandArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var response types.QueryGetSentDocumentsReceiptsResponse
				require.NoError(t, ctx.JSONCodec.UnmarshalJSON(out.Bytes(), &response))
				require.ElementsMatch(t, tt.expected, response.Receipt)
			}
		})
	}
}

func testCmdReceivedReceipts(t *testing.T) {

	for _, tt := range []struct {
		name     string
		args     []string
		flags    []string
		expected []*types.DocumentReceipt
		wantErr  bool
	}{
		{
			name:     "ok",
			args:     []string{types.ValidDocument.Sender},
			expected: documentsGenesisState.Receipts,
		},
		{
			name:    "invalid address",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name: "no receipts expected",
			args: []string{types.ValidDocumentReceiptRecipient1.Sender},
		},
		// unknown flag: --page
		// missing AddPaginationFlagsToCmd in CmdSentDocuments
		// {
		// 	name: "invalid pagination flags",
		// 	args: []string{types.ValidDocumentReceiptRecipient1.Sender},
		// 	flags: []string{
		// 		fmt.Sprintf("--%s=2", flags.FlagPage),
		// 		fmt.Sprintf("--%s=true", flags.FlagOffset),
		// 	},
		// },
		{
			name:    "no args",
			args:    []string{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commandArgs := append(tt.args, tt.flags...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdReceivedReceipts(), commandArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var response types.QueryGetReceivedDocumentsReceiptsResponse
				require.NoError(t, ctx.JSONCodec.UnmarshalJSON(out.Bytes(), &response))
				require.ElementsMatch(t, tt.expected, response.ReceiptReceived)
			}
		})
	}
}

func testCmdDocumentsReceipts(t *testing.T) {

	for _, tt := range []struct {
		name     string
		args     []string
		flags    []string
		expected []*types.DocumentReceipt
		wantErr  bool
	}{
		{
			name:     "ok",
			args:     []string{types.ValidDocument.UUID},
			expected: documentsGenesisState.Receipts,
		},
		// {
		// 	name:    "invalid uuid",
		// 	args:    []string{"abc"},
		// 	wantErr: true,
		// },
		{
			name: "no receipts expected",
			args: []string{types.ValidDocumentReceiptRecipient1.Sender},
		},
		// unknown flag: --page
		// missing AddPaginationFlagsToCmd in CmdSentDocuments
		// {
		// 	name: "invalid pagination flags",
		// 	args: []string{types.ValidDocumentReceiptRecipient1.Sender},
		// 	flags: []string{
		// 		fmt.Sprintf("--%s=2", flags.FlagPage),
		// 		fmt.Sprintf("--%s=true", flags.FlagOffset),
		// 	},
		// },
		{
			name:    "no args",
			args:    []string{},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commandArgs := append(tt.args, tt.flags...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdDocumentsReceipts(), commandArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var response types.QueryGetDocumentsReceiptsResponse
				require.NoError(t, ctx.JSONCodec.UnmarshalJSON(out.Bytes(), &response))
				require.ElementsMatch(t, tt.expected, response.Receipts)
			}
		})
	}
}
