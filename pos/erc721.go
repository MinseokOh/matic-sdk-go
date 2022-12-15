package pos

import (
	"context"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type ERC721 struct {
	*BaseToken
}

func newERC721(client *Client, address common.Address, networkType types.NetworkType) *ERC721 {
	return &ERC721{
		BaseToken: newBaseToken(client, address, networkType, types.ERC721),
	}
}

func (erc721 *ERC721) Approve(ctx context.Context, spender common.Address, tokenId *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Approve", log.Fields{
		"amount":   tokenId,
		"spender":  spender,
		"contract": erc721.address.String(),
	})

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	txHash, err := erc721.approve(ctx, spender, tokenId, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (erc721 *ERC721) ApproveAll(ctx context.Context, spender common.Address, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("ApproveAll", log.Fields{
		"spender":  spender,
		"contract": erc721.address.String(),
	})

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	if spender == (common.Address{}) {
		spender = erc721.PredicateAddress()
	}

	data, err := maticabi.ERC721.Pack("setApprovalForAll", spender, true)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(erc721.address, data, big.NewInt(0)).Build(ctx, erc721.getClient())
	if err != nil {
		return common.Hash{}, err
	}

	err = erc721.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (erc721 *ERC721) Deposit(ctx context.Context, tokenId *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Deposit", log.Fields{
		"tokenId":  tokenId,
		"contract": erc721.address.String(),
	})
	if err := erc721.checkForRoot("Deposit"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	depositData, err := maticabi.Deposit.Pack(tokenId)
	if err != nil {
		return common.Hash{}, err
	}

	txHash, err := erc721.deposit(ctx, depositData, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (erc721 *ERC721) DepositMany(ctx context.Context, tokenIds []*big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("DepositMany", log.Fields{
		"tokenIds": tokenIds,
		"contract": erc721.address.String(),
	})

	if err := erc721.checkForRoot("DepositMany"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	if err := erc721.validateMany(tokenIds); err != nil {
		return common.Hash{}, err
	}

	depositData, err := maticabi.DepositMany.Pack(tokenIds)
	if err != nil {
		return common.Hash{}, err
	}

	txHash, err := erc721.deposit(ctx, depositData, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (erc721 *ERC721) validateMany(tokenIds []*big.Int) error {
	if len(tokenIds) > 20 {
		return fmt.Errorf("can not process more than 20 tokens")
	}
	return nil
}

func (erc721 *ERC721) IsApproved(ctx context.Context, tokenId *big.Int) (bool, error) {
	erc721.Logger().Debug("IsApproved", log.Fields{
		"tokenId":  tokenId,
		"contract": erc721.address,
	})

	if err := erc721.checkForRoot("IsApproved"); err != nil {
		return false, err
	}

	getApprovedResp, err := utils.CallContract(ctx, erc721.getClient(), erc721.address, maticabi.ERC721,
		"getApproved",
		tokenId,
	)
	if err != nil {
		return false, err
	}

	if erc721.PredicateAddress() == getApprovedResp[0].(common.Address) {
		return true, nil
	}

	return false, nil
}

func (erc721 *ERC721) IsApprovedAll(ctx context.Context, address common.Address) (bool, error) {
	erc721.Logger().Debug("IsApprovedAll", log.Fields{
		"address":  address,
		"contract": erc721.address,
	})

	if err := erc721.checkForRoot("IsApprovedAll"); err != nil {
		return false, err
	}

	getApprovedResp, err := utils.CallContract(ctx, erc721.getClient(), erc721.address, maticabi.ERC721,
		"isApprovedForAll",
		address,
		erc721.PredicateAddress(),
	)
	if err != nil {
		return false, err
	}

	return getApprovedResp[0].(bool), nil
}

func (erc721 *ERC721) Withdraw(ctx context.Context, tokenId *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Withdraw", log.Fields{
		"tokenId":  tokenId,
		"contract": erc721.address,
	})

	if err := erc721.checkForChild("Withdraw"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.ERC721.Pack("withdraw", tokenId)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(erc721.address, data, big.NewInt(0)).Build(ctx, erc721.getClient())
	if err != nil {
		return common.Hash{}, err
	}

	err = erc721.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc721.Logger().Debug("Withdraw", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc721 *ERC721) WithdrawMany(ctx context.Context, tokenIds []*big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("WithdrawMany", log.Fields{
		"tokenIds": tokenIds,
		"contract": erc721.address,
	})

	if err := erc721.checkForChild("WithdrawMany"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	if err := erc721.validateMany(tokenIds); err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.ERC721.Pack("withdrawBatch", tokenIds)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(erc721.address, data, big.NewInt(0)).Build(ctx, erc721.getClient())
	if err != nil {
		return common.Hash{}, err
	}

	err = erc721.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	erc721.Logger().Debug("Withdraw", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (erc721 *ERC721) Exit(ctx context.Context, txHash common.Hash, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Exit", log.Fields{
		"txHash": txHash,
	})

	if err := erc721.checkForRoot("Exit"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	checkPointed, err := erc721.client.IsCheckPointed(ctx, txHash)
	if err != nil {
		return common.Hash{}, err
	}

	if !checkPointed {
		return common.Hash{}, fmt.Errorf("not checkpointed tx: %s", txHash.String())
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	payload, err := erc721.client.BuildPayloadForExit(ctx, txHash, types.ERC721Transfer, 0)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := erc721.exit(ctx, payload, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	erc721.Logger().Debug("Exit", log.Fields{
		"txHash": hash,
	})

	return hash, nil
}

func (erc721 *ERC721) ExitMany(ctx context.Context, txHash common.Hash, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("ExitMany", log.Fields{
		"txHash": txHash,
	})

	if err := erc721.checkForRoot("ExitMany"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	checkPointed, err := erc721.client.IsCheckPointed(ctx, txHash)
	if err != nil {
		return common.Hash{}, err
	}

	if !checkPointed {
		return common.Hash{}, fmt.Errorf("not checkpointed tx: %s", txHash.String())
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	payload, err := erc721.client.BuildPayloadForExit(ctx, txHash, types.ERC1155BatchTransfer, 0)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := erc721.exit(ctx, payload, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	erc721.Logger().Debug("ExitMany", log.Fields{
		"txHash": hash,
	})

	return hash, nil
}
