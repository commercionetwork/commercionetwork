// #![cfg(test)]

// use crate::{contract::*, test_msg_recv, test_msg_send};
// use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
// use cosmwasm_std::{from_binary, Addr};

// use crate::msg::{ExecWhitelist, ExecuteMsg, InstantiateMsg, QueryMsg, SudoMsg};
// use crate::state::{Whitelist, GOVMODULE, IBCMODULE, ADDRS_WHITELIST};

// const IBC_ADDR: &str = "IBC_MODULE";
// const GOV_ADDR: &str = "GOV_MODULE";

// const TESTING_ADDR1: &str = "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd";
// const TESTING_ADDR2: &str = "did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn";

// #[test] // Tests we can instantiate the contract and that the owners are set correctly
// fn proper_instantiation() {
//     let mut deps = mock_dependencies();

//     let msg = InstantiateMsg {
//         gov_module: Addr::unchecked(GOV_ADDR),
//         ibc_module: Addr::unchecked(IBC_ADDR),
//         addrs_whitelist: vec![Addr::unchecked(TESTING_ADDR1)],
//     };
//     let info = mock_info(IBC_ADDR, &vec![]);

//     // we can just call .unwrap() to assert this was a success
//     let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
//     assert_eq!(0, res.messages.len());

//     // The ibc and gov modules are properly stored
//     assert_eq!(IBCMODULE.load(deps.as_ref().storage).unwrap(), IBC_ADDR);
//     assert_eq!(GOVMODULE.load(deps.as_ref().storage).unwrap(), GOV_ADDR);
// }

// #[test] // Tests that only GOVMODULE and IBCMODULE can execute the contract
// fn execute_privileges() {
//     let mut deps = mock_dependencies();
//     let addr1 = Addr::unchecked(TESTING_ADDR1);

//     let msg = InstantiateMsg {
//         gov_module: Addr::unchecked(GOV_ADDR),
//         ibc_module: Addr::unchecked(IBC_ADDR),
//         addrs_whitelist: vec![],
//     };
//     let info = mock_info(GOV_ADDR, &vec![]);
//     let env = mock_env();
//     instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

//     //GOVMODULE can execute
//     let add_addr_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::AddAddrs {
//         addresses: vec![addr1.clone()],
//     });
//     execute(
//         deps.as_mut(),
//         mock_env(),
//         info.clone(),
//         add_addr_msg.clone(),
//     )
//     .unwrap();

//     //IBCMODULE can execute
//     let info = mock_info(IBC_ADDR, &vec![]);
//     execute(
//         deps.as_mut(),
//         mock_env(),
//         info.clone(),
//         add_addr_msg.clone(),
//     )
//     .unwrap();

//     //Random addr cannot execute
//     let info = mock_info("senderAddress", &vec![]);
//     let err = execute(deps.as_mut(), mock_env(), info.clone(), add_addr_msg).unwrap_err();
//     assert_eq!(
//         err.to_string(),
//         "Unauthorized! The sender senderAddress has not authorization"
//     )
// }

// #[test] // Tests that when a packet is transferred, the has permission/whitelisted
// fn sending_permission() {
//     let mut deps = mock_dependencies();
//     let addr1 = Addr::unchecked(TESTING_ADDR1);
//     let addr2 = Addr::unchecked(TESTING_ADDR2);

//     let empty_wl: Vec<Addr> = vec![];
//     let wl1 = vec![addr1.clone(), addr2.clone()];
//     let wl2 = vec![addr2.clone()];

//     let msg = InstantiateMsg {
//         gov_module: Addr::unchecked(GOV_ADDR),
//         ibc_module: Addr::unchecked(IBC_ADDR),
//         addrs_whitelist: empty_wl,
//     };
//     let info = mock_info(GOV_ADDR, &vec![]);
//     let env = mock_env();
//     let _res = instantiate(deps.as_mut(), env.clone(), info.clone(), msg.clone()).unwrap();

//     //empty whitelist -> everyone can send
//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: None
//     );
//     sudo(deps.as_mut(), mock_env(), send_msg).unwrap();

//     // only addr1 and addr2 can transfer tokens
//     let new_wl_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::New { whitelist: wl1 });
//     execute(deps.as_mut(), mock_env(), info.clone(), new_wl_msg).unwrap();

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: Some(addr1.clone())
//     );
//     sudo(deps.as_mut(), mock_env(), send_msg).unwrap();

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: Some(addr2.clone())
//     );
//     sudo(deps.as_mut(), mock_env(), send_msg).unwrap();

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: None
//     );
//     sudo(deps.as_mut(), mock_env(), send_msg.clone()).unwrap_err();
//     let err = sudo(deps.as_mut(), mock_env(), send_msg).unwrap_err();
//     assert_eq!(
//         err.to_string(),
//         "Unauthorized! The sender senderAddress has not authorization"
//     );

//     // only addr2 can transfer tokens
//     let new_wl_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::New { whitelist: wl2 });
//     execute(deps.as_mut(), mock_env(), info, new_wl_msg).unwrap();

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: Some(addr1.clone())
//     );
//     let err = sudo(deps.as_mut(), mock_env(), send_msg).unwrap_err();
//     assert_eq!(err.to_string(), "Unauthorized! The sender did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd has not authorization");

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: Some(addr2.clone())
//     );
//     sudo(deps.as_mut(), mock_env(), send_msg).unwrap();

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: None
//     );
//     let err = sudo(deps.as_mut(), mock_env(), send_msg).unwrap_err();
//     assert_eq!(
//         err.to_string(),
//         "Unauthorized! The sender senderAddress has not authorization"
//     );
// }

// #[test] // Tests we can get the current state of the trackers
// fn query_state() {
//     let mut deps = mock_dependencies();
//     let addr1 = Addr::unchecked(TESTING_ADDR1);
//     let addr2 = Addr::unchecked(TESTING_ADDR2);

//     let msg = InstantiateMsg {
//         gov_module: Addr::unchecked(GOV_ADDR),
//         ibc_module: Addr::unchecked(IBC_ADDR),
//         addrs_whitelist: vec![addr1.clone(), addr2.clone()],
//     };
//     let info = mock_info(GOV_ADDR, &vec![]);
//     let env = mock_env();
//     let _res = instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

//     let query_msg = QueryMsg::GetWhitelist {};

//     let mut res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
//     let actual_whitelist: Whitelist = from_binary(&res).unwrap();
//     assert_eq!(actual_whitelist.wl.len(), 2);
//     assert!(actual_whitelist.wl.contains(&addr1));
//     assert!(actual_whitelist.wl.contains(&addr2));

//     let send_msg = test_msg_send!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 300_u32.into(),
//         sender: Some(addr1.clone())
//     );
//     sudo(deps.as_mut(), mock_env(), send_msg.clone()).unwrap();

//     let recv_msg = test_msg_recv!(
//         channel_id: format!("channel"),
//         denom: format!("denom"),
//         funds: 30_u32.into()
//     );
//     sudo(deps.as_mut(), mock_env(), recv_msg.clone()).unwrap();

//     //Query
//     let new_wl_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::ResetWhitelist {});
//     execute(deps.as_mut(), mock_env(), info.clone(), new_wl_msg).unwrap();

//     res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
//     let actual_whitelist: Whitelist = from_binary(&res).unwrap();
//     assert!(actual_whitelist.wl.is_empty());

//     assert_eq!(ADDRS_WHITELIST.load(deps.as_ref().storage).unwrap(), actual_whitelist);

//     let new_wl_msg = ExecuteMsg::ManageWhitelist(ExecWhitelist::New {
//         whitelist: vec![addr2.clone()],
//     });
//     execute(deps.as_mut(), mock_env(), info, new_wl_msg).unwrap();
//     res = query(deps.as_ref(), mock_env(), query_msg.clone()).unwrap();
//     let actual_whitelist: Whitelist = from_binary(&res).unwrap();
//     assert_eq!(actual_whitelist.wl.len(), 1);
//     assert_eq!(actual_whitelist.wl[0], addr2);

//     assert_eq!(ADDRS_WHITELIST.load(deps.as_ref().storage).unwrap(), actual_whitelist);
// }


// #[test]
// fn test_basic_message() {
//     let json = r#"{"send_packet":{"packet":{"sequence":2,"source_port":"transfer","source_channel":"channel-0","destination_port":"transfer","destination_channel":"channel-0","data":{"denom":"ucommercio","amount":"125000000000011250","sender":"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd","receiver":"did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn"},"timeout_height":{"revision_height":100}}}}"#;
//     let _parsed: SudoMsg = serde_json_wasm::from_str(json).unwrap();
//     //println!("{parsed:?}");
// }

// #[test]
// fn test_testnet_message() {
//     let json = r#"{"send_packet":{"packet":{"sequence":4,"source_port":"transfer","source_channel":"channel-0","destination_port":"transfer","destination_channel":"channel-1491","data":{"denom":"ucommercio","amount":"100","sender":"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd","receiver":"did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn"},"timeout_height":{},"timeout_timestamp":1668024637477293371}}}"#;
//     let _parsed: SudoMsg = serde_json_wasm::from_str(json).unwrap();
//     //println!("{parsed:?}");
// }

// #[test]
// fn test_tokenfactory_message() {
//     let json = r#"{"send_packet":{"packet":{"sequence":4,"source_port":"transfer","source_channel":"channel-0","destination_port":"transfer","destination_channel":"channel-1491","data":{"denom":"transfer/channel-0/factory/osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj/czar","amount":"100000000000000000","sender":"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd","receiver":"did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn"},"timeout_height":{},"timeout_timestamp":1668024476848430980}}}"#;
//     let _parsed: SudoMsg = serde_json_wasm::from_str(json).unwrap();
//     //println!("{parsed:?}");
// }

