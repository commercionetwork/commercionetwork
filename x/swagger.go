//go:generate go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal -g swagger.go -o ../swagger

package x

import (
	_ "github.com/cosmos/cosmos-sdk/types"
)

// @title CommercioNetwork API
// @description Swagger API for CommercioNetwork
// @contact.email developers@commercio.network

// @BasePath /

type JSONResult struct {
	Height string      `json:"height" example:"1234"`
	Result interface{} `json:"result"`
}
