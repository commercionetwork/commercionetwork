package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "user-membership-verifier",
		MembershipVerifiedInvariant(k))
}

func MembershipVerifiedInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// get all the users with membership
		users, err := k.GetMembershipsSet(ctx)
		if err != nil {
			panic(err)
		}

		for _, user := range users {
			credentials := k.GetUserCredentials(ctx, user.Owner)

			// check that the user has been invited
			_, found := k.GetInvite(ctx, user.Owner)
			if !found {
				return sdk.FormatInvariant(
					types.ModuleName,
					"user-membership-verifier",
					fmt.Sprintf(
						"found user with membership but with no invite: %s",
						user.Owner.String(),
					),
				), false
			}

			// check that there are credentials for user
			if len(credentials) == 0 {
				return sdk.FormatInvariant(
					types.ModuleName,
					"user-membership-verifier",
					fmt.Sprintf(
						"found user with membership but with no credentials: %s",
						user.Owner.String(),
					),
				), false
			}

			// for each credential, check that the Verifier is actually
			// a tsp
			for _, credential := range credentials {
				if !k.IsTrustedServiceProvider(ctx, credential.Verifier) {
					return sdk.FormatInvariant(
						types.ModuleName,
						"user-membership-verifier",
						fmt.Sprintf(
							"found user whose credential was verified by a non-Verifier %s user but with no credentials: %s",
							credential.Verifier.String(),
							user.Owner.String(),
						),
					), false
				}
			}
		}

		return "", true
	}
}
