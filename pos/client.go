package pos

import (
	"context"
	"crypto/ecdsa"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ether "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type Client struct {
	ctx    context.Context
	config types.POSClientConfig
	logger *types.Logger
	Child  *ChildClient
	Root   *RootClient
}

func NewClient(config types.POSClientConfig) (*Client, error) {
	client := Client{
		ctx:    context.Background(),
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

func (client *Client) DepositEtherFor(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) (common.Hash, error) {
	client.Logger().Debug("DepositEtherFor", log.Fields{
		"amount": amount,
	})
	rootClient := client.Root

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

	data, err := maticabi.RootChainManager.Pack("depositEtherFor", address)
	if err != nil {
		return common.Hash{}, err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		GasFeeCap: gasTipCap,
		Gas:       1e5,
		Nonce:     nonce,
		To:        &client.config.Root.RootChainManager,
		Value:     amount,
		Data:      data,
	})
	if err != nil {
		return common.Hash{}, err
	}

	err = rootClient.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	client.Logger().Debug("DepositEtherFor", log.Fields{
		"txHash": tx.Hash(),
	})
	return tx.Hash(), nil
}

func (client *Client) BuildPayloadForExit(ctx context.Context, txHash common.Hash, eventSignature string) ([]byte, error) {
	client.Logger().Debug("BuildPayloadForExit", log.Fields{
		"txHash": txHash.String(),
	})

	client.Logger().Debug("TransactionReceipt", log.Fields{
		"txHash": txHash.String(),
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
	blockInfo, err := client.Root.GetRootBlockInfo(block.Number())
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

	client.Logger().Debug("GetLogIndex", nil)
	index := utils.GetLogIndex(eventSignature, receipt)

	rawReceipt, err := receipt.MarshalBinary()
	if err != nil {
		return nil, err
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
		index,
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
		"txHash": txHash.String(),
	})

	lastChildBlock, err := client.Root.GetLastChildBlock()
	if err != nil {
		return false, err
	}

	client.Logger().Debug("TransactionReceipt", log.Fields{
		"txHash": txHash.String(),
	})
	receipt, err := client.Child.TransactionReceipt(ctx, txHash)
	if err != nil {
		return false, err
	}

	if lastChildBlock.Cmp(receipt.BlockNumber) == 1 {
		client.Logger().Debug("IsCheckPointed", log.Fields{
			"checkPointed": true,
			"txHash":       txHash.String(),
		})
		return true, nil
	} else {
		client.Logger().Debug("IsCheckPointed", log.Fields{
			"checkPointed": false,
			"txHash":       txHash.String(),
		})
		return false, nil
	}
}
