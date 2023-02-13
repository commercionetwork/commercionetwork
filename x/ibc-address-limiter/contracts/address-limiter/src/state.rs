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
/*#[cfg(test)]
pub mod tests {
    use super::*;

    pub const RESET_TIME_DAILY: u64 = 60 * 60 * 24;
    pub const RESET_TIME_WEEKLY: u64 = 60 * 60 * 24 * 7;
    pub const RESET_TIME_MONTHLY: u64 = 60 * 60 * 24 * 30;

    #[test]
    fn flow() {
        let epoch = Timestamp::from_seconds(0);
        let mut flow = Flow::new(0_u32, 0_u32, epoch, RESET_TIME_WEEKLY);

        assert!(!flow.is_expired(epoch));
        assert!(!flow.is_expired(epoch.plus_seconds(RESET_TIME_DAILY)));
        assert!(!flow.is_expired(epoch.plus_seconds(RESET_TIME_WEEKLY)));
        assert!(flow.is_expired(epoch.plus_seconds(RESET_TIME_WEEKLY).plus_nanos(1)));

        assert_eq!(flow.balance(), (0_u32.into(), 0_u32.into()));
        flow.add_flow(FlowType::In, 5_u32.into());
        assert_eq!(flow.balance(), (5_u32.into(), 0_u32.into()));
        flow.add_flow(FlowType::Out, 2_u32.into());
        assert_eq!(flow.balance(), (3_u32.into(), 0_u32.into()));
        // Adding flow doesn't affect expiration
        assert!(!flow.is_expired(epoch.plus_seconds(RESET_TIME_DAILY)));

        flow.expire(epoch.plus_seconds(RESET_TIME_WEEKLY), RESET_TIME_WEEKLY);
        assert_eq!(flow.balance(), (0_u32.into(), 0_u32.into()));
        assert_eq!(flow.inflow, Uint256::from(0_u32));
        assert_eq!(flow.outflow, Uint256::from(0_u32));
        assert_eq!(flow.period_end, epoch.plus_seconds(RESET_TIME_WEEKLY * 2));

        // Expiration has moved
        assert!(!flow.is_expired(epoch.plus_seconds(RESET_TIME_WEEKLY).plus_nanos(1)));
        assert!(!flow.is_expired(epoch.plus_seconds(RESET_TIME_WEEKLY * 2)));
        assert!(flow.is_expired(epoch.plus_seconds(RESET_TIME_WEEKLY * 2).plus_nanos(1)));
    }
}*/
