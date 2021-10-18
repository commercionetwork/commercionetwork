package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func GetCmdRetrieveBlockRewardsPoolFunds() *cobra.Command {
	return &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the actual block rewards pool's total funds amount",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
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
}

func getRewardRate() *cobra.Command {
	return &cobra.Command{
		Use:   "get-reward-rate",
		Short: "Get the actual reward rate of vbr",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getRewardRateFunc(cmd)
		},
	}
}

func getRewardRateFunc(cmd *cobra.Command) error {
	cliCtx, e := client.GetClientTxContext(cmd)
	if e != nil {
		return e
	}

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
	return &cobra.Command{
		Use:   "get-automatic-withdraw",
		Short: "Get if vbr module is in automatic withdraw state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAutomaticWithdrawFunc(cmd)
		},
	}
}

func getAutomaticWithdrawFunc(cmd *cobra.Command) error {
	cliCtx, e := client.GetClientTxContext(cmd)
	if e != nil {
		return e
	}
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
