package utils

import (
	"context"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func CallContract(client types.IClient, to common.Address, abi abi.ABI, method string, args ...interface{}) ([]interface{}, error) {
	client.Logger().Debug("CallContract", log.Fields{
		"method": method,
		"args":   args,
	})

	m, ok := abi.Methods[method]
	if !ok {
		return nil, fmt.Errorf("not found method: %s", method)
	}

	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From: common.Address{},
		To:   &to,
		Data: data,
	}

	b, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	return m.Outputs.UnpackValues(b)
}
