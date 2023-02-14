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

/*#[cfg(test)]
mod tests {
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{from_binary, Addr, StdError};

    use crate::contract::{execute, query};
    use crate::helpers::tests::verify_query_response;
    use crate::msg::{ExecuteMsg, QueryMsg, QuotaMsg};
    use crate::state::{RateLimit, GOVMODULE, IBCMODULE};

    const IBC_ADDR: &str = "IBC_MODULE";
    const GOV_ADDR: &str = "GOV_MODULE";

    #[test] // Tests AddPath and RemovePath messages
    fn management_add_and_remove_path() {
        let mut deps = mock_dependencies();
        IBCMODULE
            .save(deps.as_mut().storage, &Addr::unchecked(IBC_ADDR))
            .unwrap();
        GOVMODULE
            .save(deps.as_mut().storage, &Addr::unchecked(GOV_ADDR))
            .unwrap();

        let msg = ExecuteMsg::AddPath {
            channel_id: format!("channel"),
            denom: format!("denom"),
            quotas: vec![QuotaMsg {
                name: "daily".to_string(),
                duration: 1600,
                send_recv: (3, 5),
            }],
        };
        let info = mock_info(IBC_ADDR, &vec![]);

        let env = mock_env();
        let res = execute(deps.as_mut(), env.clone(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        let query_msg = QueryMsg::GetQuotas {
            channel_id: format!("channel"),
            denom: format!("denom"),
        };

        let res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();

        let value: Vec<RateLimit> = from_binary(&res).unwrap();
        verify_query_response(
            &value[0],
            "daily",
            (3, 5),
            1600,
            0_u32.into(),
            0_u32.into(),
            env.block.time.plus_seconds(1600),
        );

        assert_eq!(value.len(), 1);

        // Add another path
        let msg = ExecuteMsg::AddPath {
            channel_id: format!("channel2"),
            denom: format!("denom"),
            quotas: vec![QuotaMsg {
                name: "daily".to_string(),
                duration: 1600,
                send_recv: (3, 5),
            }],
        };
        let info = mock_info(IBC_ADDR, &vec![]);

        let env = mock_env();
        execute(deps.as_mut(), env.clone(), info, msg).unwrap();

        // remove the first one
        let msg = ExecuteMsg::RemovePath {
            channel_id: format!("channel"),
            denom: format!("denom"),
        };

        let info = mock_info(IBC_ADDR, &vec![]);
        let env = mock_env();
        execute(deps.as_mut(), env.clone(), info, msg).unwrap();

        // The channel is not there anymore
        let err = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap_err();
        assert!(matches!(err, StdError::NotFound { .. }));

        // The second channel is still there
        let query_msg = QueryMsg::GetQuotas {
            channel_id: format!("channel2"),
            denom: format!("denom"),
        };
        let res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
        let value: Vec<RateLimit> = from_binary(&res).unwrap();
        assert_eq!(value.len(), 1);
        verify_query_response(
            &value[0],
            "daily",
            (3, 5),
            1600,
            0_u32.into(),
            0_u32.into(),
            env.block.time.plus_seconds(1600),
        );

        // Paths are overriden if they share a name and denom
        let msg = ExecuteMsg::AddPath {
            channel_id: format!("channel2"),
            denom: format!("denom"),
            quotas: vec![QuotaMsg {
                name: "different".to_string(),
                duration: 5000,
                send_recv: (50, 30),
            }],
        };
        let info = mock_info(IBC_ADDR, &vec![]);

        let env = mock_env();
        execute(deps.as_mut(), env.clone(), info, msg).unwrap();

        let query_msg = QueryMsg::GetQuotas {
            channel_id: format!("channel2"),
            denom: format!("denom"),
        };
        let res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
        let value: Vec<RateLimit> = from_binary(&res).unwrap();
        assert_eq!(value.len(), 1);

        verify_query_response(
            &value[0],
            "different",
            (50, 30),
            5000,
            0_u32.into(),
            0_u32.into(),
            env.block.time.plus_seconds(5000),
        );
    }
}*/
