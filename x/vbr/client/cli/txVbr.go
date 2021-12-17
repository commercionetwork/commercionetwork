package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

)

var _ = strconv.Itoa(0)

func CmdSetAutomaticWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-automatic-withdraw [0/1]",
		Short: "Set automatic withdraw for vbr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setAutomaticWithdrawCmdFunc(cmd, args)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func setAutomaticWithdrawCmdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	sender := cliCtx.GetFromAddress()

	aWith := false
	if args[0] == "1" {
		aWith = true
	}

	msg := types.NewMsgSetAutomaticWithdraw(sender.String(), aWith)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}


func CmdSetRewardRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-reward-rate [rate]",
		Short: "Set reward rate for vbr",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setRewardRateCmdFunc(cmd, args)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func setRewardRateCmdFunc(cmd *cobra.Command, args []string) error {
	cliCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	sender := cliCtx.GetFromAddress()
	rate, err := sdk.NewDecFromStr(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgSetRewardRate(sender.String(), rate)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

func CmdIncrementBlockRewardsPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [coin-denom] [amount]",
		Short: "Increments the block rewards pool's liquidity by the given amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			funder := clientCtx.GetFromAddress()
			amout, e := sdk.ParseCoinNormalized(args[1] + args[0])
			if e != nil {
				return e
			}
			argAmount := []sdk.Coin{amout}

			msg := types.NewMsgIncrementBlockRewardsPool(funder.String(), argAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetVbrParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-params [epoch_identifier] [vbr_earn_rate]",
		Short: "Set the vbr params with epoch identifier(i.e. \"day\" and the vbr earn rate percentage(Dec))", 
		Long: "Example usage:\n commercionetworkd tx vbr set-params day 0.500000000000000000 --from ",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gov := clientCtx.GetFromAddress()
			epochIdentifier := args[0]
			vbrEarnRate, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid vbrEarnRate (%s)", err)
			}

			msg := types.NewMsgSetVbrParams(gov.String(), epochIdentifier, vbrEarnRate)
			if err2 := msg.ValidateBasic(); err2 != nil {
				return err2
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}