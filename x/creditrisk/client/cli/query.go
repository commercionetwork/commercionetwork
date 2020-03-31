package cli

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(getPoolFundsCmd(cdc))
	return cmd
}

func getPoolFundsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Short: "query pool funds",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return queryPoolFunds(cmd, args, cdc)
		},
	}
}

func queryPoolFunds(cmd *cobra.Command, _ []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithInput(bufio.NewReader(cmd.InOrStdin()))

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	var funds sdk.Coins
	if err := cdc.UnmarshalJSON(res, &funds); err != nil {
		return err
	}

	cmd.Println(funds.String())

	return nil
}
