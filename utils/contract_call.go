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

func CallContract(ctx context.Context, client types.IClient, to common.Address, abi abi.ABI, method string, args ...interface{}) ([]interface{}, error) {
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

	b, err := client.CallContract(ctx, callMsg, nil)
	if err != nil {
		return nil, err
	}

	unpacked, err := m.Outputs.UnpackValues(b)
	if err != nil {
		return nil, err
	}

	client.Logger().Debug("CallContract", log.Fields{
		"method":   method,
		"unpacked": unpacked,
	})

	return unpacked, nil
}
