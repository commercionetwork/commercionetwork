package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	//sdk "github.com/cosmos/cosmos-sdk/types"
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
			//route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryBlockRewardsPoolFunds)
			params := &types.QueryGetBlockRewardsPoolFundsRequest{}
			res, err := queryClient.GetBlockRewardsPoolFunds(cmd.Context(), params)
			if err != nil {
				return fmt.Errorf("could not get total funds amount: %s", err)
			}

			//cmd.Println(string(res))
			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getRewardRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-reward-rate",
		Short: "Get the actual reward rate of vbr",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getRewardRateFunc(cmd)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getRewardRateFunc(cmd *cobra.Command) error {
	cliCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(cliCtx)
	params := &types.QueryGetRewardRateRequest{}
	//route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardRate)
	res, err := queryClient.GetRewardRate(cmd.Context(), params)
	if err != nil {
		return err
	}

	/*var rate sdk.Dec
	if err := cliCtx.Codec.UnmarshalJSON(res, &rate); err != nil {
		return err
	}*/

	return cliCtx.PrintProto(res)
}

func getAutomaticWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-automatic-withdraw",
		Short: "Get if vbr module is in automatic withdraw state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAutomaticWithdrawFunc(cmd)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getAutomaticWithdrawFunc(cmd *cobra.Command) error {
	cliCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(cliCtx)
	params := &types.QueryGetAutomaticWithdrawRequest{}

	res, err := queryClient.GetAutomaticWithdraw(cmd.Context(), params)
	if err != nil {
		return err
	}

	/*var auto bool
	if err := cliCtx.un.UnmarshalJSON(res, &auto); err != nil {
		return err
	}*/

	return cliCtx.PrintProto(res)
}

func getVbrParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-params",
		Short: "Get the actual params of vbr",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx:= client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(cliCtx)
		
			req := &types.QueryGetVbrParamsRequest{}
			res, err := queryClient.GetVbrParams(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("could not get total funds amount: %s", err)
			}

			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}