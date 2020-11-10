package cli

import (
	"bufio"
	"errors"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
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
		getCmdDepositIntoPool(cdc),
		getCmdGovAssignMembership(cdc),
		getCmdGovRemoveMembership(cdc),
		getCmdInvite(cdc),
		getCmdBuy(cdc),
		getCmdAddTsp(cdc),
		getCmdRemoveTsp(cdc),
	)

	return txCmd
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
		Use:   "gov-assign-membership [subscriber] [membership]",
		Short: "As government, assign membership to a user. Membership \"none\" to remove membership",
		Args:  cobra.ExactArgs(2),
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

	membership := args[1]

	if err != nil {
		return err
	}
	if membership == "none" {
		msgRemove := types.NewMsgRemoveMembership(govAddr, recipient)
		err = msgRemove.ValidateBasic()
		if err != nil {
			return err
		}
		return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msgRemove})

	}
	msgSet := types.NewMsgSetMembership(recipient, govAddr, membership)
	err = msgSet.ValidateBasic()
	if err != nil {
		return err
	}
	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msgSet})

}

func getCmdGovRemoveMembership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov-remove-membership [subscriber]",
		Short: "As government, remove membership of a user.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdGovRemoveMembershipFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdGovRemoveMembershipFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	govAddr := cliCtx.GetFromAddress()
	recipient, err := sdk.AccAddressFromBech32(args[0])

	if err != nil {
		return err
	}
	msg := types.NewMsgRemoveMembership(govAddr, recipient)
	err = msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func getCmdInvite(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invite [subscriber]",
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

func getCmdBuy(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [subscriber] [membership-type]",
		Short: "Tsp buy a membership for subscriber",
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

	tsp := cliCtx.GetFromAddress()

	buyer, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}
	membershipType := args[1]

	msg := types.NewMsgBuyMembership(membershipType, buyer, tsp)

	err2 := msg.ValidateBasic()
	if err2 != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func getCmdAddTsp(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-tsp [tsp-address]",
		Short: "Government add a tsp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdAddTspFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdAddTspFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	govAddr := cliCtx.GetFromAddress()

	tsp, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgAddTsp(tsp, govAddr)

	err2 := msg.ValidateBasic()
	if err2 != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func getCmdRemoveTsp(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-tsp [tsp-address]",
		Short: "Government remove a tsp",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdRemoveTspFunc(cdc, cmd, args)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func getCmdRemoveTspFunc(cdc *codec.Codec, cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	govAddr := cliCtx.GetFromAddress()

	tsp, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	msg := types.NewMsgRemoveTsp(tsp, govAddr)

	err2 := msg.ValidateBasic()
	if err2 != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
