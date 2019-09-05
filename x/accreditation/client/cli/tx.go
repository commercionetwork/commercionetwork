package cli

import (
	"github.com/commercionetwork/commercionetwork/x/accreditation/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Accreditation transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdSetAccrediter(cdc),
	)

	return txCmd
}

// GetCmdShareDocument is the CLI command for sending a ShareDocument transaction
func GetCmdSetAccrediter(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [user] [accrediter]",
		Short: "Sets the accrediter for the given user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			signer := cliCtx.GetFromAddress()

			user, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			accrediter, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetAccrediter(user, accrediter, signer)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
