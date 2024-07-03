use crate::state::FlowType;
use cosmwasm_std::{Addr, Uint256};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Height {
    /// Previously known as "epoch"
    revision_number: Option<u64>,

    /// The height of a block
    revision_height: Option<u64>,
}

// IBC transfer data
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct FungibleTokenData {
    pub denom: String,
    amount: Uint256,
    pub sender: Addr,
    receiver: Addr,
}

// An IBC packet
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Packet {
    pub sequence: u64,
    pub source_port: String,
    pub source_channel: String,
    pub destination_port: String,
    pub destination_channel: String,
    pub data: FungibleTokenData,
    pub timeout_height: Height,
    pub timeout_timestamp: Option<u64>,
}

fn hash_denom(denom: &str) -> String {
    let mut hasher = Sha256::new();
    hasher.update(denom.as_bytes());
    let result = hasher.finalize();
    let hash = hex::encode(result);
    format!("ibc/{}", hash.to_uppercase())
}

impl Packet {
    pub fn mock(
        source_channel: String,
        dest_channel: String,
        denom: String,
        funds: Uint256,
        sender: Option<Addr>,
    ) -> Packet {
        let mut send_packet = Packet {
            sequence: 0,
            source_port: "transfer".to_string(),
            source_channel,
            destination_port: "transfer".to_string(),
            destination_channel: dest_channel,
            data: crate::packet::FungibleTokenData {
                denom,
                amount: funds,
                sender: Addr::unchecked("senderAddress"),
                receiver: Addr::unchecked("receiverAddress"),
            },
            timeout_height: crate::packet::Height {
                revision_number: None,
                revision_height: None,
            },
            timeout_timestamp: None,
        };

        match sender {
            Some(sender_addr) => {
                send_packet.data.sender = sender_addr;
            }
            None => ()
        }

        send_packet
    }

    pub fn get_funds(&self) -> Uint256 {
        self.data.amount
    }

    fn local_channel(&self, direction: &FlowType) -> String {
        // Pick the appropriate channel depending on whether this is a send or a recv
        match direction {
            FlowType::In => self.destination_channel.clone(),
            FlowType::Out => self.source_channel.clone(),
        }
    }

    fn receiver_chain_is_source(&self) -> bool {
        self.data
            .denom
            .starts_with(&format!("transfer/{}", self.source_channel))
    }

    fn handle_denom_for_sends(&self) -> String {
        if !self.data.denom.starts_with("transfer/") {
            // For native tokens we just use what's on the packet
            return self.data.denom.clone();
        }
        // For non-native tokens, we need to generate the IBCDenom
        hash_denom(&self.data.denom)
    }

    fn handle_denom_for_recvs(&self) -> String {
        if self.receiver_chain_is_source() {
            // These are tokens that have been sent to the counterparty and are returning
            let unprefixed = self
                .data
                .denom
                .strip_prefix(&format!("transfer/{}/", self.source_channel))
                .unwrap_or_default();
            let split: Vec<&str> = unprefixed.split('/').collect();
            if split[0] == unprefixed {
                // This is a native token. Return the unprefixed token
                unprefixed.to_string()
            } else {
                // This is a non-native that was sent to the counterparty.
                // We need to hash it.
                // The ibc-go implementation checks that the denom has been built correctly. We
                // don't need to do that here because if it hasn't, the transfer module will catch it.
                hash_denom(unprefixed)
            }
        } else {
            // Tokens that come directly from the counterparty.
            // Since the sender didn't prefix them, we need to do it here.
            let prefixed = format!("transfer/{}/", self.destination_channel) + &self.data.denom;
            hash_denom(&prefixed)
        }
    }

    fn local_denom(&self, direction: &FlowType) -> String {
        match direction {
            FlowType::In => self.handle_denom_for_recvs(),
            FlowType::Out => self.handle_denom_for_sends(),
        }
    }

    pub fn path_data(&self, direction: &FlowType) -> (String, String) {
        (self.local_channel(direction), self.local_denom(direction))
    }
}

// Helpers

// Create a new packet for testing
#[cfg(test)]
#[macro_export]
macro_rules! test_msg_send {
    (channel_id: $channel_id:expr, denom: $denom:expr, funds: $funds:expr, sender: $sender:expr) => {
        $crate::msg::SudoMsg::SendPacket {
            packet: $crate::packet::Packet::mock($channel_id, $channel_id, $denom, $funds, $sender),
        }
    };
}

#[cfg(test)]
#[macro_export]
macro_rules! test_msg_recv {
    (channel_id: $channel_id:expr, denom: $denom:expr, funds: $funds:expr) => {
        $crate::msg::SudoMsg::RecvPacket {
            packet: $crate::packet::Packet::mock(
                $channel_id,
                $channel_id,
                format!("transfer/{}/{}", $channel_id, $denom),
                $funds,
                None,
            ),
        }
    };
}

#[cfg(test)]
pub mod tests {
    use crate::msg::SudoMsg;

    use super::*;

    #[test]
    fn send_native() {
        let packet = Packet::mock(
            format!("channel-17-local"),
            format!("channel-42-counterparty"),
            format!("ucommercio"),
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::Out), "ucommercio");
    }

    #[test]
    fn send_non_native() {
        // The transfer module "unhashes" the denom from
        // ibc/09E4864A262249507925831FBAD69DAD08F66FAAA0640714E765912A0751289A
        // to port/channel/denom before passing it along to the contract
        let packet = Packet::mock(
            format!("channel-17-local"),
            format!("channel-42-counterparty"),
            format!("transfer/channel-17-local/ujuno"),
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::Out),
            "ibc/09E4864A262249507925831FBAD69DAD08F66FAAA0640714E765912A0751289A"
        );
    }

    #[test]
    fn receive_non_native() {
        // The counterparty chain sends their own native token to us
        let packet = Packet::mock(
            format!("channel-42-counterparty"), // The counterparty's channel is the source here
            format!("channel-17-local"),        // Our channel is the dest channel
            format!("ujuno"),                   // This is unwrapped. It is our job to wrap it
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::In),
            "ibc/09E4864A262249507925831FBAD69DAD08F66FAAA0640714E765912A0751289A"
        );
    }

    #[test]
    fn receive_native() {
        // The counterparty chain sends us back our native token that they had wrapped
        let packet = Packet::mock(
            format!("channel-42-counterparty"), // The counterparty's channel is the source here
            format!("channel-17-local"),        // Our channel is the dest channel
            format!("transfer/channel-42-counterparty/ucommercio"),
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::In), "ucommercio");
    }

    // Let's assume we have two chains A and B (local and counterparty) connected in the following way:
    //
    // Chain A <---> channel-17-local <---> channel-42-counterparty <---> Chain B
    //
    // The following tests should pass
    //

    const WRAPPED_COM_ON_HUB_TRACE: &str = "transfer/channel-141/ucommercio";
    const WRAPPED_ATOM_ON_COMMERCIO_TRACE: &str = "transfer/channel-0/uatom";
    const WRAPPED_ATOM_ON_COMMERCIO_HASH: &str =
        "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2";
    const WRAPPED_COM_ON_HUB_HASH: &str =
        "ibc/B445D4E44F7EB8A84095F22803FF5B31DC81D0B5A10A7F0CFD77537001351A0D";

    #[test]
    fn sanity_check() {
        // uatom sent to commercionetwork
        let packet = Packet::mock(
            format!("channel-141"), // from: hub
            format!("channel-0"),   // to: commercionetwork
            format!("uatom"),
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::In),
            WRAPPED_ATOM_ON_COMMERCIO_HASH.clone()
        );

        // uatom on commercionetwork sent back to the hub
        let packet = Packet::mock(
            format!("channel-0"),                      // from: commercionetwork
            format!("channel-141"),                    // to: hub
            WRAPPED_ATOM_ON_COMMERCIO_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::In), "uatom");

        // COM sent to the hub
        let packet = Packet::mock(
            format!("channel-0"),   // from: commercionetwork
            format!("channel-141"), // to: hub
            format!("ucommercio"),
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::Out), "ucommercio");

        // COM on the hub sent back to commercionetwork
        // send
        let packet = Packet::mock(
            format!("channel-141"),                // from: hub
            format!("channel-0"),                  // to: commercionetwork
            WRAPPED_COM_ON_HUB_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::Out), WRAPPED_COM_ON_HUB_HASH);

        // receive
        let packet = Packet::mock(
            format!("channel-141"),                // from: hub
            format!("channel-0"),                  // to: commercionetwork
            WRAPPED_COM_ON_HUB_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::In), "ucommercio");

        // Now let's pretend we're the hub.
        // The following tests are from perspective of the hub (i.e.: if this contract were deployed there)
        //
        // COM sent to the hub
        let packet = Packet::mock(
            format!("channel-0"),   // from: commercionetwork
            format!("channel-141"), // to: hub
            format!("ucommercio"),
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::In), WRAPPED_COM_ON_HUB_HASH);

        // COM on the hub sent back to the commercionetwork
        let packet = Packet::mock(
            format!("channel-141"),                // from: hub
            format!("channel-0"),                  // to: commercionetwork
            WRAPPED_COM_ON_HUB_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::In), "ucommercio");

        // uatom sent to commercionetwork
        let packet = Packet::mock(
            format!("channel-141"), // from: hub
            format!("channel-0"),   // to: commercionetwork
            format!("uatom"),
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::Out), "uatom");

        // utaom on the commercionetwork sent back to the hub
        // send
        let packet = Packet::mock(
            format!("channel-0"),                      // from: commercionetwork
            format!("channel-141"),                    // to: hub
            WRAPPED_ATOM_ON_COMMERCIO_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::Out),
            WRAPPED_ATOM_ON_COMMERCIO_HASH
        );

        // receive
        let packet = Packet::mock(
            format!("channel-0"),                      // from: commercionetwork
            format!("channel-141"),                    // to: hub
            WRAPPED_ATOM_ON_COMMERCIO_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(packet.local_denom(&FlowType::In), "uatom");
    }

    #[test]
    fn sanity_double() {
        // Now let's deal with double wrapping

        let juno_wrapped_commercio_wrapped_atom_hash =
            "ibc/6CDD4663F2F09CD62285E2D45891FC149A3568E316CE3EBBE201A71A78A69388";

        // Send uatom on stored on commercionetwork to juno
        // send
        let packet = Packet::mock(
            format!("channel-42"),                     // from: commercionetwork
            format!("channel-0"),                      // to: juno
            WRAPPED_ATOM_ON_COMMERCIO_TRACE.to_string(), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::Out),
            WRAPPED_ATOM_ON_COMMERCIO_HASH
        );

        // receive
        let packet = Packet::mock(
            format!("channel-42"), // from: commercionetwork
            format!("channel-0"),  // to: juno
            WRAPPED_ATOM_ON_COMMERCIO_TRACE.to_string(),
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::In),
            juno_wrapped_commercio_wrapped_atom_hash
        );

        // Send back that multi-wrapped token to commercionetwork
        // send
        let packet = Packet::mock(
            format!("channel-0"),  // from: juno
            format!("channel-42"), // to: commercionetwork
            format!("{}{}", "transfer/channel-0/", WRAPPED_ATOM_ON_COMMERCIO_TRACE), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::Out),
            juno_wrapped_commercio_wrapped_atom_hash
        );

        // receive
        let packet = Packet::mock(
            format!("channel-0"),  // from: juno
            format!("channel-42"), // to: commercionetwork
            format!("{}{}", "transfer/channel-0/", WRAPPED_ATOM_ON_COMMERCIO_TRACE), // unwrapped before reaching the contract
            0_u128.into(),
            None,
        );
        assert_eq!(
            packet.local_denom(&FlowType::In),
            WRAPPED_ATOM_ON_COMMERCIO_HASH
        );
    }

    #[test]
    fn tokenfactory_packet() {
        let json = r#"{"send_packet":{"packet":{"sequence":4,"source_port":"transfer","source_channel":"channel-0","destination_port":"transfer","destination_channel":"channel-1491","data":{"denom":"transfer/channel-0/factory/osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj/czar","amount":"100000000000000000","sender":"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd","receiver":"did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn"},"timeout_height":{},"timeout_timestamp":1668024476848430980}}}"#;
        let parsed: SudoMsg = serde_json_wasm::from_str(json).unwrap();
        //println!("{parsed:?}");

        match parsed {
            SudoMsg::SendPacket { packet } => {
                assert_eq!(
                    packet.local_denom(&FlowType::Out),
                    "ibc/07A1508F49D0753EDF95FA18CA38C0D6974867D793EB36F13A2AF1A5BB148B22"
                );
            }
            _ => panic!("parsed into wrong variant"),
        }
    }

    #[test]
    fn packet_with_memo() {
        // extra fields (like memo) get ignored.
        let json = r#"{"recv_packet":{"packet":{"sequence":1,"source_port":"transfer","source_channel":"channel-0","destination_port":"transfer","destination_channel":"channel-0","data":{"denom":"ucommercio","amount":"1","sender":"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd","receiver":"did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn","memo":"some info"},"timeout_height":{"revision_height":100}}}}"#;
        let _parsed: SudoMsg = serde_json_wasm::from_str(json).unwrap();
        //println!("{parsed:?}");
    }
}
