package keeper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
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

// AddTrustedMinter allows to add the given minter as a trusted address that can sign the
// minting of new memberships tokens
func (keeper Keeper) AddTrustedMinter(ctx sdk.Context, minter sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	// Save the minter
	key := []byte(types.TrustworthyMinterPrefix + minter.String())
	store.Set(key, minter)
}

// GetTrustedMinters returns the list of the current addresses that are allowed to mint
// a new membership token when necessary
func (keeper Keeper) GetTrustedMinters(ctx sdk.Context) types.Minters {
	store := ctx.KVStore(keeper.StoreKey)

	var minters []sdk.AccAddress

	// Iterate over all the keys having the minter prefix
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.TrustworthyMinterPrefix))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// Add each minter to the list
		minters = append(minters, iterator.Value())
	}

	return minters
}

// Utility method that allows to retrieve the id of a token representing a membership associated to the given user
func (keeper Keeper) getMembershipTokenId(user sdk.AccAddress) string {
	return "membership-" + user.String()
}

func (keeper Keeper) getMembershipUri(membershipType string, id string) string {
	return fmt.Sprintf("membership:%s:%s", membershipType, id)
}

// AssignMembership allow to mint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
// Returns the URI of the new minted token represented the assigned membership, or an error if something goes wrong
func (keeper Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string) (string, error) {
	// Check the membership type validity
	if !types.IsMembershipTypeValid(membershipType) {
		return "", errors.New("invalid membership type")
	}

	// Find any existing membership
	if _, err := keeper.RemoveMembership(ctx, user); err != nil {
		return "", err
	}

	// Build the token information
	id := keeper.getMembershipTokenId(user)
	uri := keeper.getMembershipUri(membershipType, id)

	// Build the membership token
	membershipToken := nft.NewBaseNFT(id, user, uri)

	// Mint the token
	if err := keeper.NftKeeper.MintNFT(ctx, types.NftDenom, &membershipToken); err != nil {
		return "", err
	}

	// Return with no error
	return membershipToken.TokenURI, nil
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

// RemoveMembership allows to remove any existing membership associated with the given user.
func (keeper Keeper) RemoveMembership(ctx sdk.Context, user sdk.AccAddress) (bool, error) {
	id := keeper.getMembershipTokenId(user)

	if found, _ := keeper.NftKeeper.GetNFT(ctx, types.NftDenom, id); found == nil {
		// The token was not found, so it's trivial to delete it: simply do nothing
		return true, nil
	}

	if err := keeper.NftKeeper.DeleteNFT(ctx, types.NftDenom, keeper.getMembershipTokenId(user)); err != nil {
		// The token was found, but an error was raised during the deletion. Return the error
		return false, err
	}

	// The token was found and deleted
	return true, nil
}

// GetMembershipType returns the type of the membership represented by the given NFT token
func (keeper Keeper) GetMembershipType(membership exported.NFT) string {
	return strings.Split(membership.GetTokenURI(), ":")[1]
}

// ----------------------
// --- Genesis utils
// ----------------------

// Get GetMembershipsSet returns the list of all the memberships
// that have been minted and are currently stored inside the store
func (keeper Keeper) GetMembershipsSet(ctx sdk.Context) []types.Membership {
	var memberships []types.Membership

	collection, found := keeper.NftKeeper.GetCollection(ctx, types.NftDenom)
	if !found {
		return memberships
	}

	for _, membershipNft := range collection.NFTs {
		membership := types.Membership{
			Owner:          membershipNft.GetOwner(),
			MembershipType: keeper.GetMembershipType(membershipNft),
		}
		memberships = append(memberships, membership)
	}

	return memberships
}
