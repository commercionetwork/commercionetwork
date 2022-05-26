rm -rf ~/.ignite/relayer

ignite relayer configure -a \
--source-rpc "http://0.0.0.0:26657" \
--source-faucet "http://0.0.0.0:4500" \
--source-port "transfer" \
--source-version "ics20-1" \
--source-gasprice "0.15uccc" \
--source-prefix "did:com:" \
--source-gaslimit 300000 \
--target-rpc "http://0.0.0.0:26659" \
--target-faucet "http://0.0.0.0:4501" \
--target-port "transfer" \
--target-version "ics20-1" \
--target-gasprice "0.00025stake" \
--target-prefix "cosmos" \
--target-gaslimit 300000