package cli

import (
	"commercio-network/x/commerciodocs/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for commerciodocs module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdReadDocumentMetadata(cdc), GetCmdListAuthorizedReaders(cdc))

	return cmd
}

func GetCmdReadDocumentMetadata(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "metadata [document-reference]",
		Short: "Retrieves the metadata reference for the given document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/metadata/%s", types.QuerierRoute, name)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get metadata for document %s: \n %s", string(name), err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdListAuthorizedReaders(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "readers [document-reference]",
		Short: "Lists all the users that are allowed to read the document with the given reference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/readers/%s", types.QuerierRoute, name)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get readers for %s: \n %s", string(name), err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
