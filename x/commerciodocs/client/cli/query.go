package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdReadDocumentMetadata(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "metadata [document-reference]",
		Short: "Retrieves the metadata reference for the given document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/metadata/%s", queryRoute, name)
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

func GetCmdListAuthorizedReaders(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "readers [document-reference]",
		Short: "Lists all the users that are allowed to read the document with the given reference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/readers/%s", queryRoute, name)
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
