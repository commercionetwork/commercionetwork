package types

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPemPublicKeyFromPemPrivKeyFile(t *testing.T) {
	privKey, err := LoadRSAPrivKeyFromDisk("testfiles/priv_key.pem")
	require.NoError(t, err)

	pemPubKey := PublicKeyToPemString(&privKey.PublicKey)

	expectePemPubKey, err := ioutil.ReadFile("testfiles/pub_key.pem")
	require.NoError(t, err)

	require.Equal(t, string(expectePemPubKey), pemPubKey)
}
