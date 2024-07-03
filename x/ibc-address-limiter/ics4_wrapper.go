package ibc_address_limit

import (
	"encoding/json"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"
)

var (
	_ porttypes.Middleware         = &IBCModule{}
	_ porttypes.ICS4Wrapper        = &ICS4Wrapper{}
)

type ICS4Wrapper struct {
	channel        porttypes.ICS4Wrapper
	accountKeeper  *authkeeper.AccountKeeper
	bankKeeper     *bankkeeper.BaseKeeper
	ContractKeeper *wasmkeeper.PermissionedKeeper
	paramSpace     paramtypes.Subspace
	codec          codec.Codec
}

func NewICS4Middleware(
	channel porttypes.ICS4Wrapper,
	accountKeeper *authkeeper.AccountKeeper, contractKeeper *wasmkeeper.PermissionedKeeper,
	bankKeeper *bankkeeper.BaseKeeper, paramSpace paramtypes.Subspace, codec codec.Codec,
) ICS4Wrapper {
	return ICS4Wrapper{
		channel:        channel,
		accountKeeper:  accountKeeper,
		ContractKeeper: contractKeeper,
		bankKeeper:     bankKeeper,
		paramSpace:     paramSpace,
		codec:          codec,
	}
}

func (i *ICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {

	contract := i.GetParams(ctx)
	if contract == "" {
		// The contract has not been configured. Continue as usual
		return i.channel.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	var packetdata transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetdata); err != nil {
		return 0, errors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot unmarshal packet data: %v", data)
	}

	// We need the full packet so the contract can process it. If it can't be cast to a channeltypes.Packet, this
	// should fail. The only reason that would happen is if another middleware is modifying the packet, though. In
	// that case we can modify the middleware order or change this cast so we have all the data we need.
	// UNFORKINGTODO OQ: The full packet data is not available here. Specifically, the sequence, destPort and destChannel are not available.
	// This is silly as it means we cannot filter packets based on destination (the sequence could be obtained by calling channel.SendPacket() before checking the address limits)
	// I think this works with the current contracts as destination is not checked for sends, but would need to double check to be 100% sure.
	// Should we modify what the contracts expect so that there's no risk of them trying to rely on the missing data? Alt. just document this
	// UNFORKINGTODO N: I am setting the sequence to 0 so it can compile, but note that this needs to be addressed.
	fullPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         sourcePort,
		SourceChannel:      sourceChannel,
		DestinationPort:    "omitted",
		DestinationChannel: "omitted",
		Data:               data,
		TimeoutTimestamp:   timeoutTimestamp,
		TimeoutHeight:      timeoutHeight,
	}

	err := CheckSenderAuth(ctx, i.ContractKeeper, "send_packet", contract, fullPacket)
	if err != nil {
		return 0, err
	}

	return i.channel.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}


func (i *ICS4Wrapper) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return i.channel.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (i *ICS4Wrapper) GetParams(ctx sdk.Context) (contract string) {
	i.paramSpace.GetIfExists(ctx, []byte("contract"), &contract)
	return contract
}

func (i *ICS4Wrapper) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	//TODO implement me
	panic("implement me")
}

func (i *ICS4Wrapper) SetParams(ctx sdk.Context, params types.Params) {
	i.paramSpace.SetParamSet(ctx, &params)
}
