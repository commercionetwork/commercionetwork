package cli

import (
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

const ()

func CmdSetParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-params [epoch identifier]",
		Short: "Set the commerciokyc params with check epoch remove membership identifier",
		Long:  "Example usage:\n commercionetworkd tx commerciomint set-params day --from ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gov := clientCtx.GetFromAddress()
			epochIdentifier := args[0]

			msg := types.NewMsgSetParams(gov.String(), epochIdentifier)
			if err2 := msg.ValidateBasic(); err2 != nil {
				return err2
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
