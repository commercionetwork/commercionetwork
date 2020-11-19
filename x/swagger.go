//go:generate go run github.com/swaggo/swag/cmd/swag init --parseDependency -g swagger.go -o ../swagger

package x

import (
	_ "github.com/cosmos/cosmos-sdk/types"
)

// @title CommercioNetwork API
// @description Swagger API for CommercioNetwork
// @contact.name Gianguido Sorà
// @contact.email me+work@gsora.xyz

// @BasePath /

type JSONResult struct {
	Height string      `json:"height" example:"1234"`
	Result interface{} `json:"result"`
}
