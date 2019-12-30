package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getCmdSentDocuments(cdc),
		getCmdReceivedDocuments(cdc),
		getCmdSentReceipts(cdc),
		getCmdReceivedReceipts(cdc),
		getCmdSupportedMetadataSchemes(cdc),
		getCmdMetadataSchemesProposers(cdc),
	)

	return cmd
}

// ----------------------------------
// --- Documents
// ----------------------------------

func getCmdSentDocuments(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sent-documents [user-address]",
		Short: "Get all documents sent by user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdSentDocumentsFunc(cmd, args, cdc)
		},
	}
}

func getCmdSentDocumentsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentDocuments, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		fmt.Printf("Could not get sent documents by user: \n %s", err)
	}

	fmt.Println(string(res))

	return nil
}

func getCmdReceivedDocuments(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "received-documents [user-address]",
		Short: "Get all documents received by user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdReceivedDocumentsFunc(cmd, args, cdc)
		},
	}
}

func getCmdReceivedDocumentsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryReceivedDocuments, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		fmt.Printf("Could not get received documents by user: \n %s", err)
	}

	fmt.Println(string(res))

	return nil
}

// ----------------------------------
// --- Document metadata schemes
// ----------------------------------

func getCmdSupportedMetadataSchemes(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "metadata-schemes",
		Short: "Get all the supported metadata schemes",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdSupportedMetadataSchemesFunc(cmd, args, cdc)
		},
	}
}

func getCmdSupportedMetadataSchemesFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySupportedMetadataSchemes)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		fmt.Printf("Could not get supported metadata schemes: \n %s", err)
	}

	fmt.Println(string(res))

	return nil
}

// -----------------------------------------
// --- Document metadata schemes proposers
// -----------------------------------------

func getCmdMetadataSchemesProposers(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "metadata-schemes-proposers",
		Short: "Get all the metadata schemes proposers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdMetadataSchemesProposersFunc(cmd, args, cdc)
		},
	}
}

func getCmdMetadataSchemesProposersFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTrustedMetadataProposers)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		fmt.Printf("Could not get metadata proposers: \n %s", err)
	}

	fmt.Println(string(res))

	return nil
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func getCmdSentReceipts(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sent-receipts [user-address]",
		Short: "Get all the receipts sent from the specified user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdSentReceiptsFunc(cmd, args, cdc)
		},
	}
}

func getCmdSentReceiptsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentReceipts, args[0])
	res, _, err2 := cliCtx.QueryWithData(route, nil)
	if err2 != nil {
		fmt.Printf("Could not get any sent receipt for the given user: \n %s", err2)
	}

	fmt.Print(string(res))
	return nil
}

func getCmdReceivedReceipts(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "received-receipts [user-address] [[doc-uuid]]",
		Short: "Get the document receipt associated with given document uuid",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdReceivedReceiptsFunc(cmd, args, cdc)
		},
	}
}

func getCmdReceivedReceiptsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	addr, uuid := args[0], ""
	if len(args) == 2 {
		uuid = args[1]
	}

	route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryReceivedReceipts, addr, uuid)
	res, _, err2 := cliCtx.QueryWithData(route, nil)
	if err2 != nil {
		fmt.Printf("Could not get any received receipt for the given user or uuid: \n %s", err2)
	}

	fmt.Print(string(res))
	return nil
}
