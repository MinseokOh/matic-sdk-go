package pos

import (
	"context"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type Client struct {
	config types.POSClientConfig
	logger *types.Logger
	Child  *ChildClient
	Root   *RootClient
}

func NewClient(config types.POSClientConfig) (*Client, error) {
	client := Client{
		config: config,
		logger: types.NewLogger("pos", config.Debug),
	}
	rootClient, err := NewRootClient(config)
	if err != nil {
		return nil, err
	}
	client.Root = rootClient

	childClient, err := NewChildClient(config)
	if err != nil {
		return nil, err
	}
	client.Child = childClient

	return &client, nil
}

func (client *Client) Logger() *types.Logger { return client.logger }

func (client *Client) ERC20(address common.Address, networkType types.NetworkType) *ERC20 {
	return newERC20(client, address, networkType)
}

func (client *Client) ERC721(address common.Address, networkType types.NetworkType) *ERC721 {
	return newERC721(client, address, networkType)
}

func (client *Client) DepositEtherFor(ctx context.Context, amount *big.Int, txOption *types.TxOption) (common.Hash, error) {
	client.Logger().Debug("DepositEtherFor", log.Fields{
		"amount": amount,
	})

	if txOption == nil {
		return common.Hash{}, types.EmptyTxOption
	}

	if err := txOption.Validate(); err != nil {
		return common.Hash{}, err
	}

	data, err := maticabi.RootChainManager.Pack("depositEtherFor", txOption.From())
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(client.config.Root.RootChainManager, data, amount).Build(ctx, client.Root)
	if err != nil {
		return common.Hash{}, err
	}

	err = client.Root.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	client.Logger().Debug("DepositEtherFor", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (client *Client) ExitEther(ctx context.Context, txHash common.Hash, txOption *types.TxOption) (common.Hash, error) {
	return client.ERC20(common.Address{}, types.Root).Exit(ctx, txHash, txOption)
}

func (client *Client) BuildPayloadForExit(ctx context.Context, txHash common.Hash, eventSignature string, index int) ([]byte, error) {
	client.Logger().Debug("BuildPayloadForExit", log.Fields{
		"txHash": txHash,
	})

	client.Logger().Debug("TransactionReceipt", log.Fields{
		"txHash": txHash,
	})
	receipt, err := client.Child.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}

	client.Logger().Debug("BlockByNumber", log.Fields{
		"blockNumber": receipt.BlockNumber,
	})
	block, err := client.Child.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return nil, err
	}

	client.Logger().Debug("GetRootBlockInfo", nil)
	blockInfo, err := client.Root.GetRootBlockInfo(ctx, block.Number())
	if err != nil {
		return nil, err
	}

	client.Logger().Debug("BuildBlockProof", nil)
	blockProof, err := utils.BuildBlockProof(ctx, client.Child, receipt.BlockNumber, blockInfo.Start, blockInfo.End)
	if err != nil {
		return nil, err
	}

	client.Logger().Debug("GetReceiptProof", nil)
	path, receiptProof, err := utils.GetReceiptProof(ctx, client.Child, receipt, block)
	if err != nil {
		return nil, err
	}

	rawReceipt, err := receipt.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var logIndex uint64
	if index > 0 {
		// when token index is not 0
		client.Logger().Debug("GetLogIndices", nil)
		logIndices, err := utils.GetAllLogIndices(eventSignature, receipt)
		if err != nil {
			return nil, err
		}

		if index >= len(logIndices) {
			return nil, fmt.Errorf("index is grater than the number of tokens in this transaction")
		}
		logIndex = logIndices[index]
	} else {
		// when token index is 0
		client.Logger().Debug("GetLogIndex", nil)
		logIndex = utils.GetLogIndex(eventSignature, receipt)
	}

	payload, err := rlp.EncodeToBytes([]interface{}{
		// headerNumber - Checkpoint header block number containing the burn tx
		blockInfo.HeaderBlockNumber.Uint64(),
		// blockProof - Proof that the block header (in the child chain) is a leaf in the submitted merkle root
		blockProof,
		// blockNumber - Block number containing the burn tx on child chain
		block.Number().Uint64(),
		// blockTime - Burn tx block time
		block.Time(),
		// txRoot - Transactions root of block
		block.TxHash(),
		// receiptRoot - Receipts root of block
		block.ReceiptHash(),
		// receipt - Receipt of the burn transaction
		rawReceipt,
		// receiptProof - Merkle proof of the burn receipt
		receiptProof,
		// branchMask - 32 bits denoting the path of receipt in merkle patricia tree
		append([]byte{0}, path...),
		// receiptLogIndex - Log Index to read from the receipt
		logIndex,
	})
	if err != nil {
		return nil, err
	}

	client.logger.Debug("ExitPayload", log.Fields{
		"payload": hexutil.Encode(payload),
	})
	return payload, nil
}

func (client *Client) IsCheckPointed(ctx context.Context, txHash common.Hash) (bool, error) {
	client.Logger().Debug("IsCheckPointed", log.Fields{
		"txHash": txHash,
	})
	var err error
	lastChildBlock, err := client.Root.GetLastChildBlock(ctx)
	if err != nil {
		return false, err
	}

	client.Logger().Debug("TransactionReceipt", log.Fields{
		"txHash": txHash,
	})
	receipt, err := client.Child.TransactionReceipt(ctx, txHash)
	if err != nil {
		return false, err
	}

	if lastChildBlock.Cmp(receipt.BlockNumber) == 1 {
		client.Logger().Debug("IsCheckPointed", log.Fields{
			"checkPointed": true,
			"txHash":       txHash,
		})
		return true, nil
	} else {
		client.Logger().Debug("IsCheckPointed", log.Fields{
			"checkPointed": false,
			"txHash":       txHash,
		})
		return false, nil
	}
}
