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

	"github.com/commercionetwork/commercionetwork/x/commerciomint/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
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
	amount, err := sdk.ParseCoins(args[0])
	if err != nil {
		return err
	}

	cdpMsg := types.MsgOpenCdp{
		Depositor:       sender,
		DepositedAmount: amount,
	}

	err = cdpMsg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{cdpMsg})
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

	cdpMsg := types.MsgCloseCdp{
		Signer:    sender,
		Timestamp: timestamp,
	}

	err = cdpMsg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{cdpMsg})
}
