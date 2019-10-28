package keeper

import (
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/nft/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var membershipCosts = map[string]int64{
	types.MembershipTypeBronze: 25,
	types.MembershipTypeSilver: 250,
	types.MembershipTypeGold:   2500,
	types.MembershipTypeBlack:  25000,
}

type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	nftKeeper    nft.Keeper
	supplyKeeper supply.Keeper
}

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, nftK nft.Keeper, supplyKeeper supply.Keeper) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		nftKeeper:    nftK,
		supplyKeeper: supplyKeeper,
	}
}

// getMembershipTokenID allows to retrieve the id of a token representing a membership associated to the given user
func (k Keeper) getMembershipTokenID(user sdk.AccAddress) string {
	return "membership-" + user.String()
}

// getMembershipURI allows to returns the URI of the NFT that represents a membership of the
// given membershipType and having the given id
func (k Keeper) getMembershipURI(membershipType string, id string) string {
	return fmt.Sprintf("membership:%s:%s", membershipType, id)
}

// BuyMembership allow to mint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
func (k Keeper) BuyMembership(ctx sdk.Context, buyer sdk.AccAddress, membershipType string) sdk.Error {
	// Get the tokens from the buyer account
	membershipPrice := membershipCosts[membershipType] * 1000000 // Always multiply by one million
	membershipCost := sdk.NewCoins(sdk.NewInt64Coin(k.GetStableCreditsDenom(ctx), membershipPrice))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, membershipCost); err != nil {
		return err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, membershipCost); err != nil {
		return err
	}

	// Assign the membership
	if _, err := k.AssignMembership(ctx, buyer, membershipType); err != nil {
		return err
	}

	return nil
}

// AssignMembership allow to mint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
// Returns the URI of the new minted token represented the assigned membership, or an error if something goes w
func (k Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string) (string, sdk.Error) {
	// Check the membership type validity
	if !types.IsMembershipTypeValid(membershipType) {
		return "", sdk.ErrUnknownRequest("Invalid membership type")
	}

	// Find any existing membership
	if _, err := k.RemoveMembership(ctx, user); err != nil {
		return "", sdk.ErrUnknownRequest(err.Error())
	}

	// Build the token information
	id := k.getMembershipTokenID(user)
	uri := k.getMembershipURI(membershipType, id)

	// Build the membership token
	membershipToken := nft.NewBaseNFT(id, user, uri)

	// Mint the token
	if err := k.nftKeeper.MintNFT(ctx, types.NftDenom, &membershipToken); err != nil {
		return "", err
	}

	// Return with no error
	return membershipToken.TokenURI, nil
}

// GetMembership allows to retrieve any existent membership for the specified user.
// The second returned false (the boolean one) tells if the NFT token representing the membership was found or not
func (k Keeper) GetMembership(ctx sdk.Context, user sdk.AccAddress) (exported.NFT, bool) {
	foundToken, err := k.nftKeeper.GetNFT(ctx, types.NftDenom, k.getMembershipTokenID(user))

	// The token was not found
	if err != nil {
		return nil, false
	}

	return foundToken, true
}

// RemoveMembership allows to remove any existing membership associated with the given user.
func (k Keeper) RemoveMembership(ctx sdk.Context, user sdk.AccAddress) (bool, sdk.Error) {
	id := k.getMembershipTokenID(user)

	if found, _ := k.nftKeeper.GetNFT(ctx, types.NftDenom, id); found == nil {
		// The token was not found, so it's trivial to delete it: simply do nothing
		return true, nil
	}

	if err := k.nftKeeper.DeleteNFT(ctx, types.NftDenom, k.getMembershipTokenID(user)); err != nil {
		// The token was found, but an error was raised during the deletion. Return the error
		return false, err
	}

	// The token was found and deleted
	return true, nil
}

// GetMembershipType returns the type of the membership represented by the given NFT token
func (k Keeper) GetMembershipType(membership exported.NFT) string {
	return strings.Split(membership.GetTokenURI(), ":")[1]
}

// Get GetMembershipsSet returns the list of all the memberships
// that have been minted and are currently stored inside the store
func (k Keeper) GetMembershipsSet(ctx sdk.Context) []types.Membership {
	collection, found := k.nftKeeper.GetCollection(ctx, types.NftDenom)
	if !found {
		return nil
	}

	memberships := make([]types.Membership, len(collection.NFTs))
	for index, membershipNft := range collection.NFTs {
		memberships[index] = types.Membership{
			Owner:          membershipNft.GetOwner(),
			MembershipType: k.GetMembershipType(membershipNft),
		}
	}

	return memberships
}

// GetStableCreditsDenom returns the denom that must be used when referring to stable credits
// that can be used to purchase a membership
func (k Keeper) GetStableCreditsDenom(ctx sdk.Context) (denom string) {
	store := ctx.KVStore(k.storeKey)
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.StableCreditsStoreKey)), &denom)
	return denom
}

// SetStableCreditsDenom allows to set the denom of the coins that must be used as stable credits
// when purchasing a membership.
func (k Keeper) SetStableCreditsDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.StableCreditsStoreKey), k.cdc.MustMarshalBinaryBare(&denom))
}
