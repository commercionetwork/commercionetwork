package keeper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
)

const (
	eventNewPosition       = "new_position"
	eventBurnCCC           = "burned_ccc"
	eventSetConversionRate = "new_conversion_rate"
)

type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	supplyKeeper supply.Keeper
	govKeeper    government.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper supply.Keeper, govKeeper government.Keeper) Keeper {
	// ensure commerciomint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		supplyKeeper: supplyKeeper,
		govKeeper:    govKeeper,
	}
}

// --------------
// --- Positions
// --------------

func (k Keeper) SetPosition(ctx sdk.Context, position types.Position) error {
	store := ctx.KVStore(k.storeKey)
	key := makePositionKey(position.Owner, position.ID)

	if store.Has(key) {
		return fmt.Errorf("a position with id %s already exists", position.ID)
	}

	store.Set(key, k.cdc.MustMarshalBinaryBare(position))

	return nil
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

// NewPosition creates a new minting position for the amount deposited, credited to depositor.
func (k Keeper) NewPosition(ctx sdk.Context, depositor sdk.AccAddress, deposit sdk.Coins, id string) error {
	ucccRequested := deposit.AmountOf("uccc")
	if ucccRequested.IsZero() {
		return errors.New("no uccc requested")
	}

	conversionRate := k.GetConversionRate(ctx)

	uccDec := sdk.NewDecFromInt(ucccRequested)
	ucommercioAmount := uccDec.Mul(conversionRate).Ceil().TruncateInt()

	// Create ucccEmitted token
	ucccEmitted := sdk.NewCoin(types.CreditsDenom, ucccRequested)

	ucomAmount := sdk.NewCoin("ucommercio", ucommercioAmount)

	// Create the ETP and validate it
	position := types.NewPosition(
		depositor,
		ucomAmount.Amount,
		ucccEmitted,
		id,
		ctx.BlockTime(),
		conversionRate,
	)
	if err := position.Validate(); err != nil {
		return fmt.Errorf("could not validate position message, %w", err)
	}

	// Send the deposit from the user to the commerciomint account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(ucomAmount)); err != nil {
		return fmt.Errorf("could not move collateral amount to module account, %w", err)
	}

	// Mint the tokens and send them to the user
	creditsCoins := sdk.NewCoins(ucccEmitted)
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, creditsCoins); err != nil {
		return fmt.Errorf("could not mint coins, %w", err)
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, creditsCoins); err != nil {
		return fmt.Errorf("could not send minted coins to account, %w", err)
	}

	// Create position
	k.SetPosition(ctx, position)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventNewPosition,
		sdk.NewAttribute("depositor", depositor.String()),
		sdk.NewAttribute("amount_deposited", ucomAmount.String()),
		sdk.NewAttribute("minted_coins", creditsCoins.String()),
		sdk.NewAttribute("position_id", position.ID),
		sdk.NewAttribute("timestamp", position.CreatedAt.String()),
	))

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

// BurnCCC burns burnAmount to the conversion rate stored in the Position identified by id, and returns the
// resulting collateral amount to user.
func (k Keeper) BurnCCC(ctx sdk.Context, user sdk.AccAddress, id string, burnAmount sdk.Coin) error {
	pos, found := k.GetPosition(ctx, user, id)
	if !found {
		msg := fmt.Sprintf("position for user with address %s and id %s does not exist", user, id)
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	if pos.Credits.Amount.Sub(burnAmount.Amount).IsNegative() {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "cannot burn more tokens that those initially requested")
	}

	shouldDeletePos := pos.Credits.Amount.Sub(burnAmount.Amount).IsZero()

	// 1. burn burnAmount tokens
	// 2. decrement pos.Credits
	// 3. give user amounts of collateral back
	// 4. decrement collateral
	// 5. save or delete position

	// 1.
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, user, types.ModuleName, sdk.NewCoins(burnAmount))
	if err != nil {
		return sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "cannot send tokens from sender to module, %s", err.Error())
	}

	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnAmount))
	if err != nil {
		return sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "cannot burn coins, %s", err)
	}

	// 2.
	pos.Credits = sdk.NewCoin(
		pos.Credits.Denom,
		pos.Credits.Amount.Sub(burnAmount.Amount),
	)

	// 3.
	burnAmountDec := sdk.NewDecFromInt(burnAmount.Amount)
	collateralAmount := burnAmountDec.Mul(pos.ExchangeRate).Ceil().TruncateInt()
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(
		"ucommercio",
		collateralAmount,
	)))
	if err != nil {
		return sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "cannot send collateral from module to sender, %s", err.Error())
	}

	// 4.
	pos.Collateral = pos.Collateral.Sub(collateralAmount)

	defer func(deleted bool, ctx sdk.Context) {
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventBurnCCC,
			sdk.NewAttribute("position_id", pos.ID),
			sdk.NewAttribute("sender", user.String()),
			sdk.NewAttribute("amount", burnAmount.String()),
			sdk.NewAttribute("position_deleted", strconv.FormatBool(shouldDeletePos))))
	}(shouldDeletePos, ctx)

	// 5.
	if shouldDeletePos {
		k.deletePosition(ctx, pos)
		return nil
	}

	k.SetPosition(ctx, pos)

	return nil
}

// GetConversionRate retrieve the conversion rate.
func (k Keeper) GetConversionRate(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	var rate sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.CollateralRateKey)), &rate)
	return rate
}

// SetConversionRate store the conversion rate.
func (k Keeper) SetConversionRate(ctx sdk.Context, rate sdk.Dec) error {
	if err := types.ValidateConversionRate(rate); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.CollateralRateKey), k.cdc.MustMarshalBinaryBare(rate))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSetConversionRate,
		sdk.NewAttribute("rate", rate.String()),
	))

	return nil
}

func (k Keeper) newPositionsByOwnerIterator(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	prefix := append([]byte(types.EtpStorePrefix), owner...)
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
}

func (k Keeper) newPositionsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.EtpStorePrefix))
}

func makePositionKey(address sdk.AccAddress, id string) []byte {
	base := append([]byte(types.EtpStorePrefix), address...)
	return append(base, []byte(id)...)
}

func (k Keeper) deletePosition(ctx sdk.Context, pos types.Position) {
	store := ctx.KVStore(k.storeKey)
	key := makePositionKey(pos.Owner, pos.ID)
	if bs := store.Get(key); bs == nil {
		panic(fmt.Sprintf("no pos stored at key %s", key))
	}
	store.Delete(key)
}
