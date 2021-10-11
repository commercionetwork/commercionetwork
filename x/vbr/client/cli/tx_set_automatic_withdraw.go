package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

var _ = strconv.Itoa(0)

func CmdSetAutomaticWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-automatic-withdraw",
		Short: "Broadcast message SetAutomaticWithdraw",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			autWithdrawArg := false

			msg := types.NewMsgSetAutomaticWithdraw(clientCtx.GetFromAddress().String(), autWithdrawArg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
