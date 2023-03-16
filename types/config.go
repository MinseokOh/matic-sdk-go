package types

import (
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

type Network int

const (
	MainNet = Network(1)
	TestNet = Network(5)
)

type NetworkType int

const (
	Root  = NetworkType(1)
	Child = NetworkType(2)
)

type POSClientConfig struct {
	Child ChildConfig
	Root  RootConfig
	Debug DebugConfig
}

type ChildConfig struct {
	Rpc string
}

type RootConfig struct {
	Rpc              string
	RootChain        common.Address
	RootChainManager common.Address
}

type DebugConfig struct {
	Enable bool
	Level  log.Level
}

func NewDefaultConfig(network Network) POSClientConfig {
	contract := GetContractByNetwork(network)
	return POSClientConfig{
		Child: contract.ChildConfig("https://rpc.ankr.com/polygon"),
		Root:  contract.RootConfig("https://rpc.ankr.com/eth"),
		Debug: DebugConfig{
			Enable: true,
			Level:  DebugLevel,
		},
	}
}
