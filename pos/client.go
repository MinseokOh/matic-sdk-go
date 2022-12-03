package pos

import (
	"context"
	"crypto/ecdsa"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	childCient, err := NewChildClient(config)
	if err != nil {
		return nil, err
	}
	client.Child = childCient

	return &client, nil
}

func (client *Client) Logger() *types.Logger { return client.logger }

func (client *Client) ERC20(address common.Address) *ERC20 {
	return newERC20(address, client)
}

func (erc20 *ERC20) DepositEther(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) {

}


func (client *Client) BuildPayloadForExit(ctx context.Context, txHash common.Hash, eventSignature string) ([]byte, error) {
	client.Logger().Info("BuildPayloadForExit", log.Fields{
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
	blockProof, err := utils.BuildBlockProof(client.Child, receipt.BlockNumber, blockInfo.Start, blockInfo.End)
	if err != nil {
		return nil, err
	}

	client.Logger().Debug("GetReceiptProof", nil)
	path, receiptProof, err := utils.GetReceiptProof(client.Child, receipt, block)
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
		blockInfo.HeaderBlockNumber.Uint64(),
		blockProof,
		block.Number().Uint64(),
		block.Time(),
		block.TxHash(),
		block.ReceiptHash(),
		rawReceipt,
		receiptProof,
		append([]byte{0}, path...),
		index,
	})
	if err != nil {
		return nil, err
	}

	client.logger.Info("ExitPayload", log.Fields{
		"payload" : hexutil.Encode(payload),
	})

	return payload, nil
}
