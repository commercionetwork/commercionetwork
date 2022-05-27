rm -rf ~/.ignite/relayer

ignite relayer configure -a \
--source-rpc "http://0.0.0.0:26659" \
--source-faucet "http://0.0.0.0:4501" \
--source-port "transfer" \
--source-version "ics20-1" \
--source-gasprice "0.00025stake" \
--source-prefix "cosmos" \
--source-gaslimit 300000 \
--target-rpc "http://0.0.0.0:26657" \
--target-faucet "http://0.0.0.0:4500" \
--target-port "transfer" \
--target-version "ics20-1" \
--target-gasprice "0.15uccc" \
--target-prefix "did:com:" \
--target-gaslimit 300000



