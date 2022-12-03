package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math"
)

type MerkleTree struct {
	leaves []common.Hash
	layers [][]common.Hash
}

func NewMerkleTree(leaves []common.Hash) (*MerkleTree, error) {
	merkleTree := MerkleTree{
		leaves: make([]common.Hash, len(leaves)),
		layers: make([][]common.Hash, 0),
	}

	if len(leaves) < 1 {
		return nil, fmt.Errorf("at least 1 leaf needed")
	}

	depth := int(math.Ceil(math.Log(float64(len(leaves))) / math.Log(2)))
	if depth > 20 {
		return nil, fmt.Errorf("depth must be 20 or less")
	}

	zeroLeaves := make([]common.Hash, int(math.Pow(2, float64(depth)))-len(leaves))
	merkleTree.leaves = concat(leaves, zeroLeaves)
	merkleTree.layers = append(merkleTree.layers, merkleTree.leaves)
	merkleTree.createHashes(merkleTree.leaves)

	return &merkleTree, nil
}

func (merkleTree *MerkleTree) createHashes(nodes []common.Hash) {
	if len(nodes) == 1 {
		return
	}

	treeLevel := make([]common.Hash, 0)
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		right := nodes[i+1]

		treeLevel = append(treeLevel,
			crypto.Keccak256Hash(append(left.Bytes(), right.Bytes()...)),
		)
	}

	if len(nodes)%2 == 1 {
		treeLevel = append(treeLevel, nodes[len(nodes)-1])
	}

	merkleTree.layers = append(merkleTree.layers, treeLevel)
	merkleTree.createHashes(treeLevel)
}

func (merkleTree *MerkleTree) GetLeaves() []common.Hash   { return merkleTree.leaves }
func (merkleTree *MerkleTree) GetLayers() [][]common.Hash { return merkleTree.layers }
func (merkleTree *MerkleTree) GetRoot() common.Hash       { return merkleTree.layers[len(merkleTree.layers)-1][0] }

func concat(arr1, arr2 []common.Hash) []common.Hash {
	return append(arr1, arr2...)
}
