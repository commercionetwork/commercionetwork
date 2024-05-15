package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	uuid "github.com/satori/go.uuid"
)

const ()

func CmdMintCCC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [amount]",
		Short: "Mints a given amount of CCC\nAmount must be an integer number.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return mintCCCCmdFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func mintCCCCmdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	sender := cliCtx.GetFromAddress()
	deposit, ok := sdk.NewIntFromString(args[0])
	if !ok {
		return fmt.Errorf("amount must be an integer")
	}

	mintUUID := uuid.NewV4().String()
	postion := types.Position{
		Owner:      sender.String(),
		Collateral: deposit.Int64(),
		ID:         mintUUID,
	}

	msg := types.NewMsgMintCCC(postion)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

func CmdBurnCCC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [id] [amount]",
		Short: "Burns a given amount of tokens, associated with id.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return burnCCCCmdFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func burnCCCCmdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	sender := cliCtx.GetFromAddress()
	id := args[0]
	amount, err := sdk.ParseCoinNormalized(args[1])
	if err != nil {
		return err
	}

	msg := types.NewMsgBurnCCC(sender, id, amount)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
