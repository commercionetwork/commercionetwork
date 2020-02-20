package memberships

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, k Keeper, accountKeeper auth.AccountKeeper) []abci.ValidatorUpdate {
	ii := k.InvitesIterator(ctx)
	defer ii.Close()

	for ; ii.Valid(); ii.Next() {
		var invite types.Invite
		k.Cdc.MustUnmarshalBinaryBare(ii.Value(), &invite)

		if invite.Status != types.InviteStatusPending {
			continue
		}

		if accountKeeper.GetAccount(ctx, invite.User) != nil {
			// found a non-rewarded invite (still pending) that exists on the network, removing it!
			invite.Status = types.InviteStatusInvalid
			k.SaveInvite(ctx, invite)
		}
	}

	return []abci.ValidatorUpdate{}
}
