package utils

import (
	"github.com/MinseokOh/matic-sdk-go/types"
	maticabi "github.com/MinseokOh/matic-sdk-go/types/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	bigOne             = big.NewInt(1)
	bigTwo             = big.NewInt(2)
	checkPointInterval = big.NewInt(10000)
)

func FindRootBlockFromChild(client types.IClient, childBlockNumber *big.Int, rootChain common.Address) (*big.Int, error) {
	currentHeaderBlockResp, err := CallContract(client, rootChain, maticabi.RootChain,
		"currentHeaderBlock",
	)
	if err != nil {
		return nil, err
	}
	currentHeaderBlock := currentHeaderBlockResp[0].(*big.Int)

	// first checkpoint id = start * 10000
	start := bigOne
	// last checkpoint id = end * 10000
	end := new(big.Int).Div(currentHeaderBlock, checkPointInterval)

	// binary search on all the checkpoints to find the checkpoint that contains the childBlockNumber
	var ans *big.Int
	for start.Cmp(end) != 0 {
		mid := new(big.Int).Div(new(big.Int).Add(start, end), bigTwo)

		headerBlocksResp, err := CallContract(client, rootChain, maticabi.RootChain,
			"headerBlocks",
			new(big.Int).Mul(mid, checkPointInterval),
		)
		if err != nil {
			return nil, err
		}

		headerStart := headerBlocksResp[1].(*big.Int)
		headerEnd := headerBlocksResp[2].(*big.Int)

		if headerStart.Cmp(childBlockNumber) == -1 && childBlockNumber.Cmp(headerEnd) == -1 {
			// if childBlockNumber is between the upper and lower bounds of the headerBlock, we found our answer
			ans = mid
			break
		} else if headerStart.Cmp(childBlockNumber) == 1 {
			// childBlockNumber was checkpointed before this header
			end = mid.Sub(mid, bigOne)
		} else if headerStart.Cmp(childBlockNumber) == -1 {
			// childBlockNumber was checkpointed after this header
			start = mid.Add(mid, bigOne)
		}
	}
	if start.Cmp(end) == 0 {
		ans = start
	}

	return new(big.Int).Mul(ans, checkPointInterval), nil
}

