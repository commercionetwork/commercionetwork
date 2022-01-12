package cli

/*
import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)
*/
const ()

/*
func CmdSetConversionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-conversion-rate [rate]",
		Short: "Sets conversion rate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConversionRateCmdFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func setConversionRateCmdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)

	signer := cliCtx.GetFromAddress()
	rate, err := sdk.NewDecFromStr(args[0])
	if err != nil {
		return fmt.Errorf("cannot parse collateral rate, must be a decimal")
	}

	msg := types.NewMsgSetCCCConversionRate(signer, rate)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

func CmdSetFreezePeriod() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-freeze-period [freeze period]",
		Short: "Sets freeze period in seconds",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setFreezePeriodCmdFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func setFreezePeriodCmdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)

	signer := cliCtx.GetFromAddress()

	freezePeriod, err := time.ParseDuration(strings.TrimSpace(args[0]) + "s")
	if err != nil {
		return fmt.Errorf("cannot parse freeze period, must be an integer")
	}
	msg := types.NewMsgSetCCCFreezePeriod(signer, freezePeriod.String())
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
*/
