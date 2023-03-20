package utils

import (
	"encoding/json"
	"github.com/MinseokOh/matic-sdk-go/types"
)

func GetContractByNetwork(network types.Network) types.Contract {
	var response string
	var err error
	switch network {
	case types.TestNet:
		response, err = Get(types.TestNetContractURL, nil, nil)
	case types.MainNet:
		response, err = Get(types.MainNetContractURL, nil, nil)
	}
	if err != nil {
		return types.Contract{}
	}

	var contract types.Contract
	if err = json.Unmarshal([]byte(response), &contract); err != nil {
		return types.Contract{}
	}

	return contract
}
