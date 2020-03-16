package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getCmdGetInvites(cdc),
		getCmdGetInvitesForUser(cdc),
		getCmdGetTrustedServiceProviders(cdc),
		getCmdGetPoolFunds(cdc),
		getCmdMembershipForUser(cdc),
	)

	return cmd
}

func getCmdGetInvites(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-invites [user-address]",
		Short: "Get all membership invitations",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getInvitesFunc(cmd, args, cdc)
		},
	}
}

func getInvitesFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetInvites)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	cmd.Println(string(res))

	return nil
}

func getCmdGetInvitesForUser(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-invites-for-user [user-address]",
		Short: "Get all membership invitations for a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getInvitesForUserFunc(cmd, args, cdc)
		},
	}
}

func getInvitesForUserFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetInvites, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	cmd.Println(string(res))

	return nil
}

func getCmdGetTrustedServiceProviders(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "trusted-service-providers",
		Short: "Get all membership invitations for a user",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTrustedServiceProvidersFunc(cmd, args, cdc)
		},
	}
}

func getTrustedServiceProvidersFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetTrustedServiceProviders)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	cmd.Println(string(res))

	return nil
}

func getCmdGetPoolFunds(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the pool funds amounts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdGetPoolFundsFunc(cmd, cdc)
		},
	}
}

func getCmdGetPoolFundsFunc(cmd *cobra.Command, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetPoolFunds)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get pool funds schemes: %w", err)
	}

	cmd.Println(string(res))

	return nil
}

func getCmdMembershipForUser(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "user-membership [user]",
		Short: "Get user membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdMembershipForUserFunc(cmd, args, cdc)
		},
	}
}

func getCmdMembershipForUserFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetMembership, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get membership for user: %w", err)
	}

	cmd.Println(string(res))

	return nil
}
