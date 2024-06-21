package ibc_address_limit

import (
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
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"
)

var (
	_ porttypes.Middleware         = &IBCModule{}
	_ porttypes.ICS4Wrapper        = &ICS4Wrapper{}
	_ porttypes.PacketDataUnmarshaler = (*ICS4Wrapper)(nil)
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

	var fullPacket channeltypes.Packet
	if unmarshaler, ok := i.channel.(porttypes.PacketDataUnmarshaler); ok {
		// Use the PacketDataUnmarshaler interface to unmarshal the data
		packetData, err := unmarshaler.UnmarshalPacketData(data)
		if err != nil {
			return 0, errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot unmarshal packet data")
		}

		fullPacket, ok = packetData.(channeltypes.Packet)
		if !ok {
			return 0, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet data type")
		}
	} else {
		// Fall back to manual unmarshalling if the interface is not implemented
		if err := i.codec.Unmarshal(data, &fullPacket); err != nil {
			return 0, errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot unmarshal packet data")
		}
	}

	err := CheckSenderAuth(ctx, i.ContractKeeper, "send_packet", contract, fullPacket)
	if err != nil {
		return 0, err
	}

	return i.channel.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// UnmarshalPacketData implements types.PacketDataUnmarshaler.
func (i *ICS4Wrapper) UnmarshalPacketData(data []byte) (interface{}, error) {
	var packet channeltypes.Packet
	err := i.codec.Unmarshal(data, &packet)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot unmarshal packet data")
	}
	return packet, nil
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
