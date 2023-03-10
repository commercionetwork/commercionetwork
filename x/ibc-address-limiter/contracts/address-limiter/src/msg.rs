use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Addr;

use crate::packet::Packet;

/// Initialize the contract with the address of the IBC module.
/// Only the ibc module is allowed to execute actions on this contract
#[cw_serde]
pub struct InstantiateMsg {
    pub gov_module: Addr,
    pub ibc_module: Addr,
    pub addrs_whitelist: Vec<Addr>,
}
#[cw_serde]
pub enum ExecWhitelist {
    //add one or n addresses to whitelist
    AddAddrs{
        addresses: Vec<Addr>,
    },
    //remove one or n addresses from whitelist
    RemoveAddrs{
        addresses: Vec<Addr>,
    },
    //empty whitelist
    ResetWhitelist{},
    //Remove previous list if set and creat new one.
    New{
        whitelist: Vec<Addr>,
    }
}

/// The caller (IBC module) is responsible for correctly calculating the funds
/// being sent through the channel
#[cw_serde]
pub enum ExecuteMsg {
    ManageWhitelist(ExecWhitelist),
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(crate::state::Whitelist)]
    GetWhitelist {},
}

#[cw_serde]
pub enum SudoMsg {
    SendPacket {
        packet: Packet,
    },
    RecvPacket {
        packet: Packet,
    },
}

#[cw_serde]
pub enum MigrateMsg {}
