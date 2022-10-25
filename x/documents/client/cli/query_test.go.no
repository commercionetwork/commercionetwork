package cli_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/testutil/network"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
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

func TestQueries(t *testing.T) {

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
	t.Run("CmdReceivedDocuments", testCmdReceivedDocuments)

	t.Run("CmdSentReceipts", testCmdSentReceipts)
	t.Run("CmdReceivedReceipts", testCmdReceivedReceipts)
	t.Run("CmdDocumentsReceipts", testCmdDocumentsReceipts)
}
