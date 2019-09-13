package memberships

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgSetId = MsgAssignMembership{
	Signer:         keeper.TestSignerAddress,
	User:           keeper.TestUserAddress,
	MembershipType: keeper.TestMembershipType,
}

func TestHandler_ValidMsgAssignMembership(t *testing.T) {
	_, ctx, k := TestSetup()
	k.AddTrustedMinter(ctx, keeper.TestSignerAddress)
	var handler = NewHandler(k)
	res := handler(ctx, msgSetId)
	require.True(t, res.IsOK())
}

func TestHandler_InvalidUnknownType(t *testing.T) {
	_, ctx, k := TestSetup()
	k.AddTrustedMinter(ctx, keeper.TestSignerAddress)
	var handler = NewHandler(k)
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}

func TestHandler_InvalidMembershipType(t *testing.T) {
	_, ctx, k := TestSetup()
	k.AddTrustedMinter(ctx, keeper.TestSignerAddress)
	var handler = NewHandler(k)
	types := []string{"gren", "bronz", "slver", "gol", "blck"}
	for _, memType := range types {
		msg := MsgAssignMembership{
			Signer:         keeper.TestSignerAddress,
			User:           keeper.TestUserAddress,
			MembershipType: memType,
		}
		res := handler(ctx, msg)
		require.False(t, res.IsOK())
		require.True(t, strings.Contains(res.Log, fmt.Sprintf("Invalid membership type: %s", memType)))
	}
}

func TestHandler_MembershipUpgrade(t *testing.T) {
	_, ctx, k := TestSetup()
	k.AddTrustedMinter(ctx, keeper.TestSignerAddress)
	var handler = NewHandler(k)
	types := []string{"green", "bronze", "silver", "gold", "black"}

	for index := 1; index < len(types); index++ {
		beforeType := types[index-1]
		memType := types[index]

		_, _ = k.AssignMembership(ctx, keeper.TestUserAddress, beforeType)

		msg := MsgAssignMembership{
			Signer:         keeper.TestSignerAddress,
			User:           keeper.TestUserAddress,
			MembershipType: memType,
		}
		res := handler(ctx, msg)
		require.True(t, res.IsOK())
	}
}

func TestHandler_InvalidMembershipUpgrade(t *testing.T) {
	_, ctx, k := TestSetup()
	k.AddTrustedMinter(ctx, keeper.TestSignerAddress)
	var handler = NewHandler(k)
	types := []string{"green", "bronze", "silver", "gold", "black"}

	for _, memType := range types {
		_, _ = k.AssignMembership(ctx, keeper.TestUserAddress, memType)

		msg := MsgAssignMembership{
			Signer:         keeper.TestSignerAddress,
			User:           keeper.TestUserAddress,
			MembershipType: memType,
		}
		res := handler(ctx, msg)
		require.False(t, res.IsOK())
		require.True(t, strings.Contains(res.Log, fmt.Sprintf("Cannot upgrade from %s membership to %s", memType, memType)))
	}
}
