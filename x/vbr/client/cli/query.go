package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
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

	cmd.AddCommand(GetCmdRetrieveBlockRewardsPoolFunds(cdc, querierRoute))

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
