package cli

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"errors"

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
		getCmdDepositIntoPool(cdc),
		getCmdGovAssignMembership(cdc),
		getCmdInviteUser(cdc),
		getCmdBuyMembership(cdc),
	)

	return txCmd
}

func getCmdVerifyUser(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-user [address]",
		Short: "Sets the given address as a verified user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

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
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdDepositIntoPool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [amount]",
		Short: "Increments the membership rewards pool's liquidity by the given amount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdDepositIntoPoolFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdDepositIntoPoolFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	funder := cliCtx.GetFromAddress()
	amount, err := sdk.ParseCoins(args[0])
	if err != nil {
		return err
	}

	for _, coin := range amount {
		if coin.Denom != "ucommercio" {
			return errors.New("only ucommercio amounts are accepted")
		}
	}

	msg := types.NewMsgDepositIntoLiquidityPool(amount, funder)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func getCmdGovAssignMembership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov-assign-black-membership [subscriber]",
		Short: "As government, assign Black membership to a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdGovAssignMembershipFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdGovAssignMembershipFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	govAddr := cliCtx.GetFromAddress()
	recipient, err := sdk.AccAddressFromBech32(args[0])

	if err != nil {
		return err
	}

	msg := types.NewMsgSetBlackMembership(recipient, govAddr)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func getCmdInviteUser(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invite-user [subscriber]",
		Short: "Invite user to buy a membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdInviteUserFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdInviteUserFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	inviter := cliCtx.GetFromAddress()
	invitee, err := sdk.AccAddressFromBech32(args[0])

	if err != nil {
		return err
	}

	msg := types.NewMsgInviteUser(inviter, invitee)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func getCmdBuyMembership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-membership [membership-type]",
		Short: "Buy a membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdBuyMembershipFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdBuyMembershipFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	buyer := cliCtx.GetFromAddress()
	membershipType := args[0]

	msg := types.NewMsgBuyMembership(membershipType, buyer)

	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
