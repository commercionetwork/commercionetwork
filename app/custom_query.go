package app

import (
	"encoding/json"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	CommercioKycKeeper  "github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type getMembershipQuery struct {
	Account string		`json:"account,omitempty"`
}

type commercioCustomQuery struct {
	GetMembership 		*getMembershipQuery 	`json:"get_membership,omitempty"`
	//GetAllMemberships   *struct{}   			`json:"get_all_memberships,omitempty"`
}


// MembershipData mirrors the Rust MembershipData struct
type MembershipData struct {
	Owner          string    `json:"owner"`
	TSPAddress     string    `json:"tsp_address"`
	MembershipType string    `json:"membership_type"`
	ExpiryAt       time.Time `json:"expiry_at"`
}

type customQueryResponse struct {
	Membership MembershipData `json:"membership"`
}

func commercioPlugins(commerciokyc CommercioKycKeeper.Keeper) *wasm.QueryPlugins {
	return &wasm.QueryPlugins{
		Custom: CommerciokycQuerier(commerciokyc),
	}
}
func CommerciokycQuerier (commerciokyc CommercioKycKeeper.Keeper) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom commercioCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if custom.GetMembership != nil {
			accountAddr, err := sdk.AccAddressFromBech32(custom.GetMembership.Account)
			if err != nil {
				return nil, err
			}
			membership, err := commerciokyc.GetMembership(ctx, accountAddr)
			if err != nil {
				return nil, sdkerrors.Wrap(types.ErrInvalidMsg, err.Error())
			}
			
			return json.Marshal(customQueryResponse{ Membership: MembershipData{
							Owner: membership.Owner,
							TSPAddress: membership.TspAddress,
							MembershipType: membership.MembershipType,
							ExpiryAt: *membership.ExpiryAt,
							}})
		}

		// if custom.GetAllMemberships != nil {
		// 	memberships := commerciokyc.GetMemberships(ctx)
		// 	return json.Marshal(memberships)
		// }
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

// func performCustomQuery(_ sdk.Context, request json.RawMessage) ([]byte, error) {
// 	var custom commercioCustomQuery
// 	err := json.Unmarshal(request, &custom)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}
// 	if custom.GetMembership != nil {
// 		//msg := "Response for get membership query"
// 		return json.Marshal(customQueryResponse{Membership: MembershipData{
// 			Owner: "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
// 			TSPAddress: "fake tsp address",
// 			MembershipType: "fake membership type",
// 			//ExpiryAt: time.Now(),
// 		}})
// 	}
// 	if custom.GetAllMemberships != nil {
// 		return json.Marshal(customQueryResponse{Membership: MembershipData{
// 			Owner: "fake owner",
// 			TSPAddress: "fake tsp address",
// 			MembershipType: "fake membership type",
// 			//ExpiryAt: time.Now(),
// 		}})
// 	}
// 	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
// }
