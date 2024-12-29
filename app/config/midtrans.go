package config

import (
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/spf13/viper"
)

var clientMidtrans *coreapi.Client

type MidtransResponse struct {
	TransactionStatus string `json:"transaction_status"`
}
type MidtransConfig struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

func ReadMidtransconfig() *MidtransConfig {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading .env file: %s", err)
	}

	config := &MidtransConfig{
		IsProduction: false,
	}
	config.ServerKey = viper.GetString("MIDTRANS_SERVER_KEY")
	config.ClientKey = viper.GetString("MIDTRANS_CLIENT_KEY")
	return config
}

func NewMidtransConfig() {
	config := ReadMidtransconfig()
	midtrans.ServerKey = config.ServerKey
	midtrans.ClientKey = config.ClientKey
	midtrans.Environment = midtrans.Sandbox

	clientMidtrans = &coreapi.Client{}
}
