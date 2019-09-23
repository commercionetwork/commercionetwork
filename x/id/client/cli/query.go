package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec, moduleName, querierRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      "Querying commands for the commercioID module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdResolveIdentity(cdc, querierRoute))

	return cmd
}

func GetCmdResolveIdentity(cdc *codec.Codec, querierRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [did]",
		Short: "Resolves the given Did by returning the data associated with it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/identities/%s", querierRoute, name)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not resolve identity - %s \n", string(name))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
