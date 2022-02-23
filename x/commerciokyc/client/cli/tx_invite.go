package cli

import (
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CmdInvite returns a root CLI command handler for all x/commerciokyc transaction commands.
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

// InviteUserFunc returns a CLI command handler for creating a MsgInviteUser transaction.
func InviteUserFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	inviter := cliCtx.GetFromAddress()
	inviteeIn := args[0]
	_, err = sdk.AccAddressFromBech32(inviteeIn)
	if err != nil {
		return err
	}

	msg := types.NewMsgInviteUser(inviter.String(), inviteeIn)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
