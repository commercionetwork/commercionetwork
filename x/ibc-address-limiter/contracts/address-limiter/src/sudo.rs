use cosmwasm_std::{DepsMut, Response};

use crate::{
    packet::FungibleTokenData,
    state::{FlowType, ADDRS_WHITELIST},
    ContractError,
};

pub fn process_packet(
    deps: DepsMut,
    packet_data: FungibleTokenData,
    direction: FlowType,
) -> Result<Response, ContractError> {
    // Sudo call. Only go modules should be allowed to access this
    let mut whitelist = ADDRS_WHITELIST.load(deps.storage)?;

    let not_configured = whitelist.wl.is_empty();
    if not_configured {
        // No Addresses configured for the current path. Allowing all messages.
        return Ok(Response::new()
            .add_attribute("method", "try_transfer")
            .add_attribute("whitelist", "empty"));
    }

    //let packet_data = packet.data;
    // If it fails, allow_transfer() will return
    // ContractError::Unauthorized, which we'll propagate out
   _ = whitelist.allow_transfer(&direction,  packet_data.sender.clone())?;

    let response = Response::new()
        .add_attribute("method", "try_transfer")
        .add_attribute("sender_address", packet_data.sender.into_string());

    Ok(response)
}