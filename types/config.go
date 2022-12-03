package types

import "github.com/ethereum/go-ethereum/common"

type Network int

const (
	MainNet = Network(1)
	TestNet = Network(5)
)

type POSClientConfig struct {
	Child ChildConfig
	Root  RootConfig
	Debug bool
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
			Debug: true,
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
				ERC20Predicate:   common.HexToAddress("0x7B276A55987E3020026Bb098F15E968313Bd1aF2"),
			},
			Debug: true,
		}
	}
	return POSClientConfig{}
}
