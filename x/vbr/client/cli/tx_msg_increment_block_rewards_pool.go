package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

var _ = strconv.Itoa(0)

func CmdMsgIncrementBlockRewardsPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msg-increment-block-rewards-pool [funder] [Amount]",
		Short: "Broadcast message msgIncrementBlockRewardsPool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFunder := string(args[0])
			argAmount := []sdk.Coin{}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgIncrementBlockRewardsPool( /*clientCtx.GetFromAddress().String(), */ argsFunder, argAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
