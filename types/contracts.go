package types

import (
	"github.com/ethereum/go-ethereum/common"
)

const MaticAddress = "0x0000000000000000000000000000000000001010"

const TestNetContractURL = `https://static.matic.network/network/testnet/mumbai/index.json`
const MainNetContractURL = `https://static.matic.network/network/mainnet/v1/index.json`

type Contract struct {
	Main struct {
		NetworkName       string                `json:"NetworkName"`
		ChainID           int                   `json:"ChainId"`
		DaggerEndpoint    string                `json:"DaggerEndpoint"`
		WatcherAPI        string                `json:"WatcherAPI"`
		StakingAPI        string                `json:"StakingAPI"`
		Explorer          string                `json:"Explorer"`
		SubgraphURL       string                `json:"SubgraphUrl"`
		SupportsEIP1559   bool                  `json:"SupportsEIP1559"`
		Owner             string                `json:"Owner"`
		Contracts         RootContracts         `json:"Contracts"`
		POSContracts      POSRootContracts      `json:"POSContracts"`
		FxPortalContracts FxPortalRootContracts `json:"FxPortalContracts"`
	} `json:"Main"`
	Matic struct {
		NetworkName       string                 `json:"NetworkName"`
		ChainID           int                    `json:"ChainId"`
		RPC               string                 `json:"RPC"`
		DaggerEndpoint    string                 `json:"DaggerEndpoint"`
		Explorer          string                 `json:"Explorer"`
		NetworkAPI        string                 `json:"NetworkAPI"`
		SupportsEIP1559   bool                   `json:"SupportsEIP1559"`
		Contracts         ChildContracts         `json:"Contracts"`
		POSContracts      POSChildContracts      `json:"POSContracts"`
		FxPortalContracts FxPortalChildContracts `json:"FxPortalContracts"`
		GenesisContracts  GenesisChildContracts  `json:"GenesisContracts"`
	} `json:"Matic"`
	Heimdall struct {
		ChainID string `json:"ChainId"`
		API     string `json:"API"`
	} `json:"Heimdall"`
}

func (c Contract) RootConfig(rpc string) RootConfig {
	return RootConfig{
		Rpc:              rpc,
		RootChain:        common.HexToAddress("0x2890bA17EfE978480615e330ecB65333b880928e"),
		RootChainManager: common.HexToAddress(c.Main.POSContracts.RootChainManager),
	}
}

func (c Contract) ChildConfig(rpc string) ChildConfig {
	return ChildConfig{
		Rpc: rpc,
	}
}

type ChildContracts struct {
	EIP1559Burn string `json:"EIP1559Burn"`
	ChildChain  string `json:"ChildChain"`
	Tokens      struct {
		MaticWeth  string `json:"MaticWeth"`
		MaticToken string `json:"MaticToken"`
		TestToken  string `json:"TestToken"`
		RootERC721 string `json:"RootERC721"`
		Wmatic     string `json:"WMATIC"`
	} `json:"Tokens"`
}

type POSChildContracts struct {
	ChildChainManager      string `json:"ChildChainManager"`
	ChildChainManagerProxy string `json:"ChildChainManagerProxy"`
	Tokens                 struct {
		DummyERC20           string `json:"DummyERC20"`
		DummyERC721          string `json:"DummyERC721"`
		DummyERC1155         string `json:"DummyERC1155"`
		DummyMintableERC20   string `json:"DummyMintableERC20"`
		DummyMintableERC721  string `json:"DummyMintableERC721"`
		DummyMintableERC1155 string `json:"DummyMintableERC1155"`
		MaticWETH            string `json:"MaticWETH"`
	} `json:"Tokens"`
}

type RootContracts struct {
	BytesLib              string `json:"BytesLib"`
	Common                string `json:"Common"`
	ECVerify              string `json:"ECVerify"`
	Merkle                string `json:"Merkle"`
	MerklePatriciaProof   string `json:"MerklePatriciaProof"`
	PriorityQueue         string `json:"PriorityQueue"`
	RLPEncode             string `json:"RLPEncode"`
	RLPReader             string `json:"RLPReader"`
	SafeMath              string `json:"SafeMath"`
	Governance            string `json:"Governance"`
	GovernanceProxy       string `json:"GovernanceProxy"`
	Registry              string `json:"Registry"`
	RootChain             string `json:"RootChain"`
	RootChainProxy        string `json:"RootChainProxy"`
	ValidatorShareFactory string `json:"ValidatorShareFactory"`
	StakingInfo           string `json:"StakingInfo"`
	StakingNFT            string `json:"StakingNFT"`
	StakeManager          string `json:"StakeManager"`
	StakeManagerProxy     string `json:"StakeManagerProxy"`
	SlashingManager       string `json:"SlashingManager"`
	ValidatorShare        string `json:"ValidatorShare"`
	StateSender           string `json:"StateSender"`
	DepositManager        string `json:"DepositManager"`
	DepositManagerProxy   string `json:"DepositManagerProxy"`
	WithdrawManager       string `json:"WithdrawManager"`
	WithdrawManagerProxy  string `json:"WithdrawManagerProxy"`
	ExitNFT               string `json:"ExitNFT"`
	ERC20Predicate        string `json:"ERC20Predicate"`
	ERC721Predicate       string `json:"ERC721Predicate"`
	EIP1559Burn           string `json:"EIP1559Burn"`
	Tokens                struct {
		MaticToken string `json:"MaticToken"`
		TestToken  string `json:"TestToken"`
		RootERC721 string `json:"RootERC721"`
		MaticWeth  string `json:"MaticWeth"`
	} `json:"Tokens"`
}

type POSRootContracts struct {
	Merkle                        string `json:"Merkle"`
	MerklePatriciaProof           string `json:"MerklePatriciaProof"`
	RLPReader                     string `json:"RLPReader"`
	SafeERC20                     string `json:"SafeERC20"`
	RootChainManager              string `json:"RootChainManager"`
	RootChainManagerProxy         string `json:"RootChainManagerProxy"`
	DummyStateSender              string `json:"DummyStateSender"`
	ERC20Predicate                string `json:"ERC20Predicate"`
	ERC20PredicateProxy           string `json:"ERC20PredicateProxy"`
	ERC721Predicate               string `json:"ERC721Predicate"`
	ERC721PredicateProxy          string `json:"ERC721PredicateProxy"`
	ERC1155Predicate              string `json:"ERC1155Predicate"`
	ERC1155PredicateProxy         string `json:"ERC1155PredicateProxy"`
	EtherPredicate                string `json:"EtherPredicate"`
	EtherPredicateProxy           string `json:"EtherPredicateProxy"`
	MintableERC20Predicate        string `json:"MintableERC20Predicate"`
	MintableERC20PredicateProxy   string `json:"MintableERC20PredicateProxy"`
	MintableERC721Predicate       string `json:"MintableERC721Predicate"`
	MintableERC721PredicateProxy  string `json:"MintableERC721PredicateProxy"`
	MintableERC1155Predicate      string `json:"MintableERC1155Predicate"`
	MintableERC1155PredicateProxy string `json:"MintableERC1155PredicateProxy"`
	Tokens                        struct {
		DummyERC20           string `json:"DummyERC20"`
		DummyERC721          string `json:"DummyERC721"`
		DummyERC1155         string `json:"DummyERC1155"`
		DummyMintableERC20   string `json:"DummyMintableERC20"`
		DummyMintableERC721  string `json:"DummyMintableERC721"`
		DummyMintableERC1155 string `json:"DummyMintableERC1155"`
	} `json:"Tokens"`
}

type FxPortalRootContracts struct {
	FxRoot                      string `json:"FxRoot"`
	FxERC20RootTunnel           string `json:"FxERC20RootTunnel"`
	FxERC721RootTunnel          string `json:"FxERC721RootTunnel"`
	FxERC1155RootTunnel         string `json:"FxERC1155RootTunnel"`
	FxMintableERC20RootTunnel   string `json:"FxMintableERC20RootTunnel"`
	FxMintableERC721RootTunnel  string `json:"FxMintableERC721RootTunnel"`
	FxMintableERC1555RootTunnel string `json:"FxMintableERC1555RootTunnel"`
	Tokens                      struct {
		FxERC20Root   string `json:"FxERC20Root"`
		FxERC721Root  string `json:"FxERC721Root"`
		FxERC1155Root string `json:"FxERC1155Root"`
	} `json:"Tokens"`
}
type FxPortalChildContracts struct {
	FxChild                      string `json:"FxChild"`
	FxERC20ChildTunnel           string `json:"FxERC20ChildTunnel"`
	FxERC721ChildTunnel          string `json:"FxERC721ChildTunnel"`
	FxERC1155ChildTunnel         string `json:"FxERC1155ChildTunnel"`
	FxMintableERC20ChildTunnel   string `json:"FxMintableERC20ChildTunnel"`
	FxMintableERC721ChildTunnel  string `json:"FxMintableERC721ChildTunnel"`
	FxMintableERC1155ChildTunnel string `json:"FxMintableERC1155ChildTunnel"`
}

type GenesisChildContracts struct {
	BorValidatorSet string `json:"BorValidatorSet"`
	StateReceiver   string `json:"StateReceiver"`
}
