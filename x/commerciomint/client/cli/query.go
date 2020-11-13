package cli

import (
	"fmt"

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
		getEtps(cdc),
		getConversionRate(cdc),
	)

	return cmd
}

func getEtps(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-etps [user-addr]",
		Short: "Get all opened ETPs for an user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getEtpsFunc(args, cdc)
		},
	}
}

func getEtpsFunc(args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	sender, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetEtps, sender)
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func getConversionRate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "conversion-rate",
		Short: "Display the current conversion rate",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return getConversionRateFunc(cdc)
		},
	}
}

func getConversionRateFunc(cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryConversionRate)
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
