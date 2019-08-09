package cli

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for commerciodocs module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdSentDocuments(cdc), GetCmdReceivedDocuments(cdc))

	return cmd
}

func GetCmdRetrieveDocument(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "document [document-checksum-value]",
		Short: "Get the document information for the given checksum value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			checksumValue := args[0]

			route := fmt.Sprintf("custom/%s/document/%s", types.QuerierRoute, checksumValue)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get document from checksum %s: \n %s", checksumValue, err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdSentDocuments(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sent-documents",
		Short: "Get all documents sent by user",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/sent", types.QuerierRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get sent documents by user: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdReceivedDocuments(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "received-documents",
		Short: "Get all documents received by user",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/received", types.QuerierRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get received documents by user: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdSharedDocuments(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "shared-with [user-address]",
		Short: "Get all documents shared with given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				fmt.Printf("The given address doesn't exist or it's wrong")
			}

			route := fmt.Sprintf("custom/%s/shared-with/%s", types.QuerierRoute, address)
			res, _, err2 := cliCtx.QueryWithData(route, nil)
			if err2 != nil {
				fmt.Printf("Could not get any shared document with %s: \n %s", string(address), err2)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
