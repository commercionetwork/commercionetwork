module github.com/commercionetwork/commercionetwork

go 1.16

require (
	github.com/CosmWasm/wasmd v0.18.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/armon/go-metrics v0.3.9 // indirect
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/cosmos/cosmos-sdk v0.44.1
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gogo/protobuf v1.3.3
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cast v1.4.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.3
	github.com/tendermint/spm v0.0.0-20210524110815-6d7452d2dc4a
	github.com/tendermint/tendermint v0.34.12
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.40.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.42.9

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
