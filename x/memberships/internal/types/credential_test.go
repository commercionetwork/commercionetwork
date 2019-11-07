package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCredential_Equals(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")
	credential := types.NewCredential(user, verifier, 0)

	assert.False(t, credential.Equals(types.NewCredential(verifier, credential.Verifier, credential.Timestamp)))
	assert.False(t, credential.Equals(types.NewCredential(credential.User, credential.User, credential.Timestamp)))
	assert.False(t, credential.Equals(types.NewCredential(credential.User, credential.Verifier, credential.Timestamp+1)))
	assert.True(t, credential.Equals(credential))
}

func TestCredentials_Contains(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	credential1 := types.NewCredential(user, verifier, 0)
	credential2 := types.NewCredential(user, verifier, 1)

	assert.True(t, types.Credentials{credential1, credential2}.Contains(credential1))
	assert.True(t, types.Credentials{credential1, credential2}.Contains(credential2))
	assert.False(t, types.Credentials{}.Contains(credential1))
	assert.False(t, types.Credentials{credential1}.Contains(credential2))
}

func TestCredentials_AppendIfMissing(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	credential1 := types.NewCredential(user, verifier, 0)
	credential2 := types.NewCredential(user, verifier, 1)

	tests := []struct {
		name             string
		credentials      types.Credentials
		credential       types.Credential
		shouldBeAppended bool
	}{
		{
			name:             "Credential is appended to empty list",
			credentials:      types.Credentials{},
			credential:       credential1,
			shouldBeAppended: true,
		},
		{
			name:             "Credential is appended when missing",
			credentials:      types.Credentials{credential1},
			credential:       credential2,
			shouldBeAppended: true,
		},
		{
			name:             "Credentials is not appended when existing",
			credentials:      types.Credentials{credential1, credential2},
			credential:       credential2,
			shouldBeAppended: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, appended := test.credentials.AppendIfMissing(test.credential)
			assert.Equal(t, test.shouldBeAppended, appended)
			assert.Contains(t, result, test.credential)
		})
	}
}
