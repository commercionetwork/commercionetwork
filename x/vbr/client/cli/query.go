package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func GetQueryCmd(cdc *codec.Codec, moduleName, querierRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      "Querying commands for the VBR module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdRetrieveBlockRewardsPoolFunds(cdc, querierRoute),
		getRewardRate(cdc),
		getAutomaticWithdraw(cdc),
	)

	return cmd
}

func GetCmdRetrieveBlockRewardsPoolFunds(cdc *codec.Codec, querierRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the actual block rewards pool's total funds amount",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryBlockRewardsPoolFunds)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return fmt.Errorf("could not get total funds amount: %s", err)
			}

			cmd.Println(string(res))

			return nil
		},
	}
}

func getRewardRate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-reward-rate",
		Short: "Get the actual reward rate of vbr",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getRewardRateFunc(cdc)
		},
	}
}

func getRewardRateFunc(cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardRate)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	var rate sdk.Dec
	if err := cliCtx.Codec.UnmarshalJSON(res, &rate); err != nil {
		return err
	}

	return cliCtx.PrintOutput(rate)
}

func getAutomaticWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-automatic-withdraw",
		Short: "Get if vbr module is in automatic withdraw state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAutomaticWithdrawFunc(cdc)
		},
	}
}

func getAutomaticWithdrawFunc(cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAutomaticWithdraw)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	var auto bool
	if err := cliCtx.Codec.UnmarshalJSON(res, &auto); err != nil {
		return err
	}

	return cliCtx.PrintOutput(auto)
}
