package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestERC20_Approve(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Approve(context.Background(), common.Address{}, big.NewInt(123456789), TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_ApproveMax(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.ApproveMax(context.Background(), common.Address{}, TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_DepositFor(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Deposit(context.Background(), big.NewInt(123456789), TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_Withdraw(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(ChildDummyERC20, types.Child)
	hash, err := erc20.Withdraw(context.Background(), big.NewInt(123456789), TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_WithdrawMatic(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(Matic, types.Child)
	hash, err := erc20.Withdraw(context.Background(), big.NewInt(123456789), TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_WithdrawEther(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(ChildWETH, types.Child)
	hash, err := erc20.Withdraw(context.Background(), big.NewInt(123456789), TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_Exit(t *testing.T) {
	txHash := common.HexToHash("0xe2f5f63d36fea883fc2514e70f0f49a3c006e27e81a08acf8857da9104b15f50")
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Exit(context.Background(), txHash, TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_ExitEther(t *testing.T) {
	txHash := common.HexToHash("0xc55da852f91aad02018e92870cc440928c7ef4693e3fc5dcf8b31df58ae97f94")
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Exit(context.Background(), txHash, TestTxOption)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_Balance(t *testing.T) {
	client, err := NewClient(NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	rootBalance, err := client.ERC20(RootDummyERC20, types.Root).BalanceOf(context.Background(), TestTxOption.From())
	assert.NoError(t, err)

	childBalance, err := client.ERC20(ChildWETH, types.Child).BalanceOf(context.Background(), TestTxOption.From())
	assert.NoError(t, err)

	t.Log(rootBalance, childBalance)
}
