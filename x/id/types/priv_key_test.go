package types

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPemPublicKeyFromPemPrivKeyFile(t *testing.T) {
	privKey, err := LoadRSAPrivKeyFromDisk(filepath.Join("testdata", "priv_key.pem"))
	require.NoError(t, err)

	pemPubKey, err := PublicKeyToPemString(&privKey.PublicKey)
	require.NoError(t, err)

	expectePemPubKey, err := ioutil.ReadFile(filepath.Join("testdata", "pub_key.pem"))
	require.NoError(t, err)

	require.Equal(t, string(expectePemPubKey), pemPubKey)
}
