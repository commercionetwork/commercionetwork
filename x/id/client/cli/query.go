package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/cosmos/cosmos-sdk/codec"

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
		getCmdQueryPowerUpRequest(cdc),
	)

	return cmd
}

func getCmdQueryPowerUpRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "power-up-request [request-proof]",
		Short: "Get the power-up request by its proof",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdQueryPowerUpRequestFunc(cmd, args, cdc)
		},
	}
}

func getCmdQueryPowerUpRequestFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryResolvePowerUpRequest, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get power-up request by proof: %s", err)
	}

	cmd.Println(string(res))

	return nil
}
