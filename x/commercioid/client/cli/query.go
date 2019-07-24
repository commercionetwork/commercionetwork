package cli

/**
This file contains the functions that returns the commands allowing the user to perform queries on the data stored
inside the blockchain.

Each query available should have an associated function here returning the proper command to execute from the CLI.

The query path should be contained inside the querier.go file too.
*/

import (
	"commercio-network/x/commercioid/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the commercioID module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdResolveIdentity(cdc), GetCmdReadConnections(cdc))
}

func GetCmdResolveIdentity(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [did]",
		Short: "Resolve identity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/identities/%s", types.QuerierRoute, name)
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

func GetCmdReadConnections(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "connections [did]",
		Short: "Lists all the connections associated to the given Did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			route := fmt.Sprintf("custom/%s/connections/%s", types.QuerierRoute, name)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get connections for %s: %s \n", string(name), err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
