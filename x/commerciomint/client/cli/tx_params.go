package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdSetParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-params [conversion-rate] [freeze-period]",
		Short: "Set the commerciomint params with conversion rate and freeze-period in seconds",
		Long:  "Example usage:\n commercionetworkd tx commerciomint set-params 0.5 1814400 --from ",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gov := clientCtx.GetFromAddress()

			rate, err := math.LegacyNewDecFromStr(args[0])
			if err != nil {
				return fmt.Errorf("cannot parse collateral rate, must be a decimal")
			}

			freezePeriod, err := time.ParseDuration(strings.TrimSpace(args[1]) + "s")
			if err != nil {
				return fmt.Errorf("cannot parse freeze period, must be an integer")
			}

			msg := types.NewMsgSetParams(gov.String(), rate, freezePeriod)
			if err2 := msg.ValidateBasic(); err2 != nil {
				return err2
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
