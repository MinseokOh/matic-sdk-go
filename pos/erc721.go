package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func (erc721 *ERC721) Approve() {
	erc721.Logger().Debug("Approve", log.Fields{})
}

func (erc721 *ERC721) ApproveAll() {
	erc721.Logger().Debug("ApproveAll", log.Fields{})
}

func (erc721 *ERC721) Deposit(ctx context.Context, tokenId *big.Int, txOption *types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("Deposit", log.Fields{
		"tokenId": tokenId,
	})

	if err := erc721.checkForRoot("Deposit"); err != nil {
		return common.Hash{}, err
	}

	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return common.Hash{}, err
	}

	deposit := abi.Arguments{
		{Type: uint256Ty},
	}

	depositData, err := deposit.Pack(tokenId)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.RootChainManager.Pack("depositFor", txOption.From(), erc721.address, depositData)
	if err != nil {
		return common.Hash{}, err
	}

	txHash, err := erc721.deposit(ctx, data, txOption)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (erc721 *ERC721) DepositMany(ctx context.Context, tokenIds []*big.Int, txOption types.TxOption) (common.Hash, error) {
	erc721.Logger().Debug("DepositMany", log.Fields{
		"tokenIds": tokenIds,
	})

	if err := erc721.checkForRoot("DepositMany"); err != nil {
		return common.Hash{}, err
	}

	return common.Hash{}, nil
}

func (erc721 *ERC721) IsApproved() {
	erc721.Logger().Debug("IsApproved", log.Fields{})
}

func (erc721 *ERC721) IsApprovedAll() {
	erc721.Logger().Debug("IsApprovedAll", log.Fields{})
}

func (erc721 *ERC721) Withdraw() {
	erc721.Logger().Debug("Withdraw", log.Fields{})
}

func (erc721 *ERC721) WithdrawMany() {
	erc721.Logger().Debug("WithdrawMany", log.Fields{})
}

func (erc721 *ERC721) Exit() {
	erc721.Logger().Debug("Exit", log.Fields{})
}

func (erc721 *ERC721) ExitMany() {
	erc721.Logger().Debug("ExitMany", log.Fields{})
}
