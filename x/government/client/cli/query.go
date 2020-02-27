package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/government/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
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
		getCmdGetGovernmentAddr(cdc),
		getCmdGetTumblerAddr(cdc),
	)

	return cmd
}

func getCmdGetGovernmentAddr(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "gov-address",
		Short: "Get the government address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdGetGovernmentAddrFunc(cmd, args, cdc)
		},
	}
}

func getCmdGetGovernmentAddrFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGovernmentAddress)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get government address: %s", err)
	}

	cmd.Println(string(res))

	return nil
}

func getCmdGetTumblerAddr(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tumbler-address",
		Short: "Get the Tumbler address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdGetTumblerAddrFunc(cmd, args, cdc)
		},
	}
}

func getCmdGetTumblerAddrFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTumblerAddress)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get Tumbler address: %s", err)
	}

	cmd.Println(string(res))

	return nil
}
