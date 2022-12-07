package utils

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/MinseokOh/merkle-patricia-trie/trie"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ether "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
)

var (
	bytes32Ty, _ = abi.NewType("bytes32", "", nil)
	argsZeroHash = abi.Arguments{
		{Type: bytes32Ty},
		{Type: bytes32Ty},
	}
)

func BuildBlockProof(ctx context.Context, client types.IClient, txBlockNumber, startBlock, endBlock *big.Int) ([]byte, error) {
	client.Logger().Debug("BuildBlockProof", log.Fields{
		"txBlockNumber": txBlockNumber,
		"start":         startBlock,
		"end":           endBlock,
	})
	proof, err := getFastMerkleProof(ctx, client, txBlockNumber, startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	var buf []byte
	for i := 0; i < len(proof); i++ {
		buf = append(buf, proof[i].Bytes()...)
	}

	client.Logger().Debug("BlockProof", log.Fields{
		"proof": hexutil.Encode(buf),
	})
	return buf, nil
}

func getFastMerkleProof(ctx context.Context, client types.IClient, txBlockNumber, startBlock, endBlock *big.Int) ([]common.Hash, error) {
	start := startBlock.Int64()
	end := endBlock.Int64()
	blockNumber := txBlockNumber.Int64()

	merkleTreeDepth := int64(math.Ceil(math.Log2(float64(end - start - 1))))

	proof := make([]common.Hash, 0)

	offset := start
	targetIndex := blockNumber - offset
	leftBound := int64(0)
	rightBound := end - offset

	//fmt.Println("Searching for", targetIndex)
	for depth := int64(0); depth < merkleTreeDepth; depth++ {
		nLeaves := int64(math.Pow(2, float64(merkleTreeDepth-depth)))
		pivotLeaf := leftBound + nLeaves/2 - 1

		if targetIndex > pivotLeaf {
			newLeftBound := pivotLeaf + 1
			subTreeMerkleRoot, err := queryRootHash(ctx, client, offset+leftBound, offset+pivotLeaf)
			if err != nil {
				return nil, err
			}

			proof = append(proof, subTreeMerkleRoot)
			leftBound = newLeftBound
		} else {
			newRightBound := int64(math.Min(float64(rightBound), float64(pivotLeaf)))

			expectedHeight := merkleTreeDepth - (depth + 1)
			if rightBound <= pivotLeaf {
				subTreeMerkleRoot := recursiveZeroHash(expectedHeight)
				proof = append(proof, subTreeMerkleRoot)
			} else {
				subTreeHeight := int64(math.Ceil(math.Log2(float64(rightBound - pivotLeaf))))

				heightDifference := expectedHeight - subTreeHeight

				remainingNodesHash, err := queryRootHash(ctx, client, offset+pivotLeaf+1, offset+rightBound)
				if err != nil {
					return nil, err
				}

				leafRoots := recursiveZeroHash(subTreeHeight)

				leaves := make([]common.Hash, int(math.Pow(2, float64(heightDifference))))
				for i := 0; i < len(leaves); i++ {
					leaves[i] = leafRoots
				}
				leaves[0] = remainingNodesHash

				merkleTree, err := NewMerkleTree(leaves)
				if err != nil {
					return nil, err
				}

				subTreeMerkleRoot := merkleTree.GetRoot()
				proof = append(proof, subTreeMerkleRoot)
			}

			rightBound = newRightBound
		}
	}

	return reverseProof(proof), nil
}

func reverseProof(proofs []common.Hash) []common.Hash {
	reversed := make([]common.Hash, 0)
	for i := len(proofs) - 1; i >= 0; i-- {
		reversed = append(reversed, proofs[i])
	}
	return reversed
}

func queryRootHash(ctx context.Context, client types.IClient, startBlock, endBlock int64) (common.Hash, error) {
	var payload interface{}
	err := client.Rpc().CallContext(ctx, &payload, "eth_getRootHash",
		big.NewInt(startBlock),
		big.NewInt(endBlock),
	)

	if err != nil {
		return common.Hash{}, err
	}

	client.Logger().Debug("queryRootHash", log.Fields{
		"start": startBlock,
		"end":   endBlock,
		"hash":  common.HexToHash(payload.(string)),
	})

	return common.HexToHash(payload.(string)), nil
}

func recursiveZeroHash(n int64) common.Hash {
	if n == 0 {
		return common.Hash{}
	}
	subHash := recursiveZeroHash(n - 1)
	b, _ := argsZeroHash.Pack(subHash, subHash)
	return crypto.Keccak256Hash(b)
}

func GetReceiptProof(ctx context.Context, client types.IClient, txReceipt *ether.Receipt, block *ether.Block) ([]byte, []byte, error) {
	client.Logger().Debug("GetReceiptProof", log.Fields{
		"txReceipt": txReceipt.TxHash.String(),
		"block":     block.NumberU64(),
	})

	stateSyncTxHash := getStateSyncTxHash(block)
	receiptsTrie := trie.NewTrie()
	for _, tx := range block.Transactions() {
		if tx.Hash() == stateSyncTxHash {
			continue
		}

		client.Logger().Debug("TransactionReceipt", log.Fields{
			"txHash": tx.Hash().String(),
		})

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			continue
		}

		raw, err := receipt.MarshalBinary()
		if err != nil {
			return nil, nil, err
		}

		path, err := rlp.EncodeToBytes(receipt.TransactionIndex)
		if err != nil {
			return nil, nil, err
		}
		receiptsTrie.Put(path, raw)
	}

	path, err := rlp.EncodeToBytes(txReceipt.TransactionIndex)
	if err != nil {
		return nil, nil, err
	}

	_, stacks, ok := receiptsTrie.FindPath(path)
	if !ok {
		return nil, nil, fmt.Errorf("not found")
	}

	var rawParentNodes [][][]byte
	for _, stack := range stacks {
		var rawBytes [][]byte
		for _, raw := range stack.Raw() {
			rawBytes = append(rawBytes, raw.([]byte))
		}
		rawParentNodes = append(rawParentNodes, rawBytes)
	}

	parentNodes, err := rlp.EncodeToBytes(rawParentNodes)
	if err != nil {
		return nil, nil, err
	}

	client.Logger().Debug("ReceiptProof", log.Fields{
		"path":  hexutil.Encode(path),
		"proof": hexutil.Encode(parentNodes),
	})
	return path, parentNodes, nil
}

func getStateSyncTxHash(block *ether.Block) common.Hash {
	return GetDerivedBorTxHash(BorReceiptKey(block.NumberU64(), block.Hash()))
}

var (
	borReceiptPrefix = []byte("matic-bor-receipt-") // borReceiptPrefix + number + block hash -> bor block receipt
)

// BorReceiptKey = borReceiptPrefix + num (uint64 big endian) + hash
func BorReceiptKey(number uint64, hash common.Hash) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return append(append(borReceiptPrefix, enc...), hash.Bytes()...)
}

func GetDerivedBorTxHash(receiptKey []byte) common.Hash {
	return common.BytesToHash(crypto.Keccak256(receiptKey))
}
