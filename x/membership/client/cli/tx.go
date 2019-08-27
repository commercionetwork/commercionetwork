package cli

import (
	"github.com/commercionetwork/commercionetwork/x/membership/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

/**
This files contains the functions that take a CLI command and emit a message that will later be held in order to
perform a transaction on the blockchain.
*/

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Memberships transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdAssignMembership(cdc),
	)

	return txCmd
}

// GetCmdAssignMembership is the CLI command for sending a SetIdentity transaction
func GetCmdAssignMembership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-membership [user-address] [membership-type]",
		Short: "Associates to the specified user address a membership of the given type",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			signer := cliCtx.GetFromAddress()
			user, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgAssignMembership(signer, user, args[2])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
