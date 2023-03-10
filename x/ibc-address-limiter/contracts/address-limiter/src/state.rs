use cosmwasm_std::Addr;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cw_storage_plus::Item;

use crate::ContractError;

#[derive(Debug, Clone)]
pub enum FlowType {
    In,
    Out,
}

/// Only this address can manage the contract. This will likely be the
/// governance module, but could be set to something else if needed
pub const GOVMODULE: Item<Addr> = Item::new("gov_module");
/// Only this address can execute transfers. This will likely be the
/// IBC transfer module, but could be set to something else if needed
pub const IBCMODULE: Item<Addr> = Item::new("ibc_module");

//may use map<,>
pub const ADDRS_WHITELIST: Item<Whitelist> = Item::new("addrs_whitelist");

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Whitelist {
    pub wl: Vec<Addr>,
}

impl Whitelist {
    /// Checks if a transfer is allowed
    /// If the transfer is not allowed, it will return a Unauthorized error.
    /// Otherwise it will return ?
    pub fn allow_transfer(
        &mut self,
        direction: &FlowType,
        sender: Addr
    ) -> Result<Whitelist, ContractError> {

        match direction.clone() {
            FlowType::In => Ok(self.clone()),
            FlowType::Out => {
                if !self.wl.contains(&sender) {
                    return Err(ContractError::Unauthorized {addr:sender.into()});
                }
                Ok(self.clone())
            },
        }
    }
}