package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestPubKey_Equals(t *testing.T) {
	controller, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")
	pubKey := types.NewPubKey("id", "type", controller, "hex-value")

	assert.False(t, pubKey.Equals(types.NewPubKey(pubKey.ID+"2", pubKey.Type, pubKey.Controller, pubKey.PublicKeyHex)))
	assert.False(t, pubKey.Equals(types.NewPubKey(pubKey.ID, pubKey.Type+"other", pubKey.Controller, pubKey.PublicKeyHex)))
	controller2, _ := sdk.AccAddressFromBech32("cosmos1007jzaanx5kmqnn3akgype2jseawfj80dne9t6")
	assert.False(t, pubKey.Equals(types.NewPubKey(pubKey.ID, pubKey.Type, controller2, pubKey.PublicKeyHex)))
	assert.False(t, pubKey.Equals(types.NewPubKey(pubKey.ID, pubKey.Type, pubKey.Controller, pubKey.PublicKeyHex+"a3")))
	assert.True(t, pubKey.Equals(pubKey))
}

func TestPubKey_Validate(t *testing.T) {
	controller, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")

	err := types.NewPubKey("id", "type", controller, "13").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invalid key id, must satisfy ^cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-[0-9]+$")

	err = types.NewPubKey("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-1", "type", controller, "10").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invalid key type, must be either RsaVerificationKey2018, Secp256k1VerificationKey2018 or Ed25519VerificationKey2018")

	err = types.NewPubKey("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-1", "RsaVerificationKey2018", controller, "lkasd").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invalid publicKeyHex value")

	err = types.NewPubKey("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc#keys-1", "RsaVerificationKey2018", controller, "6369616f6369616f63").Validate()
	assert.NoError(t, err)
}
