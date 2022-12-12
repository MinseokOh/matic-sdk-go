package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestERC721_Deposit(t *testing.T) {
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	hash, err := client.ERC721(RootDummyERC721, types.Root).Deposit(context.Background(), big.NewInt(800), TestTxOption)
	assert.NoError(t, err)

	t.Log(hash)
}

func TestERC721_DepositMany(t *testing.T) {
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	ids := []*big.Int{
		big.NewInt(800),
		big.NewInt(802),
	}
	hash, err := client.ERC721(RootDummyERC721, types.Root).DepositMany(context.Background(), ids, TestTxOption)
	assert.NoError(t, err)

	t.Log(hash)
}
