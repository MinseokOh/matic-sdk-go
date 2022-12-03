package utils

import (
	sdk "github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

const emptyTopic = "0x0000000000000000000000000000000000000000000000000000000000000000"

func GetAllLogIndices(logEventSig string, receipt *types.Receipt) ([]int, error) {
	var logIndices []int
	for i, log := range receipt.Logs {
		switch logEventSig {
		case sdk.ERC20Transfer, sdk.ERC721TransferWithMetadata:
			if len(log.Topics) < 2 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[2].String()) == emptyTopic {
				logIndices = append(logIndices, i)
			}
		case sdk.ERC1155Transfer, sdk.ERC1155BatchTransfer:
			if len(log.Topics) < 3 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[3].String()) == emptyTopic {
				logIndices = append(logIndices, i)
			}
		case sdk.ERC721BatchTransfer:
			if len(log.Topics) < 2 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(sdk.ERC721Transfer) &&
				strings.ToLower(log.Topics[2].String()) == emptyTopic {
				logIndices = append(logIndices, i)
			}
		}
	}

	return logIndices, nil
}

func GetLogIndex(logEventSig string, receipt *types.Receipt) uint64 {
	logIndex := uint64(0)
	for i, log := range receipt.Logs {
		switch logEventSig {
		case sdk.ERC20Transfer, sdk.ERC721TransferWithMetadata:
			if len(log.Topics) < 2 {
				continue
			}

			if strings.ToLower(log.Topics[0].String()) == strings.ToLower(logEventSig) &&
				strings.ToLower(log.Topics[2].String()) == emptyTopic {
				logIndex = uint64(i)
			}
		case sdk.ERC1155Transfer, sdk.ERC1155BatchTransfer:
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
