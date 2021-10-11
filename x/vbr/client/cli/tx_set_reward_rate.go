package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

)

var _ = strconv.Itoa(0)

func CmdSetRewardRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-reward-rate",
		Short: "Broadcast message SetRewardRate",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			govArg := sdk.DecProto{}

			msg := types.NewMsgSetRewardRate(clientCtx.GetFromAddress().String(), govArg )
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
