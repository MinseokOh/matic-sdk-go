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

	hash, err := client.ERC721(RootDummyERC721, types.Root).Deposit(context.Background(), big.NewInt(800), &types.TxOption{
		PrivateKey: TestPrivateKey,
		TxType:     types.DynamicFeeTxType,
	})
	assert.NoError(t, err)
	t.Log(hash)
}
