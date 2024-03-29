use cosmwasm_schema::write_api;

use address_limiter::msg::{ExecuteMsg, InstantiateMsg, MigrateMsg, QueryMsg, SudoMsg};

fn main() {
    write_api! {
        instantiate: InstantiateMsg,
        query: QueryMsg,
        execute: ExecuteMsg,
        sudo: SudoMsg,
        migrate: MigrateMsg,
    }
}
