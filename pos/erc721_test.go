package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestERC721_Deposit(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	hash, err := client.ERC721(RootDummyERC721, types.Root).Deposit(context.Background(), big.NewInt(800), TestTxOption)
	assert.NoError(t, err)

	t.Log("txHash", hash)
}

func TestERC721_DepositMany(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	ids := []*big.Int{
		big.NewInt(800),
		big.NewInt(802),
	}
	hash, err := client.ERC721(RootDummyERC721, types.Root).DepositMany(context.Background(), ids, TestTxOption)
	assert.NoError(t, err)

	t.Log("txHash", hash)
}

func TestERC721_IsApproved(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	approved, err := client.ERC721(RootDummyERC721, types.Root).IsApproved(context.Background(), big.NewInt(805))
	assert.NoError(t, err)

	t.Log("approved", approved)
}

func TestERC721_IsApprovedAll(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	approved, err := client.ERC721(RootDummyERC721, types.Root).IsApprovedAll(context.Background(), TestTxOption.From())
	assert.NoError(t, err)

	t.Log("approved", approved)
}

func TestERC721_Withdraw(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc721 := client.ERC721(ChildDummyERC721, types.Child)
	hash, err := erc721.Withdraw(context.Background(), big.NewInt(801), TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC721_Exit(t *testing.T) {
	txHash := common.HexToHash("0x54f47c891b460369661e22e27eeb4afbbb5dd792c7c8b48cab758892c14ffe85")
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc721 := client.ERC721(RootDummyERC721, types.Root)
	hash, err := erc721.Exit(context.Background(), txHash, TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())

}
