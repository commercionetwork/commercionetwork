package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCredential_Equals(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")
	credential := types.NewCredential(user, verifier, 0)

	tests := []struct {
		name     string
		first    types.Credential
		second   types.Credential
		expected bool
	}{
		{
			name:     "Different user returns false",
			first:    credential,
			expected: false,
			second:   types.NewCredential(verifier, credential.Verifier, credential.Timestamp),
		},
		{
			name:     "Different verifier returns false",
			first:    credential,
			expected: false,
			second:   types.NewCredential(credential.User, credential.User, credential.Timestamp),
		},
		{
			name:     "Different timestamp returns false",
			first:    credential,
			expected: false,
			second:   types.NewCredential(credential.User, credential.Verifier, credential.Timestamp+1),
		},
		{
			name:     "Same data returns true",
			first:    credential,
			expected: true,
			second:   credential,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.first.Equals(test.second))
		})
	}

}

func TestCredentials_Contains(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	verifier, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")

	credential1 := types.NewCredential(user, verifier, 0)
	credential2 := types.NewCredential(user, verifier, 1)

	tests := []struct {
		name             string
		credentials      types.Credentials
		credential       types.Credential
		expectedContains bool
	}{

		{
			name:             "Start-list credential returns true",
			credentials:      types.Credentials{credential1, credential2},
			credential:       credential1,
			expectedContains: true,
		},
		{
			name:             "End-list credential returns true",
			credentials:      types.Credentials{credential1, credential2},
			credential:       credential2,
			expectedContains: true,
		},
		{
			name:             "Empty list returns false",
			credentials:      types.Credentials{},
			credential:       credential1,
			expectedContains: false,
		},
		{
			name:             "Not found credentials returns false",
			credentials:      types.Credentials{credential1},
			credential:       credential2,
			expectedContains: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedContains, test.credentials.Contains(test.credential))
		})
	}
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
			require.Equal(t, test.shouldBeAppended, appended)
			require.Contains(t, result, test.credential)
		})
	}
}
