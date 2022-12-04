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
	ERC20Predicate   common.Address
}

type DebugConfig struct {
	Enable bool
	Level  log.Level
}

func NewDefaultConfig(network Network) POSClientConfig {
	switch network {
	case MainNet:
		return POSClientConfig{
			Child: ChildConfig{
				Rpc: "https://rpc.ankr.com/polygon",
			},
			Root: RootConfig{
				Rpc:              "https://rpc.ankr.com/eth",
				RootChain:        common.HexToAddress("0x536c55cFe4892E581806e10b38dFE8083551bd03"),
				RootChainManager: common.HexToAddress("0x37D26DC2890b35924b40574BAc10552794771997"),
				ERC20Predicate:   common.HexToAddress("0x40ec5B33f54e0E8A33A975908C5BA1c14e5BbbDf"),
			},
			Debug: DebugConfig{
				Enable: true,
				Level:  DebugLevel,
			},
		}
	case TestNet:
		return POSClientConfig{
			Child: ChildConfig{
				Rpc: "https://rpc.ankr.com/polygon_mumbai",
			},
			Root: RootConfig{
				Rpc:              "https://rpc.ankr.com/eth_goerli",
				RootChain:        common.HexToAddress("0x2890bA17EfE978480615e330ecB65333b880928e"),
				RootChainManager: common.HexToAddress("0xBbD7cBFA79faee899Eaf900F13C9065bF03B1A74"),
				ERC20Predicate:   common.HexToAddress("0xdD6596F2029e6233DEFfaCa316e6A95217d4Dc34"),
			},
			Debug: DebugConfig{
				Enable: true,
				Level:  DebugLevel,
			},
		}
	}
	return POSClientConfig{}
}
