package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func GetCmdRetrieveBlockRewardsPoolFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the actual block rewards pool's total funds amount",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			
			queryClient := types.NewQueryClient(cliCtx)
			params := &types.QueryGetBlockRewardsPoolFundsRequest{}
			res, err := queryClient.GetBlockRewardsPoolFunds(cmd.Context(), params)
			if err != nil {
				return fmt.Errorf("could not get total funds amount: %s", err)
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-params",
		Short: "Get the actual params of vbr",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx:= client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(cliCtx)
		
			req := &types.QueryGetParamsRequest{}
			res, err := queryClient.GetParams(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("could not get total funds amount: %s", err)
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}