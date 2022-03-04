package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/documents/client/cli"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

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
			args:    []string{types.ValidDocumentReceiptRecipient1.UUID},
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
		flags    []string
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSentDocuments(), commandArgs)
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

func testCmdReceivedDocuments(t *testing.T) {

	for _, tt := range []struct {
		name     string
		args     []string
		flags    []string
		expected []*types.Document
		wantErr  bool
	}{
		{
			name:     "ok",
			args:     []string{types.ValidDocumentReceiptRecipient1.Sender},
			expected: documentsGenesisState.Documents,
		},
		{
			name:    "invalid address",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name: "no documents expected",
			args: []string{types.ValidDocument.Sender},
		},
		// unknown flag: --page
		// missing AddPaginationFlagsToCmd in CmdReceivedDocuments
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdReceivedDocuments(), commandArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var response types.QueryGetReceivedDocumentResponse
				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &response))
				require.ElementsMatch(t, tt.expected, response.ReceivedDocument)
			}
		})
	}
}
