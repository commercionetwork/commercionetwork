package v3_0_0

import (
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

const (
	ModuleName = "did"
)

type DidDocuments []*types.DidDocumentNew

func (didDocuments DidDocuments) AppendIfMissingID(i *types.DidDocumentNew) DidDocuments {
	for _, ele := range didDocuments {
		if ele.ID == i.ID {
			return didDocuments
		}
	}
	return append(didDocuments, i)
}
