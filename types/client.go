package types

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type IClient interface {
	Rpc() *rpc.Client
	Logger() *Logger

	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}
