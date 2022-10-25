package keeper

import (
	"reflect"
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SetIdentity(t *testing.T) {

	type args struct {
		identity types.Identity
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty store",
			args: args{
				identity: types.ValidIdentity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			k.SetIdentity(ctx, tt.args.identity)

			identity, err := k.GetIdentity(ctx, tt.args.identity.DidDocument.ID, tt.args.identity.Metadata.Updated)

			require.NoError(t, err)
			require.Equal(t, tt.args.identity, *identity)
		})
	}
}

func TestKeeper_GetIdentity(t *testing.T) {
	type args struct {
		address   string
		timestamp string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Identity
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				address:   types.ValidIdentity.DidDocument.ID,
				timestamp: types.ValidIdentity.Metadata.Updated,
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				address:   types.ValidIdentity.DidDocument.ID,
				timestamp: types.ValidIdentity.Metadata.Updated,
			},
			want: &types.ValidIdentity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			if tt.want != nil {
				k.SetIdentity(ctx, types.ValidIdentity)
			}

			got, err := k.GetIdentity(ctx, tt.args.address, tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.GetIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetIdentity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func identitiesAtIncreasingMoments(nWithSameID, nWithDifferentID int) []*types.Identity {

	identities := []*types.Identity{}

	timestamp, err := time.Parse(time.RFC3339, types.ValidIdentity.Metadata.Created)
	if err != nil {
		panic("could not parse time")
	}

	timestampForSameID := timestamp
	for i := 0; i < nWithSameID; i++ {
		timestampForSameID = timestampForSameID.Add(time.Hour)
		identities = append(identities, validIdentityWithMetadataUpdated(timestampForSameID))
	}

	timestampForDifferentID := timestamp
	for i := 0; i < nWithDifferentID; i++ {
		timestampForDifferentID = timestampForDifferentID.Add(time.Hour)
		identities = append(identities, validIdentityWithDifferentIDMetadataUpdated(timestampForDifferentID))
	}

	return identities
}

func validIdentityWithMetadataUpdated(timestamp time.Time) *types.Identity {
	metadataNew := types.Metadata{
		Created: types.ValidIdentity.Metadata.Created,
		Updated: timestamp.UTC().Format(types.ComplaintW3CTime),
	}
	identityNew := types.ValidIdentity
	identityNew.Metadata = &metadataNew

	return &identityNew
}

func validIdentityWithDifferentIDMetadataUpdated(timestamp time.Time) *types.Identity {

	_, _, addr := testdata.KeyTestPubAddr()
	didDocumentNew := *types.ValidIdentity.DidDocument
	didDocumentNew.ID = addr.String()

	identityNew := *validIdentityWithMetadataUpdated(timestamp)
	identityNew.DidDocument = &didDocumentNew

	return &identityNew
}

func TestKeeper_GetIdentityHistoryOfAddress(t *testing.T) {

	type args struct {
		ID string
	}
	tests := []struct {
		name       string
		args       args
		identities []*types.Identity
	}{
		{
			name: "empty store",
			args: args{
				ID: types.ValidIdentity.DidDocument.ID,
			},
			identities: identitiesAtIncreasingMoments(0, 0),
		},
		{
			name: "one update",
			args: args{
				ID: types.ValidIdentity.DidDocument.ID,
			},
			identities: identitiesAtIncreasingMoments(1, 0),
		},
		{
			name: "three updates",
			args: args{
				ID: types.ValidIdentity.DidDocument.ID,
			},
			identities: identitiesAtIncreasingMoments(3, 0),
		},
		{
			name: "two updates, among identities with other IDs",
			args: args{
				ID: types.ValidIdentity.DidDocument.ID,
			},
			identities: identitiesAtIncreasingMoments(2, 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			expected := []*types.Identity{}
			for _, identity := range tt.identities {
				k.SetIdentity(ctx, *identity)

				if identity.DidDocument.ID == tt.args.ID {
					expected = append(expected, identity)
				}
			}

			result := k.GetIdentityHistoryOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			assert.Equal(t, expected, result)
		})
	}
}

func TestKeeper_GetAllIdentities(t *testing.T) {

	tests := []struct {
		name string
		want []*types.Identity
	}{
		{
			name: "empty",
			want: identitiesAtIncreasingMoments(0, 0),
		},
		{
			name: "one",
			want: identitiesAtIncreasingMoments(1, 0),
		},
		{
			name: "more",
			want: identitiesAtIncreasingMoments(2, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, identity := range tt.want {
				k.SetIdentity(ctx, *identity)
			}

			got := k.GetAllIdentities(ctx)

			require.Len(t, got, len(tt.want))
			for _, identity := range tt.want {
				require.Contains(t, got, identity)
			}
		})
	}
}
