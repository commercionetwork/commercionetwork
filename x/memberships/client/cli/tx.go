package cli

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Accreditations transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		getCmdVerifyUser(cdc),
	)

	return txCmd
}

func getCmdVerifyUser(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-user [address]",
		Short: "Sets the given address as a verified user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdVerifyUserFunc(cmd, args, cdc)
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func getCmdVerifyUserFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

	verifier := cliCtx.GetFromAddress()
	user, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgSetUserVerified(user, verifier)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
