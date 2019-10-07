package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
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

	cmd.AddCommand(
		getCmdResolveIdentity(cdc, querierRoute),
		getCmdResolveDepositRequest(cdc, querierRoute),
		getCmdResolvePowerUpRequest(cdc, querierRoute),
	)

	return cmd
}

func getCmdResolveIdentity(cdc *codec.Codec, querierRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [did]",
		Short: "Resolves the given Did by returning the data associated with it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveDid, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not resolve identity - %s \n", args[0])
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdResolveDepositRequest(cdc *codec.Codec, querierRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit-request [proof]",
		Short: "Returns the deposit request having the given proof",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolveDepositRequest, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not resolve identity - %s \n", args[0])
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdResolvePowerUpRequest(cdc *codec.Codec, querierRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "power-up-request [proof]",
		Short: "Returns the power up request having the given proof",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolvePowerUpRequest, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not resolve identity - %s \n", args[0])
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
