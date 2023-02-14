use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug, PartialEq)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized! The sender {addr} has not authorization")]
    Unauthorized {
        addr: String
    },

    #[error("Empy whitelist")]
    EmptyWhitelist {},

    #[error("Whitelist not found")]
    WhitelistNotFound {},

    #[error("Empy list of addresses")]
    EmptyList {},
}
