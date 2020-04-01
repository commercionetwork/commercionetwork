package cli

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
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
		getCdp(cdc),
		getCdps(cdc),
	)

	return cmd
}

func getCdp(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-cdp [user-addr] [block-height]",
		Short: "Get a CDP opened by a user at a given block height",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCdpFunc(cmd, args, cdc)
		},
	}
}

func getCdpFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	sender, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	timestamp, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("timestamp must be a number")
	}

	route := fmt.Sprintf("custom/%s/%s/%s/%d", types.QuerierRoute, types.QueryGetCdp, sender.String(), timestamp)
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func getCdps(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-cdps [user-addr]",
		Short: "Get all opened CDPs for an user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCdpsFunc(cmd, args, cdc)
		},
	}
}

func getCdpsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	sender, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetCdps, sender)
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func getCdpCollateralRate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "collateral-rate",
		Short: "Display the current Cdp collateral rate",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return getCdpCollateralRateFunc(cdc)
		},
	}
}

func getCdpCollateralRateFunc(cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCollateralRate)
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
