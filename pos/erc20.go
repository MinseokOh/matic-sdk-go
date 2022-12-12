package pos

import (
	"context"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type ERC20 struct {
	*BaseToken
}

func newERC20(client *Client, address common.Address, networkType types.NetworkType) *ERC20 {
	return &ERC20{
		BaseToken: newBaseToken(client, address, networkType, "erc20"),
	}
}

// Approve : approve to spender, when spender is zero address, approve to predicate address
func (erc20 *ERC20) Approve(ctx context.Context, spender common.Address, amount *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc20.Logger().Debug("Approve", log.Fields{
		"amount":   amount,
		"contract": erc20.address.String(),
	})
	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	txHash, err := erc20.approve(ctx, spender, amount, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Debug("Approve", log.Fields{
		"txHash": txHash,
	})
	return txHash, nil
}

// ApproveMax : approve max to spender, when spender is zero address, approve to predicate address
func (erc20 *ERC20) ApproveMax(ctx context.Context, spender common.Address, txOption *types.TxOption) (common.Hash, error) {
	erc20.Logger().Debug("ApproveMax", log.Fields{
		"contract": erc20.address.String(),
	})
	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	amount, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	txHash, err := erc20.approve(ctx, spender, amount, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Debug("Approve", log.Fields{
		"txHash": txHash,
	})
	return txHash, nil
}

func (erc20 *ERC20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	balanceOfResp, err := utils.CallContract(ctx, erc20.getClient(), erc20.address, maticabi.ERC20, "allowance", owner, spender)
	if err != nil {
		return nil, err
	}
	allowance := balanceOfResp[0].(*big.Int)

	erc20.Logger().Debug("Allowance", log.Fields{
		"allowance": allowance,
	})

	return allowance, nil
}

func (erc20 *ERC20) Deposit(ctx context.Context, amount *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc20.Logger().Debug("Deposit", log.Fields{
		"amount":   amount,
		"contract": erc20.address.String(),
	})
	if err := erc20.checkForRoot("DepositFor"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	client := erc20.getClient()

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

	data, err := maticabi.RootChainManager.Pack("depositFor", txOption.From(), erc20.address, depositData)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(erc20.config.Root.RootChainManager, data, big.NewInt(0)).Build(ctx, client)
	if err != nil {
		return common.Hash{}, err
	}

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Debug("DepositFor", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc20 *ERC20) Withdraw(ctx context.Context, amount *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc20.Logger().Debug("Withdraw", log.Fields{
		"amount":   amount,
		"contract": erc20.address.String(),
	})
	if err := erc20.checkForChild("Withdraw"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
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

	tx, err := txOption.SetTxData(erc20.address, data, value).Build(ctx, erc20.getClient())
	if err != nil {
		return common.Hash{}, err
	}

	err = erc20.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Debug("Withdraw", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc20 *ERC20) Exit(ctx context.Context, txHash common.Hash, txOption *types.TxOption) (common.Hash, error) {
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

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	payload, err := erc20.client.BuildPayloadForExit(ctx, txHash, types.ERC20Transfer, 0)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := erc20.exit(ctx, payload, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	erc20.Logger().Debug("Exit", log.Fields{
		"txHash": hash,
	})

	return hash, err
}

func (erc20 *ERC20) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	balanceOfResp, err := utils.CallContract(ctx, erc20.getClient(), erc20.address, maticabi.ERC20, "balanceOf", address)
	if err != nil {
		return nil, err
	}
	balance := balanceOfResp[0].(*big.Int)

	erc20.Logger().Debug("BalanceOf", log.Fields{
		"balance": balance,
	})

	return balance, nil
}
