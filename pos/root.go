package pos

import (
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

func (root *RootClient) GetRootBlockInfo(txBlockNumber *big.Int) (types.RootBlockInfo, error) {
	root.Logger().Debug("GetRootBlockInfo",
		log.Fields{
			"txBlockNumber": txBlockNumber,
		},
	)

	rootBlockNumber, err := utils.FindRootBlockFromChild(root, txBlockNumber, root.config.RootChain)
	if err != nil {
		return types.RootBlockInfo{}, err
	}

	headerBlocksResp, err := utils.CallContract(root, root.config.RootChain, maticabi.RootChain,
		"headerBlocks",
		rootBlockNumber,
	)

	headerBlock := types.RootBlockInfo{
		HeaderBlockNumber: rootBlockNumber,
		Start:             headerBlocksResp[1].(*big.Int),
		End:               headerBlocksResp[2].(*big.Int),
	}

	root.Logger().Info("RootBlockInfo",
		log.Fields{
			"headerBlock": headerBlock,
		},
	)

	return headerBlock, nil
}
