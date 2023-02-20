use crate::msg::{ExecWhitelist};
use crate::state::{Whitelist, GOVMODULE, IBCMODULE, ADDRS_WHITELIST};
use crate::ContractError;
use cosmwasm_std::{Addr, DepsMut, Response};

pub fn try_manage_whitelist(
    deps: DepsMut, 
    sender: Addr, 
    exec_whitelist: ExecWhitelist,
) -> Result<Response, ContractError> {
    let ibc_module = IBCMODULE.load(deps.storage)?;
    let gov_module = GOVMODULE.load(deps.storage)?;
    if sender != ibc_module && sender != gov_module {
        return Err(ContractError::Unauthorized {addr: sender.into()});
    }

    match exec_whitelist {
        ExecWhitelist::AddAddrs { addresses } => {
            let mut actual_whitelist = ADDRS_WHITELIST.load(deps.storage)?;
            if addresses.is_empty() {
                return Err(ContractError::EmptyList {  });
            }

            addresses.iter().for_each(|addr|{
                if !actual_whitelist.wl.contains(addr){
                    actual_whitelist.wl.push((*addr).clone());
                }
            });
            ADDRS_WHITELIST.save(deps.storage, &actual_whitelist)?;
        },

        ExecWhitelist::RemoveAddrs { addresses } => {
            ADDRS_WHITELIST.update(deps.storage, |mut whitelist| {
                if whitelist.wl.len() == 0{
                    return Err(ContractError::EmptyWhitelist {});
                }
                addresses.iter().for_each(|addr| {
                    let index = whitelist.wl.iter().position(|x| *x == *addr);
                    match index {
                        Some(i) => {
                            _ = whitelist.wl.swap_remove(i)
                        },
                        None => (),
                    }
                });

                Ok(whitelist)
            })?;
        },

        ExecWhitelist::ResetWhitelist {} => {
            let empty = Whitelist{wl: Vec::new()};
            ADDRS_WHITELIST.save(deps.storage, &empty)?;
        },

        ExecWhitelist::New { whitelist } => {
            if whitelist.len() == 0 {
                return Err(ContractError::EmptyList {});
            }
            let new_list = Whitelist{wl: whitelist};
            ADDRS_WHITELIST.save(deps.storage, &new_list)?;
        },
    }
    
    Ok(Response::new()
        .add_attribute("method", "try_manage_whitelist"))
}

#[cfg(test)]
mod tests {
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{from_binary, Addr};

    use crate::contract::{instantiate, execute, query};
    use crate::msg::{InstantiateMsg, ExecuteMsg, QueryMsg, ExecWhitelist};
    use crate::state::Whitelist;

    const IBC_ADDR: &str = "IBC_MODULE";
    const GOV_ADDR: &str = "GOV_MODULE";

    const TESTING_ADDR1: &str = "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd";
    const TESTING_ADDR2: &str = "did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn";
   
#[test]
fn execute_add_addrs() {
    let mut deps = mock_dependencies();
    let addr1 = Addr::unchecked(TESTING_ADDR1);
    let addr3 = Addr::unchecked("sender");

    let msg = InstantiateMsg {
        gov_module: Addr::unchecked(GOV_ADDR),
        ibc_module: Addr::unchecked(IBC_ADDR),
        addrs_whitelist: vec![addr1.clone()],
    };
    let info = mock_info(GOV_ADDR, &vec![]);
    let env = mock_env();
    instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

    //add empty list to whitelist
    let add_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::AddAddrs { addresses: vec![] });
    let err = execute(deps.as_mut(), mock_env(), info.clone(), add_addr_msg).unwrap_err();
    assert_eq!(err.to_string(), "Empty list of addresses");

    let query_msg = QueryMsg::GetWhitelist {};
    let mut res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    let mut actual_whitelist: Whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 1);
    assert_eq!(actual_whitelist.wl[0], addr1.clone());

    //add one address to whitelist
    let add_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::AddAddrs {
        addresses: vec![addr3.clone()],
    });
    execute(deps.as_mut(), mock_env(), info.clone(), add_addr_msg).unwrap();

    res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    actual_whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 2);
    assert!(actual_whitelist.wl.contains(&addr1));
    assert!(actual_whitelist.wl.contains(&addr3));

    //Duplications are not allowed
    let add_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::AddAddrs {
        addresses: vec![addr3.clone()],
    });
    execute(deps.as_mut(), mock_env(), info.clone(), add_addr_msg).unwrap();

    res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    actual_whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 2);
    assert!(actual_whitelist.wl.contains(&addr1));
    assert!(actual_whitelist.wl.contains(&addr3));
}

#[test]
fn execute_remove_addrs() {
    let mut deps = mock_dependencies();
    let addr1 = Addr::unchecked(TESTING_ADDR1);
    let addr2 = Addr::unchecked(TESTING_ADDR2);
    let addr3 = Addr::unchecked("sender");

    let msg = InstantiateMsg {
        gov_module: Addr::unchecked(GOV_ADDR),
        ibc_module: Addr::unchecked(IBC_ADDR),
        addrs_whitelist: vec![addr1.clone(), addr2.clone()],
    };
    let info = mock_info(GOV_ADDR, &vec![]);
    let env = mock_env();
    instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

    //remove unknown address from whitelist
    let remove_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::RemoveAddrs {
        addresses: vec![addr3],
    });
    execute(deps.as_mut(), mock_env(), info.clone(), remove_addr_msg).unwrap();

    let query_msg = QueryMsg::GetWhitelist {};

    let mut res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    let mut actual_whitelist : Whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 2);

    //remove addresses from whitelist
    let remove_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::RemoveAddrs {
        addresses: vec![addr1.clone(), addr2.clone()],
    });
    execute(deps.as_mut(), mock_env(), info.clone(), remove_addr_msg).unwrap();

    res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    actual_whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 0);

    //remove from empty whitelist
    let remove_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::RemoveAddrs { 
        addresses: vec![addr2.clone()] ,
    });
    let err = execute(deps.as_mut(), mock_env(), info.clone(), remove_addr_msg).unwrap_err();
    assert_eq!(err.to_string(), "Empty whitelist");
}

#[test]
fn execute_reset_whiltelist() {
    let mut deps = mock_dependencies();
    let addr1 = Addr::unchecked(TESTING_ADDR1);
    let addr2 = Addr::unchecked(TESTING_ADDR2);

    let msg = InstantiateMsg {
        gov_module: Addr::unchecked(GOV_ADDR),
        ibc_module: Addr::unchecked(IBC_ADDR),
        addrs_whitelist: vec![addr1.clone(), addr2.clone()],
    };
    let info = mock_info(GOV_ADDR, &vec![]);
    let env = mock_env();
    instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

    let query_msg = QueryMsg::GetWhitelist {};
    let mut res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    let mut actual_whitelist: Whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 2);
    assert!(actual_whitelist.wl.contains(&addr1));
    assert!(actual_whitelist.wl.contains(&addr2));

    //reset whitelist
    let reset_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::ResetWhitelist {});
    execute(deps.as_mut(), mock_env(), info.clone(), reset_msg).unwrap();

    res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    actual_whitelist = from_binary(&res).unwrap();

    assert!(actual_whitelist.wl.is_empty());
}

#[test]
fn execute_new_whitelist() {
    let mut deps = mock_dependencies();
    let addr1 = Addr::unchecked(TESTING_ADDR1);
    let addr2 = Addr::unchecked(TESTING_ADDR2);
    let addr3 = Addr::unchecked("sender");

    let msg = InstantiateMsg {
        gov_module: Addr::unchecked(GOV_ADDR),
        ibc_module: Addr::unchecked(IBC_ADDR),
        addrs_whitelist: vec![],
    };
    let info = mock_info(GOV_ADDR, &vec![]);
    let env = mock_env();
    instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

    let query_msg = QueryMsg::GetWhitelist {};
    let mut res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    let mut actual_whitelist: Whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 0);

    //set new whitelist
    let new_wl_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::New {
        whitelist: vec![addr1.clone(), addr3.clone(), addr2.clone()],
    });
    execute(deps.as_mut(), mock_env(), info.clone(), new_wl_msg).unwrap();

    res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    actual_whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 3);
    assert!(actual_whitelist.wl.contains(&addr1));
    assert!(actual_whitelist.wl.contains(&addr3));
    assert!(actual_whitelist.wl.contains(&addr2));

    //empty new whitelist
    let new_wl_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::New {
        whitelist: vec![],
    });
    let err = execute(deps.as_mut(), mock_env(), info.clone(), new_wl_msg).unwrap_err();
    assert_eq!(err.to_string(), "Empty list of addresses");

    res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
    actual_whitelist = from_binary(&res).unwrap();

    assert_eq!(actual_whitelist.wl.len(), 3);
}

}
