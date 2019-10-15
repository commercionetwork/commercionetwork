package keeper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/nft/exported"
)

type Keeper struct {
	cdc        *codec.Codec
	StoreKey   sdk.StoreKey
	BankKeeper bank.Keeper
	NftKeeper  nft.Keeper
}

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, nftK nft.Keeper, bankK bank.Keeper) Keeper {
	return Keeper{
		StoreKey:   storeKey,
		BankKeeper: bankK,
		NftKeeper:  nftK,
		cdc:        cdc,
	}
}

// -------------------------
// --- Invites
// -------------------------

func (keeper Keeper) getInviteStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.InviteStorePrefix + user.String())
}

// InviteUser allows to set a given user as being invited by the given invite sender.
func (keeper Keeper) InviteUser(ctx sdk.Context, recipient, sender sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)
	inviteKey := keeper.getInviteStoreKey(recipient)

	if store.Has(inviteKey) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("%s has already been invited", recipient.String()))
	}

	// Build the accreditation
	accreditation := types.Invite{
		Sender:   sender,
		User:     recipient,
		Rewarded: false,
	}

	// Save the accreditation
	accreditationBz := keeper.cdc.MustMarshalBinaryBare(&accreditation)
	store.Set(inviteKey, accreditationBz)
	return nil
}

// GetInvite allows to get the invitation related to a user
func (keeper Keeper) GetInvite(ctx sdk.Context, user sdk.AccAddress) (invite types.Invite, found bool) {
	store := ctx.KVStore(keeper.StoreKey)
	key := keeper.getInviteStoreKey(user)

	if store.Has(key) {
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &invite)
		return invite, true
	}

	return types.Invite{}, false
}

// GetInvites returns all the invites ever made
func (keeper Keeper) GetInvites(ctx sdk.Context) (invites []types.Invite) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.InviteStorePrefix))

	for ; iterator.Valid(); iterator.Next() {
		var invite types.Invite
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &invite)
		invites = append(invites, invite)
	}

	return
}

func (keeper Keeper) SaveInvite(ctx sdk.Context, invite types.Invite) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(keeper.getInviteStoreKey(invite.User), keeper.cdc.MustMarshalBinaryBare(&invite))
}

// ---------------------
// --- Verifications
// ---------------------

func (keeper Keeper) getCredentialsStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.CredentialsStorePrefix + user.String())
}

func (keeper Keeper) SaveCredential(ctx sdk.Context, credential types.Credential) {
	credentials := keeper.GetUserCredentials(ctx, credential.User)
	if credentials, edited := credentials.AppendIfMissing(credential); edited {
		store := ctx.KVStore(keeper.StoreKey)
		store.Set(keeper.getCredentialsStoreKey(credential.User), keeper.cdc.MustMarshalBinaryBare(&credentials))
	}
}

func (keeper Keeper) GetUserCredentials(ctx sdk.Context, user sdk.AccAddress) types.Credentials {
	store := ctx.KVStore(keeper.StoreKey)

	var credentials types.Credentials
	bz := store.Get(keeper.getCredentialsStoreKey(user))
	keeper.cdc.MustUnmarshalBinaryBare(bz, &credentials)

	return credentials
}

func (keeper Keeper) GetCredentials(ctx sdk.Context) (credentials types.Credentials) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CredentialsStorePrefix))

	credentials = types.Credentials{}
	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &credential)
		credentials, _ = credentials.AppendIfMissing(credential)
	}

	return credentials
}

// ---------------------
// --- Reward pool
// ---------------------

// DepositIntoPool allows anyone to deposit into the liquidity pool that
// will be used when giving out rewards for accreditations.
func (keeper Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) error {
	if amount.IsAnyNegative() {
		return errors.New("amount cannot be negative")
	}

	store := ctx.KVStore(keeper.StoreKey)

	// Remove the coins from the user wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, depositor, amount)
	if err != nil {
		return err
	}

	// Add the amount to the pool
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolStoreKey)), &pool)
	pool = pool.Add(amount)
	store.Set([]byte(types.LiquidityPoolStoreKey), keeper.cdc.MustMarshalBinaryBare(&pool))

	return nil
}

// SetPoolFunds allows to set the current pool funds amount
func (keeper Keeper) SetPoolFunds(ctx sdk.Context, pool sdk.Coins) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), keeper.cdc.MustMarshalBinaryBare(&pool))
}

// GetPoolFunds return the current pool funds for the given context
func (keeper Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(keeper.StoreKey)
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolStoreKey)), &pool)
	return pool
}

// -----------------------------------
// --- Trusted Service Providers
// -----------------------------------

// AddTrustedServiceProvider allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (keeper Keeper) AddTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	signers := keeper.GetTrustedServiceProviders(ctx)
	if signers, success := signers.AppendIfMissing(tsp); success {
		newSignersBz := keeper.cdc.MustMarshalBinaryBare(&signers)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}
}

// GetTrustedServiceProviders returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (keeper Keeper) GetTrustedServiceProviders(ctx sdk.Context) (signers ctypes.Addresses) {
	store := ctx.KVStore(keeper.StoreKey)

	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	keeper.cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	return
}

// IsTrustedServiceProvider tells if the given signer is a trusted one or not
func (keeper Keeper) IsTrustedServiceProvider(ctx sdk.Context, signer sdk.AccAddress) bool {
	signers := keeper.GetTrustedServiceProviders(ctx)
	return signers.Contains(signer)
}

// ----------------------
// --- Memberships
// ----------------------

func (keeper Keeper) GetStableCreditsDenom(ctx sdk.Context) (denom string) {
	store := ctx.KVStore(keeper.StoreKey)
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.StableCreditsStoreKey)), &denom)
	return denom
}

func (keeper Keeper) SetStableCreditsDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.StableCreditsStoreKey), keeper.cdc.MustMarshalBinaryBare(&denom))
}

// Utility method that allows to retrieve the id of a token representing a membership associated to the given user
func (keeper Keeper) getMembershipTokenId(user sdk.AccAddress) string {
	return "membership-" + user.String()
}

func (keeper Keeper) getMembershipUri(membershipType string, id string) string {
	return fmt.Sprintf("membership:%s:%s", membershipType, id)
}

var membershipCosts = map[string]int64{
	types.MembershipTypeBronze: 25,
	types.MembershipTypeSilver: 250,
	types.MembershipTypeGold:   2500,
	types.MembershipTypeBlack:  25000,
}

// BuyMembership allow to mint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
func (keeper Keeper) BuyMembership(ctx sdk.Context, buyer sdk.AccAddress, membershipType string) sdk.Error {
	// Get the tokens from the buyer account
	membershipPrice := membershipCosts[membershipType] * 1000000 // Always multiply by one million
	membershipCost := sdk.NewCoins(sdk.NewInt64Coin(keeper.GetStableCreditsDenom(ctx), membershipPrice))
	if _, err := keeper.BankKeeper.SubtractCoins(ctx, buyer, membershipCost); err != nil {
		return err
	}

	if _, err := keeper.AssignMembership(ctx, buyer, membershipType); err != nil {
		return err
	}

	return nil
}

// AssignMembership allow to mint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
// Returns the URI of the new minted token represented the assigned membership, or an error if something goes w
func (keeper Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string) (string, sdk.Error) {
	// Check the membership type validity
	if !types.IsMembershipTypeValid(membershipType) {
		return "", sdk.ErrUnknownRequest("Invalid membership type")
	}

	// Find any existing membership
	if _, err := keeper.RemoveMembership(ctx, user); err != nil {
		return "", sdk.ErrUnknownRequest(err.Error())
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
func (keeper Keeper) RemoveMembership(ctx sdk.Context, user sdk.AccAddress) (bool, sdk.Error) {
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
