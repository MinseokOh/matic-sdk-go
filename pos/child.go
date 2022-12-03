package pos

import (
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

type ChildClient struct {
	*ethclient.Client
	rpc    *rpc.Client
	config types.ChildConfig
	logger *types.Logger
}

func NewChildClient(config types.POSClientConfig) (*ChildClient, error) {
	child := ChildClient{
		config: config.Child,
		logger: types.NewLogger("child", config.Debug),
	}
	var err error
	child.rpc, err = rpc.Dial(child.config.Rpc)
	if err != nil {
		return nil, err
	}
	child.Client = ethclient.NewClient(child.rpc)

	child.Logger().Debug("NewChildClient", log.Fields{
		"rpc" : child.config.Rpc,
	})

	return &child, nil
}

func (child *ChildClient) Rpc() *rpc.Client      { return child.rpc }
func (child *ChildClient) Logger() *types.Logger { return child.logger }
