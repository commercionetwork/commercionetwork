package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
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

	cmd.AddCommand(getCmdGetCDPs(cdc), getCmdGetCDP(cdc))

	return cmd
}

func getCmdGetCDP(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "CDP [owner-address] [timestamp]",
		Short: "Get the CDP associated with the given address and timestamp",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryGetCDP, args[0], args[1])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get CDP with %s timestamp", args[0])
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdGetCDPs(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "CDPs [owner-address]",
		Short: "Get all the CDPs associated with the given address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCDPs, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get CDPs of %s", args[0])
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
