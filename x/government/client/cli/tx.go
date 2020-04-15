package cli

import (
	"bufio"

	"github.com/commercionetwork/commercionetwork/x/government/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "government transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		getCmdSetTumblerAddress(cdc),
	)

	return txCmd
}

func getCmdSetTumblerAddress(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-tumbler-address [address]",
		Short: "Let government set the tumbler address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdSetTumblerAddressFunc(cmd, cdc, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdSetTumblerAddressFunc(cmd *cobra.Command, cdc *codec.Codec, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	governmentAddr := cliCtx.GetFromAddress()
	recipient, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgSetTumblerAddress(governmentAddr, recipient)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
