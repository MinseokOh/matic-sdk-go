package pos

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ether "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type ERC20 struct {
	client      *Client
	config      types.POSClientConfig
	networkType types.NetworkType
	logger      *types.Logger
	address     common.Address
}

func newERC20(client *Client, address common.Address, networkType types.NetworkType) *ERC20 {
	return &ERC20{
		client:      client,
		config:      client.config,
		address:     address,
		networkType: networkType,
		logger:      types.NewLogger("erc20", client.config.Debug),
	}
}
func (erc20 *ERC20) Logger() *types.Logger { return erc20.logger }

func (erc20 *ERC20) Approve(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	erc20.Logger().Debug("Approve", log.Fields{
		"amount":   amount,
		"contract": erc20.address.String(),
	})
	rootClient := erc20.client.Root

	address := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	chainId, err := rootClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := rootClient.PendingNonceAt(ctx, address)
	if err != nil {
		return common.Hash{}, err
	}

	gasTipCap, err := rootClient.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.ERC20.Pack("approve", erc20.config.Root.ERC20Predicate, amount)
	if err != nil {
		return common.Hash{}, err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		GasFeeCap: gasTipCap,
		Gas:       6e4,
		Nonce:     nonce,
		To:        &erc20.address,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	err = rootClient.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Info("Approve", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc20 *ERC20) ApproveMax(ctx context.Context, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	erc20.Logger().Debug("ApproveMax", log.Fields{
		"contract": erc20.address.String(),
	})
	rootClient := erc20.client.Root

	amount, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	address := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	chainId, err := rootClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := rootClient.PendingNonceAt(ctx, address)
	if err != nil {
		return common.Hash{}, err
	}

	gasTipCap, err := rootClient.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.ERC20.Pack("approve", erc20.config.Root.ERC20Predicate, amount)
	if err != nil {
		return common.Hash{}, err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		GasFeeCap: gasTipCap,
		Gas:       6e4,
		Nonce:     nonce,
		To:        &erc20.address,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	err = rootClient.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Info("ApproveMax", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc20 *ERC20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	balanceOfResp, err := utils.CallContract(ctx, erc20.getClient(), erc20.address, maticabi.ERC20, "allowance", owner, spender)
	if err != nil {
		return nil, err
	}
	allowance := balanceOfResp[0].(*big.Int)

	erc20.Logger().Info("Allowance", log.Fields{
		"allowance": allowance,
	})

	return allowance, nil
}

func (erc20 *ERC20) DepositFor(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	erc20.Logger().Debug("DepositFor", log.Fields{
		"amount":   amount,
		"contract": erc20.address.String(),
	})
	if err := erc20.checkForRoot("DepositFor"); err != nil {
		return common.Hash{}, err
	}
	rootClient := erc20.getClient()

	address := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	chainId, err := rootClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := rootClient.PendingNonceAt(ctx, address)
	if err != nil {
		return common.Hash{}, err
	}

	gasTipCap, err := rootClient.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return common.Hash{}, err
	}

	deposit := abi.Arguments{
		{Type: uint256Ty},
	}

	depositData, err := deposit.Pack(amount)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.RootChainManager.Pack("depositFor", address, erc20.address, depositData)
	if err != nil {
		return common.Hash{}, err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		GasFeeCap: gasTipCap,
		Gas:       11e4,
		Nonce:     nonce,
		To:        &erc20.config.Root.RootChainManager,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	err = rootClient.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Info("DepositFor", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc20 *ERC20) Withdraw(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	erc20.Logger().Debug("Withdraw", log.Fields{
		"amount":   amount,
		"contract": erc20.address.String(),
	})
	if err := erc20.checkForChild("Withdraw"); err != nil {
		return common.Hash{}, err
	}
	childClient := erc20.getClient()

	address := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	chainId, err := childClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := childClient.PendingNonceAt(ctx, address)
	if err != nil {
		return common.Hash{}, err
	}

	gasTipCap, err := childClient.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.ERC20.Pack("withdraw", amount)
	if err != nil {
		return common.Hash{}, err
	}

	value := big.NewInt(0)
	if erc20.address == common.HexToAddress(types.MaticAddress) {
		value = amount
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		GasFeeCap: gasTipCap,
		Gas:       7e4,
		Nonce:     nonce,
		To:        &erc20.address,
		Value:     value,
		Data:      data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	err = childClient.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Info("Withdraw", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc20 *ERC20) Exit(ctx context.Context, txHash common.Hash, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	erc20.Logger().Debug("Exit", log.Fields{
		"txHash":   txHash.String(),
		"contract": erc20.address.String(),
	})
	if err := erc20.checkForRoot("Exit"); err != nil {
		return common.Hash{}, err
	}

	checkPointed, err := erc20.client.IsCheckPointed(ctx, txHash)
	if err != nil {
		return common.Hash{}, err
	}

	if !checkPointed {
		return common.Hash{}, fmt.Errorf("not checkpointed tx: %s", txHash.String())
	}

	rootClient := erc20.getClient()

	address := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	chainId, err := rootClient.ChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := rootClient.PendingNonceAt(ctx, address)
	if err != nil {
		return common.Hash{}, err
	}

	gasTipCap, err := rootClient.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	payload, err := erc20.client.BuildPayloadForExit(ctx, txHash, types.ERC20Transfer)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.RootChainManager.Pack("exit", payload)
	if err != nil {
		return common.Hash{}, err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		GasFeeCap: gasTipCap,
		Gas:       1e6,
		Nonce:     nonce,
		To:        &erc20.config.Root.RootChainManager,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	err = rootClient.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Info("Exit", log.Fields{
		"txHash": tx.Hash(),
	})

	return tx.Hash(), err
}

func (erc20 *ERC20) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	balanceOfResp, err := utils.CallContract(ctx, erc20.getClient(), erc20.address, maticabi.ERC20, "balanceOf", address)
	if err != nil {
		return nil, err
	}
	balance := balanceOfResp[0].(*big.Int)

	erc20.Logger().Info("BalanceOf", log.Fields{
		"balance": balance,
	})

	return balance, nil
}

func (erc20 *ERC20) getClient() types.IClient {
	if erc20.networkType == types.Root {
		return erc20.client.Root
	} else {
		return erc20.client.Child
	}
}

func (erc20 *ERC20) checkForRoot(method string) error {
	if erc20.networkType != types.Root {
		return fmt.Errorf("allowed on root %s", method)
	}
	return nil
}

func (erc20 *ERC20) checkForChild(method string) error {
	if erc20.networkType != types.Child {
		return fmt.Errorf("allowed on child %s", method)
	}
	return nil
}
