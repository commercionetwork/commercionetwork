package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
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

	cmd.AddCommand(GetCmdRetrieveBlockRewardsPoolFunders(cdc), GetCmdRetrieveBlockRewardsPoolFunds(cdc))

	return cmd
}

func GetCmdRetrieveBlockRewardsPoolFunders(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "funders",
		Short: "Get the authorized funders of the block rewards pool",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryBlockRewardsPoolFunders)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get any block rewards pool's funders: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func GetCmdRetrieveBlockRewardsPoolFunds(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the actual block rewards pool's total funds amount",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryBlockRewardsPoolFunds)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get total funds amount: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
