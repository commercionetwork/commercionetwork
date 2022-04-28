package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CmdAssignMemebrship() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-membership [subscriber] [membership]",
		Short: "As government, assign membership to a user.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return AssignMemebrshipFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func AssignMemebrshipFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	govAddr := cliCtx.GetFromAddress()
	recipient, err := sdk.AccAddressFromBech32(args[0])
	membership := args[1]

	if err != nil {
		return err
	}

	msg := types.NewMsgSetMembership(recipient.String(), govAddr.String(), membership)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

func CmdRemoveMemebrship() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-membership [subscriber]",
		Short: "As government, remove membership of a user.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RemoveMemebrshipFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func RemoveMemebrshipFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	govAddr := cliCtx.GetFromAddress()
	recipient, err := sdk.AccAddressFromBech32(args[0])

	if err != nil {
		return err
	}
	msg := types.NewMsgRemoveMembership(govAddr.String(), recipient.String())
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

func CmdAddTsp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-tsp [tsp-address]",
		Short: "Government add a tsp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return AddTspFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func AddTspFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	govAddr := cliCtx.GetFromAddress()
	tsp, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgAddTsp(tsp.String(), govAddr.String())

	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

// CmdRemoveTsp remove the given tsp from tsps list
func CmdRemoveTsp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-tsp [tsp-address]",
		Short: "Government remove a tsp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RemoveTspFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func RemoveTspFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	govAddr := cliCtx.GetFromAddress()

	tsp, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgRemoveTsp(tsp.String(), govAddr.String())

	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

// CmdDepositIntoPool add the given amount of commercio token to the pool
func CmdDepositIntoPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [amount]",
		Short: "Increments the membership rewards pool's liquidity by the given amount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return depositIntoPoolFunc(cmd, args)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func depositIntoPoolFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)

	funder := cliCtx.GetFromAddress()
	amount, err := sdk.ParseCoinsNormalized(args[0])
	if err != nil {
		return err
	}

	for _, coin := range amount {
		if coin.Denom != "ucommercio" {
			return errors.New("only ucommercio amounts are accepted")
		}
	}

	msg := types.NewMsgDepositIntoLiquidityPool(amount, funder.String())
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
