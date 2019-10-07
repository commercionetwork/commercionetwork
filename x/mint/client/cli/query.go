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

	cmd.AddCommand(getCmdGetCdps(cdc), getCmdGetCdp(cdc))

	return cmd
}

func getCmdGetCdp(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "Cdp [owner-address] [timestamp]",
		Short: "Get the Cdp associated with the given address and timestamp",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryGetCdp, args[0], args[1])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get Cdp with %s timestamp", args[0])
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdGetCdps(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "Cdps [owner-address]",
		Short: "Get all the Cdps associated with the given address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCdps, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get Cdps of %s", args[0])
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
