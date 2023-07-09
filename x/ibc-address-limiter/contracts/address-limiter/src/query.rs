use cosmwasm_std::{to_binary, Binary, Deps, StdResult};

use crate::state::ADDRS_WHITELIST;

pub fn get_whitelist(
    deps: Deps,
) -> StdResult<Binary> {
    to_binary(&ADDRS_WHITELIST.load(deps.storage)?)
}