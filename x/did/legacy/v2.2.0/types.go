package v2_2_0

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "id"
)

type DidDocument struct {
	Context string         `json:"@context" example:"https://www.w3.org/ns/did/v1"`
	ID      sdk.AccAddress `json:"id" swaggertype:"string" example:"did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf"`
	PubKeys PubKeys        `json:"publicKey"`

	// To a future reader, to mark a DidDocument field as optional:
	//  - tag it with `omitempty` if it's a simple type (i.e. not a struct)
	//  - make it a pointer if it's a complex type (i.e. a struct)

	// Proof is **NOT** optional, we need it to have omitempty/pointer to make the signature procedure more straightforward,
	// i.e. DidDocument.Validate() will check if proof is empty, and throw an error if true.
	Proof *Proof `json:"proof,omitempty"`

	Service Services `json:"service,omitempty"` // Services are optional
}

type DidPowerUpRequest struct {
	Status   *RequestStatus `json:"status"`
	Claimant sdk.AccAddress `json:"claimant"`
	Amount   sdk.Coins      `json:"amount"`
	Proof    string         `json:"proof" example:"S5hyg4slMxm9fK8PTNDs8tHmQcBfWXG0vqrNHLXY5K1qUz3QwZYjR9nzJoNDJh18aPsXper7rNBbyZPOm5K//x8Bqm2EJkdnHd7woa5eFqpziGaHxqvgPaLGspH47tnVilARTeF23L2NVHWcEWuo9U5cWg52l1lOixOG+DehT3vC9KjLqg0YqBoL2u0LTLqQMON4UUjC8JwzT/RMs30OYGsWuLc9s48RtJCQJZ+yAg3U6jZn3OokGwWWjYxF9tAsMR48KilHsPigsa9WPnaAyCMSJ05hOqjBxWiSHYiH1nAefFqHtNFXhJF3LRUCJ2xnSHxJC5Ndj4HFzUjyK4aiV1mtRlRcsqmXU80HEk7IzI74HYpW74F8LzXNsh8Pbl7HXoIzEiOHB5XStFnrxkIL3sYAJGH/pGbX3SxeyfoZhY4ikEyqX3OB7Pat2yHh/63XSPThRVpD7g0gy5N2aKBz3vrHCPhe3QQTzWmKlJOcg1FE5ZtSUEHdVQbm1GD9zP6KZDfbekh9+xU0EFczW9JF/we61LTvMF1KoxaBpL46O/J6ROEOQsb03hLEMadBKxZ+XaqAHiQWKu6G5YH2opNTGKcvSyNfDInOvAygUOfzLgTCWp7JOU09hWBKW1ya2yJNJMZ6q9giEAlqS/qqYy4gAqZKjt7nF0siOb3Vz6zEaXdhCcqrfnNN6n/kFXWz24yAucW+/EHt+hsygEVUZQ=="`
	ID       string         `json:"id" example:"d423c645-fd50-4841-8138-192ee8e23dde"`
	ProofKey string         `json:"proof_key" example:"L0QIWxtHeWeUQhmfWqB2n+MZXFqEYctltilM0j69tBd1drUoUSz/vUkaPadQAdKqtQOD43Py7/JZt5IFyx7iDdphzJEX7bqq+B6nC2DQUeISEiXwtDmJYMp20/N23DY2T7L/Z/dzbxRZDWoUhtr9fRPeJL8NHtPqU9YZw2f1tgMk2t/ZMKtBhYzO5BnF8Crmshjw6b6KA3fK+j7YrmF8fVpVFCdz5jd7cprf5RIqwVjt4w1cYZWeKvGLWeGVX3oiCB67EzXZVUCsD03evr90GDY9qGLfUaWJdBkNjByDotLY0OhrKpcZ+O0IZyZv1+YKx7ZDoPAsEJqpqw4M9bGQRg=="`
}

type PubKey struct {
	ID           string         `json:"id" example:"did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx#keys-1"`
	Type         string         `json:"type" example:"RsaVerificationKey2018"`
	Controller   sdk.AccAddress `json:"controller" swaggertype:"string" example:"did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf"`
	PublicKeyPem string         `json:"publicKeyPem" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHg\nkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScAD\nG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUc\nzhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7\nZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0\nO2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfK\nHQIDAQAB\n-----END PUBLIC KEY-----\n"`
}
type PubKeys []PubKey

type Proof struct {
	Type               string    `json:"type" example:"EcdsaSecp256k1VerificationKey2019"`
	Created            time.Time `json:"created" example:"2020-04-22T04:23:50.73112321Z"`
	ProofPurpose       string    `json:"proofPurpose" example:"authentication"`
	Controller         string    `json:"controller" example:"did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx"`
	VerificationMethod string    `json:"verificationMethod" example:"did:com:pub1addwnpepqt6lnn5v0c3rys49v5v9f4kvcchehnu7kyk8t8vce5lsxfy7e2pxwyvmf6t"`
	SignatureValue     string    `json:"signatureValue" example:"nIgRvObXlF2OIbktZcQJw0UU7zDEku8cEBq7194YOjhEvD5wBZ+TcNu9GNRZucC6OyuplHfK6uo57+3lVQbpgA=="`
}

type RequestStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type PowerUpRequestProof struct {
	SenderDid   sdk.AccAddress `json:"sender_did"`
	PairwiseDid sdk.AccAddress `json:"pairwise_did"`
	Timestamp   int64          `json:"timestamp"`
	Signature   string         `json:"signature"`
}

type DidDocumentUnsigned DidDocument

// Service represents a service type needed for DidDocument.
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type Services []Service

type GenesisState struct {
	DidDocuments    []DidDocument       `json:"did_documents"`
	PowerUpRequests []DidPowerUpRequest `json:"power_up_requests"`
}
