package cli

import (
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CmdInvite() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invite [subscriber]",
		Short: "Invite user to buy a membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return InviteUserFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func InviteUserFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	inviter := cliCtx.GetFromAddress()
	invitee, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgInviteUser(inviter.String(), invitee.String())
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
