package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
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

	cmd.AddCommand(getCmdOraclesList(cdc), getCmdCurrentPrice(cdc), getCmdCurrentPrices(cdc), getCmdBlacklistedDenoms(cdc))

	return cmd
}

func getCmdOraclesList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "oracles",
		Short: "Get a list of all trusted oracles",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetOracles)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get oracles' list: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdCurrentPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "price [token-name]",
		Short: "Get the current price of the token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCurrentPrice, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get token's price: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdCurrentPrices(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "prices",
		Short: "Get the current prices of all tokens",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetCurrentPrices)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get tokens' price: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdBlacklistedDenoms(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "blacklisted-denoms",
		Short: "Get the current blacklisted denoms",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdBlacklistedDenomsFunc(cmd, cdc)
		},
	}
}

func getCmdBlacklistedDenomsFunc(cmd *cobra.Command, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetBlacklistedDenoms)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		cmd.Println("could not get blacklisted denoms:", err)
	}

	cmd.Println(string(res))

	return nil
}
