package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type RootClient struct {
	*ethclient.Client
	config types.RootConfig
	logger *types.Logger
	rpc    *rpc.Client
}

func NewRootClient(config types.POSClientConfig) (*RootClient, error) {
	root := RootClient{
		config: config.Root,
		logger: types.NewLogger("root", config.Debug),
	}
	var err error
	root.rpc, err = rpc.Dial(root.config.Rpc)
	if err != nil {
		return nil, err
	}
	root.Client = ethclient.NewClient(root.rpc)

	root.Logger().Debug("NewRootClient", log.Fields{
		"rpc": root.config.Rpc,
	})

	return &root, nil
}

func (root *RootClient) Rpc() *rpc.Client      { return root.rpc }
func (root *RootClient) Logger() *types.Logger { return root.logger }

func (root *RootClient) GetRootBlockInfo(ctx context.Context, txBlockNumber *big.Int) (types.RootBlockInfo, error) {
	root.Logger().Debug("GetRootBlockInfo",
		log.Fields{
			"txBlockNumber": txBlockNumber,
		},
	)

	rootBlockNumber, err := utils.FindRootBlockFromChild(ctx, root, txBlockNumber, root.config.RootChain)
	if err != nil {
		return types.RootBlockInfo{}, err
	}

	headerBlocksResp, err := utils.CallContract(ctx, root, root.config.RootChain, maticabi.RootChain,
		"headerBlocks",
		rootBlockNumber,
	)
	if err != nil {
		return types.RootBlockInfo{}, err
	}

	headerBlock := types.RootBlockInfo{
		HeaderBlockNumber: rootBlockNumber,
		Start:             headerBlocksResp[1].(*big.Int),
		End:               headerBlocksResp[2].(*big.Int),
	}

	root.Logger().Debug("RootBlockInfo",
		log.Fields{
			"headerBlock": headerBlock,
		},
	)

	return headerBlock, nil
}

func (root *RootClient) GetLastChildBlock(ctx context.Context) (*big.Int, error) {
	root.Logger().Debug("GetLastChildBlock", nil)

	getLastChildBlockResp, err := utils.CallContract(ctx, root, root.config.RootChain, maticabi.RootChain,
		"getLastChildBlock",
	)
	if err != nil {
		return nil, err
	}
	lastChildBlock := getLastChildBlockResp[0].(*big.Int)

	root.Logger().Debug("GetLastChildBlock",
		log.Fields{
			"blockNumber": lastChildBlock,
		},
	)
	return lastChildBlock, nil
}
