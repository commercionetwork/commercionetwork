package v3_0_0

import (
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

const (
	ModuleName = "did"
)

type DidDocuments []*types.DidDocument

func (didDocuments DidDocuments) AppendIfMissingID(ins *types.DidDocument) {
	/*for _, present := range didDocuments {
		if present.ID == ins.ID {
			return didDocuments
		}
	}*/
	didDocuments = append(didDocuments, ins)
}
