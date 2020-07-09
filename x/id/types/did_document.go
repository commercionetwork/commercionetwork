package types

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// DidDocument is the concrete serialization of the data model, according to a particular syntax.
// The Did Document contains attributes or claims about the DID Subject, and the DID itself is contained in
// the id property.
type DidDocument struct {
	Context string         `json:"@context"`
	ID      sdk.AccAddress `json:"id"`
	PubKeys PubKeys        `json:"publicKey"`

	// To a future reader, to mark a DidDocument field as optional:
	//  - tag it with `omitempty` if it's a simple type (i.e. not a struct)
	//  - make it a pointer if it's a complex type (i.e. a struct)

	// Proof is **NOT** optional, we need it to have omitempty/pointer to make the signature procedure more straightforward,
	// i.e. DidDocument.Validate() will check if proof is empty, and throw an error if true.
	Proof *Proof `json:"proof,omitempty"`

	Service Services `json:"service,omitempty"` // Services are optional
}

// DidDocumentUnsigned is an intermediate type used to check for proof correctness
// It is identical to a DidDocument, it's kept for logical compartimentization.
type DidDocumentUnsigned DidDocument

// SignaturePrice represent the price a Signature service asks to sign the associated document.
// Price can be nil to indicate a free signature service for the profile.
type SignaturePrice struct {
	CertificateProfile   string             `json:"certificate_profile"`
	Price                *sdk.Coin          `json:"price,omitempty"`
	MembershipMultiplier map[string]sdk.Dec `json:"membership_multiplier"`
}

// Equals checks that ss is equal to s.
func (s SignaturePrice) Equal(ss SignaturePrice) bool {
	if len(s.MembershipMultiplier) != len(ss.MembershipMultiplier) {
		return false
	}

	for k, v := range s.MembershipMultiplier {
		if !ss.MembershipMultiplier[k].Equal(v) {
			return false
		}
	}

	return s.CertificateProfile == ss.CertificateProfile &&
		s.Price.IsEqual(*ss.Price)
}

// SignaturePrices is an array of SignaturePrice
type SignaturePrices []SignaturePrice

// Price returns the SignaturePrice for a given CertificateProfile, given a Membership.
func (sp SignaturePrices) Price(cp string, m string) (*sdk.Coin, error) {
	// a nil price with a nil error represents a free signature certificate profile
	if len(sp) == 0 {
		return nil, nil
	}

	var selPrice *SignaturePrice

	for _, p := range sp {
		if strings.TrimSpace(cp) == strings.TrimSpace(p.CertificateProfile) {
			selPrice = &p
			break
		}
	}

	if selPrice == nil {
		return nil, fmt.Errorf("no price for \"%s\" certificate profile", cp)
	}

	memMul, memMulFound := selPrice.MembershipMultiplier[m]
	// if there aren't any membership multiplier or m is not included in selPrice's multipliers,
	// just return the price
	if len(selPrice.MembershipMultiplier) == 0 || !memMulFound {
		return selPrice.Price, nil
	}

	newPriceInt := memMul.MulInt(selPrice.Price.Amount).TruncateInt()

	newPrice := sdk.NewCoin(selPrice.Price.Denom, newPriceInt)

	return &newPrice, nil
}

// Equal checks that o equals to sp.
// Both arrays must be equal index-wise.
func (sp SignaturePrices) Equal(o SignaturePrices) bool {
	if len(sp) != len(o) {
		return false
	}

	for i := 0; i < len(sp); i++ {
		if !sp[i].Equal(o[i]) {
			return false
		}
	}

	return true
}

// Service represents a service type needed for DidDocument.
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`

	// W3C says: " Each service endpoint MUST have id, type, and serviceEndpoint properties, and MAY include additional properties."
	// We can add whatever property we deem important here.

	// SignaturePrices holds informations about signature price for each signature type supported by the service provider.
	SignaturePrices SignaturePrices `json:"signature_prices,omitempty"`
}

// Validate returns error when Service is not valid.
func (s Service) Validate() error {
	if s.ID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"id\" is required")
	}

	if s.Type == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"type\" is required")
	}

	if s.ServiceEndpoint == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}

	if _, err := url.Parse(s.ServiceEndpoint); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"serviceEndpoint\" does not contain a valid URL")
	}

	// One can only define signature prices when service type is "service"
	if len(s.SignaturePrices) != 0 && s.Type != SignatureService {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "signature_prices present but service type not \"signature\"")
	}

	return nil
}

// SignatureEnabled returns true if the selected service can sign documents.
func (s Services) SignatureEnabled() (Service, error) {
	for _, ss := range s {
		if ss.Type == SignatureService {
			return ss, nil
		}
	}

	return Service{}, errors.New("no signature-enabled service")
}

// Equals returns true if s is equal to otherService.
func (s Service) Equals(otherService Service) bool {
	return s.ServiceEndpoint == otherService.ServiceEndpoint &&
		s.Type == otherService.Type &&
		s.ID == otherService.ID &&
		s.SignaturePrices.Equal(otherService.SignaturePrices)
}

// Services is a slice of services.
type Services []Service

// Equals returns true if s is equal to other.
func (s Services) Equals(other Services) bool {
	if len(s) != len(other) {
		return false
	}

	for key, value := range other {
		if !s[key].Equals(value) {
			return false
		}
	}

	return true
}

// Validate checks that all the Service instance inside s are valid.
func (s Services) Validate() error {
	for i, service := range s {
		err := service.Validate()
		if err != nil {
			return fmt.Errorf("service %d validation failed: %w", i, err)
		}
	}

	return nil
}

// noDuplicates checks that no services in s have the same ID.
func (s Services) noDuplicates() error {
	switch len(s) {
	case 0, 1:
		return nil
	case 2:
		if s[0].ID == s[1].ID {
			return errors.New("services 0 and 1 have the same ID")
		}
	default:
		for i, k := range s {
			for l, m := range s {
				if i == l {
					continue
				}

				if k.ID == m.ID {
					return fmt.Errorf("services %d and %d have the same ID", i, l)
				}

			}
		}
	}

	return nil
}

// noMultipleSignature checks that there are no more than 1 service providing signature type
func (s Services) noMultipleSignature() error {
	switch len(s) {
	case 0, 1:
		return nil
	case 2:
		if (s[0].Type == SignatureService) && (s[1].Type == SignatureService) {
			return errors.New("services 0 and 1 both provide Signature service")
		}
	default:
		foundSignature := false
		for i, k := range s {
			if k.Type == SignatureService {
				if !foundSignature {
					foundSignature = true
					continue
				}

				return fmt.Errorf("service %d duplicates Signature type", i)
			}
		}
	}

	return nil
}

// Equals returns true iff didDocument and other contain the same data
func (didDocument DidDocument) Equals(other DidDocument) bool {
	return didDocument.Context == other.Context &&
		didDocument.ID.Equals(other.ID) &&
		didDocument.PubKeys.Equals(other.PubKeys) &&
		didDocument.Proof.Equals(*other.Proof) &&
		didDocument.Service.Equals(other.Service)
}

// Validate checks the data present inside this Did Document and returns an
// error if something is wrong
func (didDocument DidDocument) Validate() error {

	if didDocument.ID.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (didDocument.ID.String()))
	}

	if didDocument.Context != ContextDidV1 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid context, must be https://www.w3.org/ns/did/v1")
	}

	for _, key := range didDocument.PubKeys {
		if err := key.Validate(); err != nil {
			return err
		}

		if !didDocument.ID.Equals(key.Controller) {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Public key controller must match did document id")
		}
	}

	if err := didDocument.PubKeys.noDuplicates(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	if !didDocument.PubKeys.HasVerificationAndSignatureKey() {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "specified public keys are not in the correct format")
	}

	if didDocument.Proof == nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "proof not provided")
	}

	if err := didDocument.Proof.Validate(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("proof validation error: %s", err.Error()))
	}

	// we have some service, we should validate 'em
	if didDocument.Service != nil {
		if err := didDocument.Service.Validate(); err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnauthorized, err.Error())
		}

		// As per W3C DID spec, "The value of serviceEndpoint MUST NOT contain multiple entries with the same id."
		if err := didDocument.Service.noDuplicates(); err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}

		// Only one "signature" type service is allowed
		if err := didDocument.Service.noMultipleSignature(); err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}
	}

	if err := didDocument.VerifyProof(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, err.Error())
	}

	return nil
}

// VerifyProof verifies d's Proof against its content.
// The Proof is constructed as follows:
//  - let K be the Bech32 Account public key, embedded in the Proof "Verification Method" field
//  - let S be K transformed in a raw Secp256k1 public key
//  - let B be the SHA-256 (as defined in the FIPS 180-4) of the JSON representation of d minus the Proof field
//  - let L be the Proof Signature Value, decoded from Base64 encoding
// The Proof is verified if K.Verify(B, L) is verified.
func (didDocument DidDocument) VerifyProof() error {
	u := DidDocumentUnsigned(didDocument)

	// Explicitly zero out the Proof field.
	//
	// Here we leverage the fact that encoding/json do not encode nil pointers,
	// effectively giving us DidDocument-(Proof field).
	u.Proof = nil

	oProof := didDocument.Proof

	// get a public key object
	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, oProof.VerificationMethod)
	if err != nil {
		return err
	}

	// get a seck256k1 public key
	sk := pk.(secp256k1.PubKeySecp256k1)

	// marshal u in json
	data, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("could not marshal unsigned did document during proof verification: %w", err)
	}

	// get signature bytes from base64
	sigBytes, err := base64.StdEncoding.DecodeString(oProof.SignatureValue)
	if err != nil {
		return fmt.Errorf("could not decode base64 signature value: %w", err)
	}

	verified := sk.VerifyBytes(data, sigBytes)

	if !verified {
		return fmt.Errorf("proof signature verification failed")
	}

	return nil
}
