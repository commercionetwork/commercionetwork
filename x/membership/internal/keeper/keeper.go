package keeper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/membership/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/nft/exported"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	// NFT keeper to mint tokens
	NftKeeper nft.Keeper

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	Cdc *codec.Codec
}

// NewKeeper creates new instances of the membership module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, nftKeeper nft.Keeper) Keeper {
	return Keeper{
		StoreKey:  storeKey,
		NftKeeper: nftKeeper,
		Cdc:       cdc,
	}
}

// Utility method that allows to retrieve the id of a token representing a membership associated to the given user
func (keeper Keeper) getMembershipTokenId(user sdk.AccAddress) string {
	return "membership-" + user.String()
}

func (keeper Keeper) getMembershipUri(membershipType string, id string) string {
	return fmt.Sprintf("membership:%s:%s", membershipType, id)
}

// AssignMembership allow to mint and assign a membership of the given membershipType to the specified user
func (keeper Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string) error {

	// Check the membership type validity
	if !types.IsMembershipTypeValid(membershipType) {
		return errors.New("invalid membership type")
	}

	// Make sure the user does not yet have a membership
	membership, found := keeper.GetMembership(ctx, user)
	if found && !types.CanUpgrade(keeper.GetMembershipType(membership), membershipType) {
		return errors.New("user already has a membership")
	}

	// Build the token information
	id := keeper.getMembershipTokenId(user)
	uri := keeper.getMembershipUri(membershipType, id)

	// Build the membership token
	membershipToken := nft.NewBaseNFT(id, user, uri)

	// Mint the token
	if err := keeper.NftKeeper.MintNFT(ctx, types.NftDenom, &membershipToken); err != nil {
		return err
	}

	// Return with no error
	return nil
}

// GetMembership allows to retrieve any existent membership for the specified user.
// The second returned false (the boolean one) tells if the NFT token representing the membership was found or not
func (keeper Keeper) GetMembership(ctx sdk.Context, user sdk.AccAddress) (exported.NFT, bool) {
	foundToken, err := keeper.NftKeeper.GetNFT(ctx, types.NftDenom, keeper.getMembershipTokenId(user))

	// The token was not found
	if err != nil {
		return nil, false
	}

	return foundToken, true
}

func (keeper Keeper) GetMembershipType(membership exported.NFT) string {
	return strings.Split(membership.GetTokenURI(), ":")[1]
}
