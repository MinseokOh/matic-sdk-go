package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	RootDummyERC20  = common.HexToAddress("0x655f2166b0709cd575202630952d71e2bb0d61af")
	ChildDummyERC20 = common.HexToAddress("0xfe4f5145f6e09952a5ba9e956ed0c25e3fa4c7f1")
	ChildWETH       = common.HexToAddress("0xA6FA4fB5f76172d178d61B04b0ecd319C5d1C0aa")
	Matic           = common.HexToAddress(types.MaticAddress)
)

func TestERC20_Approve(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Approve(context.Background(), big.NewInt(123456789), privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_DepositFor(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.DepositFor(context.Background(), big.NewInt(123456789), privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_Withdraw(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(ChildDummyERC20, types.Child)
	hash, err := erc20.Withdraw(context.Background(), big.NewInt(123456789), privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_WithdrawMatic(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(Matic, types.Child)
	hash, err := erc20.Withdraw(context.Background(), big.NewInt(123456789), privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_WithdrawEther(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(ChildWETH, types.Child)
	hash, err := erc20.Withdraw(context.Background(), big.NewInt(123456789), privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_Exit(t *testing.T) {
	txHash := common.HexToHash("0xc55da852f91aad02018e92870cc440928c7ef4693e3fc5dcf8b31df58ae97f94")
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Exit(context.Background(), txHash, privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_ExitEther(t *testing.T) {
	txHash := common.HexToHash("0xc55da852f91aad02018e92870cc440928c7ef4693e3fc5dcf8b31df58ae97f94")
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20 := client.ERC20(RootDummyERC20, types.Root)
	hash, err := erc20.Exit(context.Background(), txHash, privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestERC20_Balance(t *testing.T) {
	address := common.HexToAddress("0x97524878fa80e4c335b64b8f07c888811d8a201e")
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	rootBalance, err := client.ERC20(RootDummyERC20, types.Root).BalanceOf(context.Background(), address)
	assert.NoError(t, err)
	childBalance, err := client.ERC20(ChildWETH, types.Child).BalanceOf(context.Background(), address)
	assert.NoError(t, err)

	t.Log(rootBalance, childBalance)
}
