package types

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPemPublicKeyFromPemPrivKeyFile(t *testing.T) {
	privKey, err := LoadRSAPrivKeyFromDisk(path.Join("testdata", "priv_key.pem"))
	require.NoError(t, err)

	pemPubKey := PublicKeyToPemString(&privKey.PublicKey)

	expectePemPubKey, err := ioutil.ReadFile(path.Join("testdata", "pub_key.pem"))
	require.NoError(t, err)

	require.Equal(t, string(expectePemPubKey), pemPubKey)
}
