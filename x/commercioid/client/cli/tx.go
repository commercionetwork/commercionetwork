package cli

import (
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"
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
		Short:                      "CommercioID transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdSetIdentity(cdc),
	)

	return txCmd
}

// GetCmdSetIdentity is the CLI command for sending a SetIdentity transaction
func GetCmdSetIdentity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-identity [did-document-uri]",
		Short: "Associates the given did document reference to your Did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//todo check if we must need NewCLIContextWithFrom
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			account := cliCtx.GetFromAddress()

			msg := types.NewMsgSetIdentity(account, args[0])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
