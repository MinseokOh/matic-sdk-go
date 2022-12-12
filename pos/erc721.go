package pos

import (
	"context"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type ERC721 struct {
	*BaseToken
}

func newERC721(client *Client, address common.Address, networkType types.NetworkType) *ERC721 {
	return &ERC721{
		BaseToken: newBaseToken(client, address, networkType, "erc721"),
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

	return common.Hash{}, nil
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

func (erc721 *ERC721) IsApproved() {
	erc721.Logger().Debug("IsApproved", log.Fields{})
}

func (erc721 *ERC721) IsApprovedAll() {
	erc721.Logger().Debug("IsApprovedAll", log.Fields{})
}

func (erc721 *ERC721) Withdraw(ctx context.Context, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Withdraw", log.Fields{})

	if err := erc721.checkForChild("Withdraw"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	return common.Hash{}, nil
}

func (erc721 *ERC721) WithdrawMany(ctx context.Context, tokenIds []*big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("WithdrawMany", log.Fields{})

	if err := erc721.checkForChild("WithdrawMany"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	if err := erc721.validateMany(tokenIds); err != nil {
		return common.Hash{}, err
	}

	return common.Hash{}, nil
}

func (erc721 *ERC721) Exit(ctx context.Context, txHash common.Hash, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Exit", log.Fields{})

	if err := erc721.checkForChild("Exit"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	return common.Hash{}, nil
}

func (erc721 *ERC721) ExitMany(ctx context.Context, txHash common.Hash, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("ExitMany", log.Fields{})

	if err := erc721.checkForChild("ExitMany"); err != nil {
		return common.Hash{}, err
	}

	if err := types.ValidateTxOption(txOption); err != nil {
		return common.Hash{}, err
	}

	return common.Hash{}, nil
}
