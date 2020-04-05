package cli

import (
	"bufio"
	"fmt"
	"strconv"

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
		openCDPCmd(cdc),
		closeCDPCmd(cdc),
		setCollateralRateCmd(cdc),
	)

	return txCmd
}

func openCDPCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-cdp [amount]",
		Short: "Opens a CDP for the given amount of ucommercio coins",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openCDPCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func openCDPCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	sender := cliCtx.GetFromAddress()
	deposit, err := sdk.ParseCoins(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgOpenCdp(sender, deposit)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func closeCDPCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-cdp [block-height]",
		Short: "Closes a CDP for a user, opened at a given block height",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return closeCDPCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func closeCDPCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	sender := cliCtx.GetFromAddress()

	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("timestamp must be a number")
	}

	msg := types.NewMsgCloseCdp(sender, timestamp)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func setCollateralRateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-collateral-rate [rate]",
		Short: "Set CDP collateral rate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setCollateralRateCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func setCollateralRateCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	signer := cliCtx.GetFromAddress()
	rate, err := sdk.NewDecFromStr(args[0])
	if err != nil {
		return err
	}
	msg := types.NewMsgSetCdpCollateralRate(signer, rate)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
