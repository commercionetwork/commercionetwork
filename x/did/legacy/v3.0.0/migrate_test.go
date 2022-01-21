package v3_0_0

import (
	"reflect"
	"testing"
	"time"

	v220did "github.com/commercionetwork/commercionetwork/x/did/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	rsaVerificationKey2018          = "-----BEGIN PUBLIC KEY----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
	multibaseRsaVerificationKey2018 = "mMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB"
	rsaSignature2018                = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHg\nkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScAD\nG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUc\nzhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7\nZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0\nO2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfK\nHQIDAQAB\n-----END PUBLIC KEY-----\n"
	multibaseRsaSignature2018       = "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHgkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScADG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUczhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7ZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0O2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfKHQIDAQAB"
)

func Test_publicKeyPemToMultiBase(t *testing.T) {
	tests := []struct {
		name            string
		pkPem           string
		wantPkMultiBase string
	}{
		{
			"RsaVerificationKey2018",
			rsaVerificationKey2018,
			types.ValidIdentity.DidDocument.VerificationMethod[0].PublicKeyMultibase,
		},
		{
			"RsaSignature2018",
			rsaSignature2018,
			types.ValidIdentity.DidDocument.VerificationMethod[1].PublicKeyMultibase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPkMultiBase := publicKeyPemToMultiBase(tt.pkPem); gotPkMultiBase != tt.wantPkMultiBase {
				t.Errorf("publicKeyPemToMultiBase() = %v, want %v", gotPkMultiBase, tt.wantPkMultiBase)
			}
		})
	}
}

var didAccAddress sdk.AccAddress

var v220Service = v220did.Services{
	v220did.Service{
		ID:              types.ValidIdentity.DidDocument.Service[0].ID,
		Type:            types.ValidIdentity.DidDocument.Service[0].Type,
		ServiceEndpoint: types.ValidIdentity.DidDocument.Service[0].ServiceEndpoint,
	},
	v220did.Service{
		ID:              types.ValidIdentity.DidDocument.Service[1].ID,
		Type:            types.ValidIdentity.DidDocument.Service[1].Type,
		ServiceEndpoint: types.ValidIdentity.DidDocument.Service[1].ServiceEndpoint,
	},
}

var v220ddo v220did.DidDocument

var v220PubKeys v220did.PubKeys

var expectedIdentity types.Identity

var v220GenState v220did.GenesisState

func init() {
	didAccAddress, err := sdk.AccAddressFromBech32(types.ValidIdentity.DidDocument.ID)
	if err != nil {
		panic(err)
	}

	v220PubKeys = v220did.PubKeys{
		v220did.PubKey{
			ID:           types.ValidIdentity.DidDocument.VerificationMethod[0].ID,
			Type:         types.ValidIdentity.DidDocument.VerificationMethod[0].Type,
			Controller:   didAccAddress,
			PublicKeyPem: rsaVerificationKey2018,
		},
		v220did.PubKey{
			ID:           types.ValidIdentity.DidDocument.VerificationMethod[1].ID,
			Type:         types.ValidIdentity.DidDocument.VerificationMethod[1].Type,
			Controller:   didAccAddress,
			PublicKeyPem: rsaSignature2018,
		},
	}

	created, err := time.Parse(types.ComplaintW3CTime, types.ValidIdentity.Metadata.Created)
	if err != nil {
		panic(err)
	}

	v220ddo = v220did.DidDocument{
		Context: types.ContextDidV1,
		ID:      didAccAddress,
		PubKeys: v220PubKeys,
		Proof: &v220did.Proof{
			Type:               "EcdsaSecp256k1VerificationKey2019",
			Created:            created,
			ProofPurpose:       "authentication",
			Controller:         types.ValidIdentity.DidDocument.ID,
			VerificationMethod: types.ValidIdentity.DidDocument.ID,
			SignatureValue:     "QNB13Y7Q91tzjn4w==",
		},
		Service: v220Service,
	}

	expectedDDO := *types.ValidIdentity.DidDocument
	expectedDDO.Context = []string{types.ValidIdentity.DidDocument.Context[0]}
	expectedDDO.AssertionMethod = []string{}
	expectedDDO.Authentication = []string{}
	expectedDDO.KeyAgreement = []string{}
	expectedDDO.CapabilityInvocation = []string{}
	expectedDDO.CapabilityDelegation = []string{}

	expectedIdentity = types.Identity{
		DidDocument: &expectedDDO,
		Metadata:    types.ValidIdentity.Metadata,
	}

	v220GenState = v220did.GenesisState{
		DidDocuments:    []v220did.DidDocument{v220ddo},
		PowerUpRequests: []v220did.DidPowerUpRequest{},
	}

}

func Test_convertService(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if gotServices300 := convertService(v220Service); !reflect.DeepEqual(gotServices300, types.ValidIdentity.DidDocument.Service) {
			t.Errorf("convertService() = %v, want %v", gotServices300, types.ValidIdentity.DidDocument.Service)
		}
	})
}

func Test_convertPubKeys(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if gotVerificationMethods := convertPubKeys(v220PubKeys); !reflect.DeepEqual(gotVerificationMethods, types.ValidIdentity.DidDocument.VerificationMethod) {
			t.Errorf("convertPubKeys() = %v, want %v", gotVerificationMethods, types.ValidIdentity.DidDocument.VerificationMethod)
		}
	})
}

func Test_fromDidDocumentToIdentity(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		if gotIdentity := fromDidDocumentToIdentity(v220ddo); !reflect.DeepEqual(*gotIdentity, expectedIdentity) {
			t.Errorf("fromDidDocumentToIdentity() = %v, want %v", *gotIdentity, expectedIdentity)
		}
	})
}

func TestMigrate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {

		gs := &types.GenesisState{
			Identities: []*types.Identity{&expectedIdentity},
		}

		if got := Migrate(v220GenState); !reflect.DeepEqual(got, gs) {
			t.Errorf("Migrate() = %v, want %v", got, gs)
		}
		got := Migrate(v220GenState)
		require.Equal(t, gs, got)
	})

}
