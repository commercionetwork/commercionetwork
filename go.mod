module github.com/commercionetwork/commercionetwork

go 1.16

require (
	github.com/CosmWasm/wasmd v0.27.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cosmos/cosmos-sdk v0.45.4
	github.com/cosmos/go-bip39 v1.0.0
	//github.com/cosmos/ibc-go/v2 v2.2.0
	github.com/cosmos/ibc-go/v3 v3.0.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/swag v1.7.3
	github.com/tendermint/spm v0.0.0-20210524110815-6d7452d2dc4a
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.7
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
