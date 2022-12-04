package utils

import (
	"github.com/MinseokOh/matic-sdk-go/types"
	ether "github.com/ethereum/go-ethereum/core/types"
	"strings"
)

const emptyTopic = "0x0000000000000000000000000000000000000000000000000000000000000000"

func GetAllLogIndices(logEventSig string, receipt *ether.Receipt) ([]int, error) {
	var logIndices []int
	for i, log := range receipt.Logs {
		switch logEventSig {
		case types.ERC20Transfer, types.ERC721TransferWithMetadata:
			if len(log.Topics) < 2 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[2].String()) == emptyTopic {
				logIndices = append(logIndices, i)
			}
		case types.ERC1155Transfer, types.ERC1155BatchTransfer:
			if len(log.Topics) < 3 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[3].String()) == emptyTopic {
				logIndices = append(logIndices, i)
			}
		case types.ERC721BatchTransfer:
			if len(log.Topics) < 2 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(types.ERC721Transfer) &&
				strings.ToLower(log.Topics[2].String()) == emptyTopic {
				logIndices = append(logIndices, i)
			}
		}
	}

	return logIndices, nil
}

func GetLogIndex(logEventSig string, receipt *ether.Receipt) uint64 {
	logIndex := uint64(0)
	for i, log := range receipt.Logs {
		switch logEventSig {
		case types.ERC20Transfer, types.ERC721TransferWithMetadata:
			if len(log.Topics) < 2 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[2].String()) == emptyTopic {
				logIndex = uint64(i)
			}
		case types.ERC1155Transfer, types.ERC1155BatchTransfer:
			if len(log.Topics) < 3 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[3].String()) == emptyTopic {
				logIndex = uint64(i)
			}
		default:
			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) {
				logIndex = uint64(i)
			}
		}
	}

	return logIndex
}
