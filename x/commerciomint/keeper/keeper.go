package keeper

import (
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"
	uuid "github.com/satori/go.uuid"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	creditrisk "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	pricefeed "github.com/commercionetwork/commercionetwork/x/pricefeed/keeper"
)

type Keeper struct {
	cdc             *codec.Codec
	storeKey        sdk.StoreKey
	priceFeedKeeper pricefeed.Keeper
	supplyKeeper    supply.Keeper
	govKeeper       government.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper supply.Keeper, pk pricefeed.Keeper, govKeeper government.Keeper) Keeper {
	// ensure commerciomint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        key,
		priceFeedKeeper: pk,
		supplyKeeper:    supplyKeeper,
		govKeeper:       govKeeper,
	}
}

// --------------
// --- Positions
// --------------

func (k Keeper) SetPosition(ctx sdk.Context, position types.Position) {
	store := ctx.KVStore(k.storeKey)
	key := makePositionKey(position.Owner, position.ID)
	if bs := store.Get(key); bs != nil {
		panic(fmt.Errorf("cannot overwrite position at key %s", key))
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(position))
}

func (k Keeper) GetPosition(ctx sdk.Context, owner sdk.AccAddress, id string) (types.Position, bool) {
	position := types.Position{}
	key := makePositionKey(owner, id)
	store := ctx.KVStore(k.storeKey)
	bs := store.Get(key)
	if bs == nil {
		return position, false
	}
	k.cdc.MustUnmarshalBinaryBare(bs, &position)
	return position, true
}

func (k Keeper) GetAllPositionsOwnedBy(ctx sdk.Context, owner sdk.AccAddress) []types.Position {
	var positions []types.Position
	i := k.newPositionsByOwnerIterator(ctx, owner)
	defer i.Close()
	for ; i.Valid(); i.Next() {
		var position types.Position
		k.cdc.MustUnmarshalBinaryBare(i.Value(), &position)
		positions = append(positions, position)
	}
	return positions
}

// NewPosition subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (k Keeper) NewPosition(ctx sdk.Context, depositor sdk.AccAddress, deposit sdk.Coins) error {
	ucomDeposit := deposit.AmountOf("ucommercio")
	if ucomDeposit.IsZero() {
		return errors.New("no ucommercio deposited")
	}

	conversionRate := k.GetConversionRate(ctx)
	creditsAmount := ucomDeposit.Quo(conversionRate)

	// Create credits token
	credits := sdk.NewCoin(types.CreditsDenom, creditsAmount)

	id := uuid.NewV4()

	// Create the CDP and validate it
	position := types.NewPosition(
		depositor,
		deposit,
		credits,
		id.String(),
		time.Now(),
		conversionRate,
	)
	if err := position.Validate(); err != nil {
		return fmt.Errorf("could not validate position message, %w", err)
	}

	// Send the deposit from the user to the commerciomint account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, deposit); err != nil {
		return fmt.Errorf("could not move collateral amount to module account, %w", err)
	}

	// Mint the tokens and send them to the user
	creditsCoins := sdk.NewCoins(credits)
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, creditsCoins); err != nil {
		return fmt.Errorf("could not mint coins, %w", err)
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, creditsCoins); err != nil {
		return fmt.Errorf("could not send minted coins to account, %w", err)
	}

	// Create position
	k.SetPosition(ctx, position)

	return nil
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []types.Position {
	var positions []types.Position
	iterator := k.newPositionsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pos types.Position
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &pos)
		positions = append(positions, pos)
	}

	return positions
}

// CloseCdp subtract the Position's liquidity amount (commercio cash credits) from user's wallet, after that sends the
// deposited amount back to it. If these two operations ends without errors, the Position get closed.
// Errors occurs if:k.GetCdpsByOwner(ctx, testCdpOwner)
// - cdp doesnt exist
// - subtracting or adding fund to account don't end well
// TODO: this thing should burn tokens, too.
func (k Keeper) CloseCdp(ctx sdk.Context, user sdk.AccAddress, id string) error {
	pos, found := k.GetPosition(ctx, user, id)
	if !found {
		msg := fmt.Sprintf("position for user with address %s and id %s does not exist", user, id)
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	// Send the coins from the user to the module and then burn them
	creditsCoins := sdk.NewCoins(pos.Credits)
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, pos.Owner, types.ModuleName, creditsCoins); err != nil {
		return err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, creditsCoins); err != nil {
		return err
	}

	// Get the user the deposited amount
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, pos.Owner, pos.Deposit); err != nil {
		return err
	}

	// Delete the CDP
	k.deletePosition(ctx, pos)

	return nil
}

// GetConversionRate retrieve the conversion rate.
func (k Keeper) GetConversionRate(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	var rate sdk.Int
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.CollateralRateKey)), &rate)
	return rate
}

// SetConversionRate store the conversion rate.
func (k Keeper) SetConversionRate(ctx sdk.Context, rate sdk.Int) error {
	if err := types.ValidateConversionRate(rate); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.CollateralRateKey), k.cdc.MustMarshalBinaryBare(rate))
	return nil
}

func (k Keeper) newPositionsByOwnerIterator(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	prefix := []byte(fmt.Sprintf("%s%s:", types.CdpStorePrefix, owner.String()))
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
}

func (k Keeper) newPositionsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.CdpStorePrefix))
}

func (k Keeper) calculateFiatValue(ctx sdk.Context, deposits sdk.Coins) (sdk.Dec, error) {
	fiatValue := sdk.ZeroDec()
	for _, deposit := range deposits {
		assetPrice, found := k.priceFeedKeeper.GetCurrentPrice(ctx, deposit.Denom)
		if !found {
			return fiatValue, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("no current price for given denom: %s", deposit.Denom))
		}
		fiatValue = fiatValue.Add(deposit.Amount.ToDec().Mul(assetPrice.Value))
	}
	return fiatValue, nil
}

func (k Keeper) liquidate(ctx sdk.Context, pos types.Position) error {
	// Send the coins from the user to the module and then burn them
	if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, creditrisk.ModuleName, pos.Deposit); err != nil {
		return err
	}
	// Delete the CDP
	k.deletePosition(ctx, pos)
	return nil
}

func makePositionKey(address sdk.AccAddress, id string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", types.CdpStorePrefix, id, address.String()))
}

func (k Keeper) deletePosition(ctx sdk.Context, pos types.Position) {
	store := ctx.KVStore(k.storeKey)
	key := makePositionKey(pos.Owner, pos.ID)
	if bs := store.Get(key); bs == nil {
		panic(fmt.Sprintf("no pos stored at key %s", key))
	}
	store.Delete(key)
}
