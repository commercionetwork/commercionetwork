package cli
/*
import (
    "strconv"
	"github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
)

var _ = strconv.Itoa(0)

func CmdDeleteUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-upgrade [proposer]",
		Short: "Broadcast message DeleteUpgrade",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
      argsProposer := string(args[0])
      
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteUpgrade(clientCtx.GetFromAddress().String(), string(argsProposer))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

    return cmd
}*/