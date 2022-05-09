package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdFixSupply(),
	)
	return cmd
}

func CmdFixSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fixsupply [amount] [sub]",
		Short: "Fix supply to avoid invariant broken on upgrade chain",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fixSupplymdFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func fixSupplymdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	sender := cliCtx.GetFromAddress()
	amount, err := sdk.ParseCoinNormalized(args[0])
	if err != nil {
		return err
	}
	sub := false
	if len(args) == 2 {
		sub = true
	}

	msg := types.NewMsgFixSupplys(sender, amount, sub)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
