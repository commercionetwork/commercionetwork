package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func (k Keeper) SetPosition(ctx sdk.Context, position types.Position) error {
	store := ctx.KVStore(k.storeKey)
	owner, err := sdk.AccAddressFromBech32(position.Owner)
	if err != nil {
		return err
	}
	key := makePositionKey(owner, position.ID)

	if store.Has(key) {
		return fmt.Errorf("a position with id %s already exists", position.ID)
	}

	store.Set(key, k.cdc.MustMarshal(&position))

	return nil
}

func (k Keeper) UpdatePosition(ctx sdk.Context, position types.Position) error {
	store := ctx.KVStore(k.storeKey)
	owner, err := sdk.AccAddressFromBech32(position.Owner)
	if err != nil {
		return err
	}
	key := makePositionKey(owner, position.ID)

	if !store.Has(key) {
		return fmt.Errorf("a position with id %s doesn't exists", position.ID)
	}

	store.Set(key, k.cdc.MustMarshal(&position))

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
	k.cdc.MustUnmarshal(bs, &position)
	return position, true
}

func (k Keeper) GetPositionById(ctx sdk.Context, id string) (types.Position, bool) {
	position := types.Position{}
	positions := k.GetAllPositions(ctx)
	for _, p := range positions {
		if p.ID == id {
			return *p, true
		}
	}
	return position, false
}

func (k Keeper) GetAllPositionsOwnedBy(ctx sdk.Context, owner sdk.AccAddress) []*types.Position {
	var positions []*types.Position
	i := k.newPositionsByOwnerIterator(ctx, owner)
	defer i.Close()
	for ; i.Valid(); i.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(i.Value(), &position)
		positions = append(positions, &position)
	}
	return positions
}

// NewPosition creates a new minting position for the amount deposited, credited to depositor.
func (k Keeper) NewPosition(ctx sdk.Context, depositor string, deposit sdk.Coins, id string) error {
	owner, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return err
	}
	ucccRequested := deposit.AmountOf(types.CreditsDenom)
	if ucccRequested.IsZero() {
		return fmt.Errorf("no %s requested", types.CreditsDenom)
	}

	conversionRate := k.GetConversionRate(ctx)

	uccDec := sdk.NewDecFromInt(ucccRequested)
	ucommercioAmount := uccDec.Mul(conversionRate).Ceil().TruncateInt()

	// Create ucccEmitted token
	ucccEmitted := sdk.NewCoin(types.CreditsDenom, ucccRequested)

	ucomAmount := sdk.NewCoin(types.BondDenom, ucommercioAmount)

	// Create the ETP and validate it
	position := types.NewPosition(
		owner,
		ucomAmount.Amount,
		ucccEmitted,
		id,
		ctx.BlockTime(),
		conversionRate,
	)
	createAt := ctx.BlockTime()
	position.CreatedAt = &createAt
	position.ExchangeRate = conversionRate
	position.Credits = &ucccEmitted

	if err := position.Validate(); err != nil {
		return fmt.Errorf("could not validate position message, %w", err)
	}

	// Send the deposit from the user to the commerciomint account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(ucomAmount)); err != nil {
		return fmt.Errorf("could not move collateral amount to module account, %w", err)
	}

	// Mint the tokens and send them to the user
	creditsCoins := sdk.NewCoins(ucccEmitted)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, creditsCoins); err != nil {
		return fmt.Errorf("could not mint coins, %w", err)
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, creditsCoins); err != nil {
		return fmt.Errorf("could not send minted coins to account, %w", err)
	}

	// Create position
	k.SetPosition(ctx, position)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventNewPosition,
		sdk.NewAttribute("depositor", owner.String()),
		sdk.NewAttribute("amount_deposited", ucomAmount.String()),
		sdk.NewAttribute("minted_coins", creditsCoins.String()),
		sdk.NewAttribute("position_id", position.ID),
		sdk.NewAttribute("timestamp", position.CreatedAt.String()),
	))
	logger := k.Logger(ctx)
	logger.Debug("mint successful")

	return nil
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []*types.Position {
	positions := []*types.Position{}
	iterator := k.newPositionsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pos types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &pos)
		positions = append(positions, &pos)
	}

	return positions
}

// RemoveCCC burns burnAmount to the conversion rate stored in the Position identified by id, and returns the
// resulting collateral amount to user.
func (k Keeper) RemoveCCC(ctx sdk.Context, user sdk.AccAddress, id string, burnAmount sdk.Coin) (sdk.Int, error) {
	pos, found := k.GetPosition(ctx, user, id)
	residualAmount := sdk.NewInt(0)
	if !found {
		msg := fmt.Sprintf("position for user with address %s and id %s does not exist", user, id)
		return residualAmount, errors.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	// Control if position is almost in freezing period
	freezePeriod := k.GetFreezePeriod(ctx)
	createdAt := *pos.CreatedAt // TODO CHECK FORMAT AND ERROR
	if ctx.BlockTime().Sub(createdAt) <= freezePeriod {
		return residualAmount, errors.Wrap(sdkErr.ErrInvalidRequest, "cannot burn position yet in the freeze period")
	}

	// Control if tokens requested to burn are more than initially requested
	if pos.Credits.Amount.Sub(burnAmount.Amount).IsNegative() {
		return residualAmount, errors.Wrap(sdkErr.ErrInvalidRequest, "cannot burn more tokens that those initially requested")
	}

	shouldDeletePos := pos.Credits.Amount.Sub(burnAmount.Amount).IsZero()

	// 1. burn burnAmount tokens
	// 2. decrement pos.Credits
	// 3. give user amounts of collateral back
	// 4. decrement collateral
	// 5. save or delete position

	// 1.
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, user, types.ModuleName, sdk.NewCoins(burnAmount))
	if err != nil {
		return residualAmount, errors.Wrapf(sdkErr.ErrInvalidRequest, "cannot send tokens from sender to module, %s", err.Error())
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnAmount))
	if err != nil {
		return residualAmount, errors.Wrapf(sdkErr.ErrInvalidRequest, "cannot burn coins, %s", err)
	}

	// 2.
	residualAmount = pos.Credits.Amount.Sub(burnAmount.Amount)
	*pos.Credits = sdk.NewCoin(
		pos.Credits.Denom,
		residualAmount,
	)

	// 3.
	burnAmountDec := sdk.NewDecFromInt(burnAmount.Amount)
	collateralAmount := burnAmountDec.Mul(pos.ExchangeRate).Ceil().TruncateInt()
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(
		types.BondDenom,
		collateralAmount,
	)))
	if err != nil {
		return residualAmount, errors.Wrapf(sdkErr.ErrInvalidRequest, "cannot send collateral from module to sender, %s", err.Error())
	}

	// 4.
	posCollateral := sdk.NewInt(pos.Collateral)
	diffCollateral := posCollateral.Sub(collateralAmount) // TODO CONVERT INT64 TO COIN
	pos.Collateral = diffCollateral.Int64()               // TODO Should panic

	// TODO Review events
	defer func(deleted bool, ctx sdk.Context) {
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventBurnCCC,
			sdk.NewAttribute("position_id", pos.ID),
			sdk.NewAttribute("sender", user.String()),
			sdk.NewAttribute("amount", burnAmount.String()),
			sdk.NewAttribute("position_deleted", strconv.FormatBool(shouldDeletePos))))
		logger := k.Logger(ctx)
		logger.Debug("burn successful")

	}(shouldDeletePos, ctx)

	// 5.
	if shouldDeletePos {
		residualAmount = sdk.NewInt(0)
		k.deletePosition(ctx, pos)
		return residualAmount, nil
	}

	k.UpdatePosition(ctx, pos)

	return residualAmount, nil
}

func (k Keeper) newPositionsByOwnerIterator(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	prefix := append([]byte(types.EtpStorePrefix), owner...)
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
}

func getEtpsByOwnerStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.EtpStorePrefix), user...)
}

func makePositionKey(address sdk.AccAddress, id string) []byte {
	base := append([]byte(types.EtpStorePrefix), address...)
	return append(base, []byte(id)...)
}

func (k Keeper) newPositionsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.EtpStorePrefix))
}

func (k Keeper) deletePosition(ctx sdk.Context, pos types.Position) {
	store := ctx.KVStore(k.storeKey)
	owner, _ := sdk.AccAddressFromBech32(pos.Owner)
	key := makePositionKey(owner, pos.ID)
	if bs := store.Get(key); bs == nil {
		panic(fmt.Sprintf("no pos stored at key %s", key))
	}
	store.Delete(key)
}
