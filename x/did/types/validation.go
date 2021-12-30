package types

import (
	"crypto/x509"
	"encoding/base64"
	fmt "fmt"
	"strings"

	commons "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

func isValidDidCom(did string) error {

	if !strings.HasPrefix(did, "did:com:") {
		return errors.Errorf("invalid ID address (%s), must have 'did:com:' prefix", did)
	}

	if _, err := sdk.AccAddressFromBech32(did); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ID address (%s), %e", did, err)
	}

	return nil
}

func (s *Service) isValid() error {
	if s == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service is not defined")
	}

	// validate id
	// Required
	// A string that conforms to the rules of RFC3986 for URIs.
	if IsEmpty(s.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"id\" is required")
	}
	if !IsValidRFC3986Uri(s.ID) {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, "service field \"id\" must conform to the rules of RFC3986 for URIs")
	}
	if len(s.ID) > serviceLenghtLimitID {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, fmt.Sprint("service field \"id\" must be smaller than ", serviceLenghtLimitID, " characters"))
	}

	// validate type
	// Required
	// A string.
	// W3C recommendation: In order to maximize interoperability, the service type and its associated properties SHOULD be registered in the DID Specification Registries
	if IsEmpty(s.Type) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"type\" is required")
	}
	if len(s.Type) > serviceLenghtLimitType {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, fmt.Sprint("service field \"type\" must be smaller than ", serviceLenghtLimitType, " characters"))
	}

	// validate serviceEndpoint
	// Required
	// A string that conforms to the rules of RFC3986 for URIs.
	if IsEmpty(s.ServiceEndpoint) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}
	if !IsValidRFC3986Uri(s.ServiceEndpoint) {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, "service field \"serviceEndpoint\" must conform to the rules of RFC3986 for URIs")
	}
	if len(s.ServiceEndpoint) > serviceLenghtLimitServiceEndpoint {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, fmt.Sprint("service field \"serviceEndpoint\" must be smaller than ", serviceLenghtLimitServiceEndpoint, " characters"))
	}

	return nil
}

func (v *VerificationMethod) isValid(subject string) error {
	if v == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod is not defined")
	}

	// validate ID
	// Required
	// A string that conforms to the rules for DID URL.
	if IsEmpty(v.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" is required")
	}
	if !IsValidDIDURL(v.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" must conform to the rules for DID URL")
	}

	// validate type
	// Required
	// A string that references exactly one verification method type.
	// W3C recommendation: In order to maximize global interoperability, the verification method type SHOULD be registered in the DID Specification Registries -> https://www.w3.org/TR/did-spec-registries/#verification-method-types
	if IsEmpty(v.Type) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"type\" is required")
	}
	if !commons.Strings(verificationMethodTypes).Contains(v.Type) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"type\" is not supported")
	}

	// validate controller
	// Required
	// A string that conforms to the rules of DID Syntax.
	// commercionetwork: same as the subject i.e. the ID field of DID document
	if IsEmpty(v.Controller) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"controller\" is required")
	}
	if !IsValidDID(v.Controller) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"controller\" must conform to the rules of DID Syntax")
	}
	if v.Controller != subject {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"controller\" must be equal to the subject i.e. the ID field of DID document")
	}

	// validate publicKeyMultibase
	// A verification method MUST NOT contain multiple verification material properties for the same material. For example, expressing key material in a verification method using both publicKeyJwk and publicKeyMultibase at the same time is prohibited.
	// -> using only publicKeyMultibase
	if IsEmpty(v.PublicKeyMultibase) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"publicKeyMultibase\" is required")
	}

	// commercionetwork: keys of type RsaVerificationKey2018 must be with suffix #keys-1, and must be a valid RSA PKIX public key
	if v.Type == RsaVerificationKey2018 {
		if !strings.HasSuffix(v.ID, RsaVerificationKey2018NameSuffix) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid fields \"type\" and \"id\": keys of type "+RsaVerificationKey2018+" must be with suffix "+RsaVerificationKey2018NameSuffix)
		}
	}
	// commercionetwork: keys of type RsaSignatureKey2018 must be with suffix #keys-2, and must be a valid RSA PKIX public key
	if v.Type == RsaSignature2018 {
		if !strings.HasSuffix(v.ID, RsaSignature2018NameSuffix) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid fields \"type\" and \"id\": keys of type "+RsaSignature2018+" must be with suffix "+RsaSignature2018NameSuffix)
		}
	}

	if v.Type == RsaVerificationKey2018 || v.Type == RsaSignature2018 {
		if v.PublicKeyMultibase[0] != MultibaseCodeBase64 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid field \"publicKeyMultibase\" must start with multibase code "+string(MultibaseCodeBase64))
		}
		if err := validateRSAPubkey([]byte(v.PublicKeyMultibase[1:])); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid field \"publicKeyMultibase\" must be a valid RSA public key: %e", err))
		}
	}

	return nil
}

func validateRSAPubkey(key []byte) error {
	pemBytes := make([]byte, base64.StdEncoding.DecodedLen(len(key)))
	_, err := base64.StdEncoding.Decode(pemBytes, key)
	if err != nil {
		return err
	}
	_, err = x509.ParsePKIXPublicKey(pemBytes)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	return nil
}
