package v3_0_0

import (
	"reflect"
	"testing"
	"time"

	v220did "github.com/commercionetwork/commercionetwork/x/did/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// package initialization for correct validation of commercionetwork addresses
func init() {
	configTestPrefixes()
}

func configTestPrefixes() {
	AccountAddressPrefix := "did:com:"
	AccountPubKeyPrefix := AccountAddressPrefix + "pub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.Seal()
}

const (
	rsaVerificationKey2018 = "-----BEGIN PUBLIC KEY----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
	rsaSignature2018       = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHg\nkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScAD\nG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUc\nzhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7\nZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0\nO2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfK\nHQIDAQAB\n-----END PUBLIC KEY-----\n"

	validBase64RsaVerificationKey2018 = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB"
	validBase64RsaSignature2018       = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHgkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScADG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUczhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7ZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0O2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfKHQIDAQAB"

	validDidSubject      = "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc"
	validDidOtherSubject = "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"

	validDateCreated = "2019-03-23T06:35:22Z"
)

var (
	validContext = []string{
		types.ContextDidV1,
		"https://w3id.org/security/suites/ed25519-2018/v1",
		"https://w3id.org/security/suites/x25519-2019/v1",
	}

	validVerificationMethodRsaVerificationKey2018 = types.VerificationMethod{
		ID:                 validDidSubject + types.RsaVerificationKey2018NameSuffix,
		Type:               types.RsaVerificationKey2018,
		Controller:         validDidSubject,
		PublicKeyMultibase: string(types.MultibaseCodeBase64) + validBase64RsaVerificationKey2018,
	}

	validVerificationMethodRsaSignature2018 = types.VerificationMethod{
		ID:                 validDidSubject + types.RsaSignature2018NameSuffix,
		Type:               types.RsaSignature2018,
		Controller:         validDidSubject,
		PublicKeyMultibase: string(types.MultibaseCodeBase64) + validBase64RsaSignature2018,
	}

	validVerificationMethods = []*types.VerificationMethod{
		&validVerificationMethodRsaVerificationKey2018,
		&validVerificationMethodRsaSignature2018,
	}

	validServiceBar = types.Service{
		ID:              "https://bar.example.com",
		Type:            "agent",
		ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
	}

	validServiceFoo = types.Service{
		ID:              "https://foo.example.com",
		Type:            "xdi",
		ServiceEndpoint: "https://commerc.io/xdi/serviceEndpoint/",
	}

	validServices = []*types.Service{
		&validServiceBar,
		&validServiceFoo,
	}

	validDidDocument = types.DidDocument{
		Context:            validContext,
		ID:                 validDidSubject,
		VerificationMethod: validVerificationMethods,
		Authentication: []string{
			validDidSubject + types.RsaVerificationKey2018NameSuffix,
		},
		AssertionMethod: []string{
			validDidSubject + types.RsaSignature2018NameSuffix,
		},
		KeyAgreement: []string{
			types.RsaVerificationKey2018NameSuffix,
		},
		CapabilityInvocation: []string{
			types.RsaSignature2018NameSuffix,
		},
		CapabilityDelegation: nil,
		Service:              validServices,
	}

	validMsgSetDidDocument = types.MsgSetIdentity{
		DidDocument: &validDidDocument,
	}

	validMetadata = types.Metadata{
		Created: validDateCreated,
		Updated: validDateCreated,
	}

	validIdentity = types.Identity{
		DidDocument: &validDidDocument,
		Metadata:    &validMetadata,
	}

	v220ddo                      v220did.DidDocument
	v220ddoOtherSubject          v220did.DidDocument
	v220PubKeys                  v220did.PubKeys
	expectedIdentity             types.Identity
	expectedIdentityOtherSubject types.Identity
)

func init() {
	didAccAddress, err := sdk.AccAddressFromBech32(validDidSubject)
	if err != nil {
		panic(err)
	}

	didAccAddressOtherSubject, err := sdk.AccAddressFromBech32(validDidOtherSubject)
	if err != nil {
		panic(err)
	}

	v220PubKeys = v220did.PubKeys{
		v220did.PubKey{
			ID:           validVerificationMethodRsaVerificationKey2018.ID,
			Type:         validVerificationMethodRsaVerificationKey2018.Type,
			Controller:   didAccAddress,
			PublicKeyPem: rsaVerificationKey2018,
		},
		v220did.PubKey{
			ID:           validVerificationMethodRsaSignature2018.ID,
			Type:         validVerificationMethodRsaSignature2018.Type,
			Controller:   didAccAddress,
			PublicKeyPem: rsaSignature2018,
		},
	}

	created, err := time.Parse(types.ComplaintW3CTime, validIdentity.Metadata.Created)
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
			Controller:         validIdentity.DidDocument.ID,
			VerificationMethod: validIdentity.DidDocument.ID,
			SignatureValue:     "QNB13Y7Q91tzjn4w==",
		},
		Service: v220Service,
	}

	v220ddoOtherSubject = v220ddo
	v220ddoOtherSubject.ID = didAccAddressOtherSubject

	expectedDDO := *validIdentity.DidDocument
	expectedDDO.Context = []string{validIdentity.DidDocument.Context[0]}
	expectedDDO.AssertionMethod = []string{}
	expectedDDO.Authentication = []string{}
	expectedDDO.KeyAgreement = []string{}
	expectedDDO.CapabilityInvocation = []string{}
	expectedDDO.CapabilityDelegation = []string{}

	expectedDDOOtherSubject := expectedDDO
	expectedDDOOtherSubject.ID = validDidOtherSubject

	expectedIdentity = types.Identity{
		DidDocument: &expectedDDO,
		Metadata:    validIdentity.Metadata,
	}

	expectedIdentityOtherSubject = types.Identity{
		DidDocument: &expectedDDOOtherSubject,
		Metadata:    validIdentity.Metadata,
	}

}

func Test_publicKeyPemToMultiBase(t *testing.T) {
	tests := []struct {
		name            string
		pkPem           string
		wantPkMultiBase string
	}{
		{
			"RsaVerificationKey2018",
			rsaVerificationKey2018,
			validVerificationMethodRsaVerificationKey2018.PublicKeyMultibase,
		},
		{
			"RsaSignature2018",
			rsaSignature2018,
			validVerificationMethodRsaSignature2018.PublicKeyMultibase,
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

var v220Service = v220did.Services{
	v220did.Service{
		ID:              validIdentity.DidDocument.Service[0].ID,
		Type:            validIdentity.DidDocument.Service[0].Type,
		ServiceEndpoint: validIdentity.DidDocument.Service[0].ServiceEndpoint,
	},
	v220did.Service{
		ID:              validIdentity.DidDocument.Service[1].ID,
		Type:            validIdentity.DidDocument.Service[1].Type,
		ServiceEndpoint: validIdentity.DidDocument.Service[1].ServiceEndpoint,
	},
}

func Test_convertService(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if gotServices300 := convertService(v220Service); !reflect.DeepEqual(gotServices300, validIdentity.DidDocument.Service) {
			t.Errorf("convertService() = %v, want %v", gotServices300, validIdentity.DidDocument.Service)
		}
	})
}

func Test_convertPubKeys(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if gotVerificationMethods := convertPubKeys(v220PubKeys); !reflect.DeepEqual(gotVerificationMethods, validIdentity.DidDocument.VerificationMethod) {
			t.Errorf("convertPubKeys() = %v, want %v", gotVerificationMethods, validIdentity.DidDocument.VerificationMethod)
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
	type args struct {
		oldGenState v220did.GenesisState
	}
	tests := []struct {
		name string
		args args
		want *types.GenesisState
	}{
		{
			name: "empty",
			args: args{},
			want: &types.GenesisState{
				Identities: []*types.Identity{},
			},
		},
		{
			name: "ok",
			args: args{
				oldGenState: v220did.GenesisState{
					DidDocuments: []v220did.DidDocument{v220ddo, v220ddoOtherSubject},
				},
			},
			want: &types.GenesisState{
				Identities: []*types.Identity{&expectedIdentity, &expectedIdentityOtherSubject},
			},
		},
		{
			name: "duplicate not added",
			args: args{
				oldGenState: v220did.GenesisState{
					DidDocuments: []v220did.DidDocument{v220ddo, v220ddoOtherSubject, v220ddo},
				},
			},
			want: &types.GenesisState{
				Identities: []*types.Identity{&expectedIdentity, &expectedIdentityOtherSubject},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Migrate(tt.args.oldGenState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Migrate() = %v, want %v", got, tt.want)
			}
		})
	}
}
