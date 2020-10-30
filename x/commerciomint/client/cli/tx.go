package cli

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CommercioMINT transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		mintCCCCmd(cdc),
		burnCCCCmd(cdc),
		setConversionRateCmd(cdc),
	)

	return txCmd
}

func mintCCCCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [amount]",
		Short: "Mints a given amount of CCC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return mintCCCCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func mintCCCCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	// TODO: implicit uccc instead of a coin
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	sender := cliCtx.GetFromAddress()
	deposit, err := sdk.ParseCoins(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgMintCCC(sender, deposit)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func burnCCCCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [id] [amount]",
		Short: "Burns a given amount of tokens, associated with id.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return burnCCCCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func burnCCCCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	sender := cliCtx.GetFromAddress()
	id := args[0]

	amount, err := sdk.ParseCoin(args[1])
	if err != nil {
		return err
	}

	msg := types.NewMsgBurnCCC(sender, id, amount)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func setConversionRateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-conversion-rate [rate]",
		Short: "Sets conversion rate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConversionRateCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func setConversionRateCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	signer := cliCtx.GetFromAddress()
	rate, ok := sdk.NewIntFromString(args[0])
	if !ok {
		return fmt.Errorf("cannot parse collateral rate, must be an integer")
	}
	msg := types.NewMsgSetCCCConversionRate(signer, rate)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
