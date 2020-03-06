package types_test

import (
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestPubKey_Equals(t *testing.T) {
	controller, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")
	controller2, _ := sdk.AccAddressFromBech32("cosmos1007jzaanx5kmqnn3akgype2jseawfj80dne9t6")
	pubKey := types.NewPubKey("id", "type", controller, "hex-value")

	tests := []struct {
		name  string
		us    types.PubKey
		them  types.PubKey
		equal bool
	}{
		{
			"different id",
			pubKey,
			types.NewPubKey(pubKey.ID+"2", pubKey.Type, pubKey.Controller, pubKey.PublicKey),
			false,
		},
		{
			"different type",
			pubKey,
			types.NewPubKey(pubKey.ID, pubKey.Type+"other", pubKey.Controller, pubKey.PublicKey),
			false,
		},
		{
			"different controller",
			pubKey,
			types.NewPubKey(pubKey.ID, pubKey.Type, controller2, pubKey.PublicKey),
			false,
		},
		{
			"different pubkey",
			pubKey,
			types.NewPubKey(pubKey.ID, pubKey.Type, pubKey.Controller, pubKey.PublicKey+"a3"),
			false,
		},
		{
			"two equal pubkeys",
			pubKey,
			pubKey,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}

func TestPubKey_Validate(t *testing.T) {
	controller, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")

	tests := []struct {
		name string
		pk   types.PubKey
		want error
	}{
		{
			"invalid key id",
			types.NewPubKey("id", "type", controller, "13"),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid key id, must satisfy ^cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-[0-9]+$"),
		},
		{
			"invalid key type",
			types.NewPubKey("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-1", "type", controller, "10"),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid key type, must be either RsaVerificationKey2018, Secp256k1VerificationKey2018 or Ed25519VerificationKey2018"),
		},
		{
			"invaliod pubkey hex value",
			types.NewPubKey("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-1", "RsaVerificationKey2018", controller, "lkasd"),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid publicKeyHex value"),
		},
		{
			"valid pubkey",
			types.NewPubKey("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-1", "RsaVerificationKey2018", controller, "6369616f6369616f63"),
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				require.EqualError(t, tt.pk.Validate(), tt.want.Error())
			} else {
				require.NoError(t, tt.pk.Validate())
			}
		})
	}
}

func TestPubKeys_Equals(t *testing.T) {
	controller, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")

	first := types.NewPubKey("id-1", "type-1", controller, "hexValue-1")
	second := types.NewPubKey("id-2", "type-2", controller, "hexValue-2")

	tests := []struct {
		name  string
		us    types.PubKeys
		them  types.PubKeys
		equal bool
	}{
		{
			"two empty pubkeys",
			types.PubKeys{},
			types.PubKeys{},
			true,
		},
		{
			"two pubkeys with same elements and same element ordering",
			types.PubKeys{first, second},
			types.PubKeys{first, second},
			true,
		},
		{
			"two pubkeys with same elements but different element ordering",
			types.PubKeys{first, second},
			types.PubKeys{second, first},
			false,
		},
		{
			"two pubkeys with different elements",
			types.PubKeys{first},
			types.PubKeys{first, second},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}

func TestPubKeys_FindByID(t *testing.T) {
	controller, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")

	first := types.NewPubKey("id-1", "type-1", controller, "hexValue-1")
	second := types.NewPubKey("id-2", "type-2", controller, "hexValue-2")

	tests := []struct {
		name    string
		pubKeys types.PubKeys
		id      string
		wantPk  types.PubKey
		found   bool
	}{
		{
			"key not found in empty pubkeys",
			types.PubKeys{},
			first.ID,
			types.PubKey{},
			false,
		},
		{
			"key not found in non-empty pubkeys",
			types.PubKeys{first},
			second.ID,
			types.PubKey{},
			false,
		},
		{
			"key found in non-empty pubkeys",
			types.PubKeys{first, second},
			first.ID,
			first,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			pk, foundVal := tt.pubKeys.FindByID(tt.id)
			require.Equal(t, tt.found, foundVal)
			require.Equal(t, tt.wantPk, pk)
		})
	}
}

func TestPubKeys_HasVerificationAndSignatureKey(t *testing.T) {
	tests := []struct {
		name    string
		pubKeys types.PubKeys
		want    bool
	}{
		{
			"empty keys",
			types.PubKeys{},
			false,
		},
		{
			"keys-1 and keys-2 present, only keys-1 of proper type",
			types.PubKeys{
				types.PubKey{
					ID:   "#keys-1",
					Type: types.KeyTypeRsaVerification,
				},
				types.PubKey{
					ID:   "#keys-2",
					Type: "NotRsaSignature",
				},
			},
			false,
		},
		{
			"keys-1 and keys-2 present, only keys-2 of proper type",
			types.PubKeys{
				types.PubKey{
					ID:   "#keys-1",
					Type: "NotRsaVerification",
				},
				types.PubKey{
					ID:   "#keys-2",
					Type: types.KeyTypeRsaSignature,
				},
			},
			false,
		},
		{
			"keys-1 and keys-2 present, both proper type",
			types.PubKeys{
				types.PubKey{
					ID:   "#keys-1",
					Type: types.KeyTypeRsaVerification,
				},
				types.PubKey{
					ID:   "#keys-2",
					Type: types.KeyTypeRsaSignature,
				},
			},
			true,
		},
		{
			"keys-1 and keys-2 present, both not proper type",
			types.PubKeys{
				types.PubKey{
					ID:   "#keys-1",
					Type: "NotRsaVerification",
				},
				types.PubKey{
					ID:   "#keys-2",
					Type: "NotRsaSignature",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.pubKeys.HasVerificationAndSignatureKey())
		})
	}
}
