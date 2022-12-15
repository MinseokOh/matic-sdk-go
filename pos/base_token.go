package pos

import (
	"context"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/MinseokOh/matic-sdk-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"math/big"
)

type BaseToken struct {
	client           *Client
	config           types.POSClientConfig
	networkType      types.NetworkType
	tokenType        types.TokenType
	logger           *types.Logger
	address          common.Address
	predicateAddress common.Address
}

func newBaseToken(client *Client, address common.Address, networkType types.NetworkType, tokenType types.TokenType) *BaseToken {
	return &BaseToken{
		client:      client,
		config:      client.config,
		address:     address,
		networkType: networkType,
		logger:      types.NewLogger(tokenType.String(), client.config.Debug),
	}
}

func (token *BaseToken) Logger() *types.Logger { return token.logger }

func (token *BaseToken) approve(ctx context.Context, spender common.Address, value *big.Int, txOption *types.TxOption) (common.Hash, error) {
	if spender == (common.Address{}) {
		spender = token.PredicateAddress()
	}

	data, err := maticabi.ERC20.Pack("approve", spender, value)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(token.address, data, big.NewInt(0)).Build(ctx, token.getClient())
	if err != nil {
		return common.Hash{}, err
	}

	err = token.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (token *BaseToken) deposit(ctx context.Context, depositData []byte, txOption *types.TxOption) (common.Hash, error) {
	token.logger.Debug("depositFor", log.Fields{
		"from": txOption.From(),
		"to":   token.address,
		"data": hexutil.Encode(depositData),
	})

	data, err := maticabi.RootChainManager.Pack("depositFor", txOption.From(), token.address, depositData)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(token.config.Root.RootChainManager, data, big.NewInt(0)).Build(ctx, token.getClient())
	if err != nil {
		return common.Hash{}, err
	}

	err = token.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (token *BaseToken) exit(ctx context.Context, payload []byte, txOption *types.TxOption) (common.Hash, error) {
	token.logger.Debug("exit", log.Fields{
		"payload": hexutil.Encode(payload),
	})
	data, err := maticabi.RootChainManager.Pack("exit", payload)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := txOption.SetTxData(token.config.Root.RootChainManager, data, big.NewInt(0)).Build(ctx, token.getClient())
	if err != nil {
		return common.Hash{}, err
	}
	
	err = token.getClient().SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (token *BaseToken) PredicateAddress() common.Address {
	if err := token.checkForRoot("PredicateAddress"); err != nil {
		token.Logger().Error("PredicateAddress", log.Fields{
			"msg": err,
		})
		return common.Address{}
	}

	if token.predicateAddress != (common.Address{}) {
		token.Logger().Debug("PredicateAddress", log.Fields{
			"address": token.predicateAddress,
		})

		return token.predicateAddress
	}

	tokenTypeResp, err := utils.CallContract(context.Background(), token.client.Root, token.config.Root.RootChainManager, maticabi.RootChainManager,
		"tokenToType",
		token.address,
	)
	if err != nil {
		token.Logger().Error("tokenToType", log.Fields{
			"error": err.Error(),
		})
		return common.Address{}
	}

	typeToPredicateResp, err := utils.CallContract(context.Background(), token.client.Root, token.config.Root.RootChainManager, maticabi.RootChainManager,
		"typeToPredicate",
		tokenTypeResp[0],
	)
	if err != nil {
		token.Logger().Error("typeToPredicate", log.Fields{
			"error": err.Error(),
		})
		return common.Address{}
	}
	token.predicateAddress = typeToPredicateResp[0].(common.Address)

	token.Logger().Debug("PredicateAddress", log.Fields{
		"address": token.predicateAddress,
	})
	return token.predicateAddress
}

func (token *BaseToken) getClient() types.IClient {
	if token.networkType == types.Root {
		return token.client.Root
	} else {
		return token.client.Child
	}
}

func (token *BaseToken) checkForRoot(method string) error {
	if token.networkType != types.Root {
		return fmt.Errorf("allowed on root %s", method)
	}
	return nil
}

func (token *BaseToken) checkForChild(method string) error {
	if token.networkType != types.Child {
		return fmt.Errorf("allowed on child %s", method)
	}
	return nil
}
