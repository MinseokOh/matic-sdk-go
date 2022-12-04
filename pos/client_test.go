package pos

import (
	"context"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestClient_BuildPayloadForExit(t *testing.T) {
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	txHash := common.HexToHash("0xb005d8db45f33836c422ee18286fa8ebe49b4ec7b9930e673d85ecd081cc3b8e")
	payload, err := client.BuildPayloadForExit(context.Background(), txHash, types.ERC20Transfer)
	assert.NoError(t, err)

	expectedPayload := `0xf9073e8423d1ec00b901000068067197f28641cdbdf52915e48cc4eb65a6a74b37694bbee0ad96ce467519432bcb168e3932a52c0d460fb2e93a4fc730e5bd0ccc9e42c2755a04a0800bc7fd0bb97fc4896508954b02c19b90fd244a391e50fe0cd6b979ff1afa1ac2a2e69bf8ddfc998afe2723c5f99b9cca0c960953809513203979e077808d59620556564de85c8772fbd0b15dce8168ef9965d53cef0947d80fe8e3d562c5fccdffa8ea7fbf33ba48e44e44e2e2b5ca8a61bb32bcb7f14f7eb7e853f6698681a0356da1c79f16dc17a2aaab85e5f759e18a40db750f262155e489c5f9b8dfa759ef1c607e346e173cb8dcb5786265cc1264562462e3d0d4a14122bc088c6e35c2316784013f8974846180c3c8a01c231c504b86cd2bbf2360e61bf747207ee2fbfc6f35005c328005c600e9255ca0ea31d8804cd84c283be9702b4de7eff615884681cc7c4cfdc65a46673eeb566cb902eaf902e701828547b9010000000000400000020000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000008000000800000000000000000000100000000000000000000020000000000000000001800000000000000000080000010000000000000000000010000000000000000000000040000000000000000000000000000200000000000000020000000000000010001000000000000000000000000004000000003000000000001000000000000000000000000000000100000000020000000000000000000000000000000000000000000000000000000000000100000f901ddf89b94fe4f5145f6e09952a5ba9e956ed0c25e3fa4c7f1f863a0ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa0000000000000000000000000bbca830ee5dcabde33db24496b5524b9c5a69fe6a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000011f9013d940000000000000000000000000000000000001010f884a04dfe1bbbcf077ddc3e01291eea2d5c70c2b422b415d95645b9adcfd678cb1d63a00000000000000000000000000000000000000000000000000000000000001010a0000000000000000000000000bbca830ee5dcabde33db24496b5524b9c5a69fe6a0000000000000000000000000c26880a0af2ea0c7e8130e6ec47af756465452e8b8a000000000000000000000000000000000000000000000000000007c1fcb8018000000000000000000000000000000000000000000000000056bc077dda68980000000000000000000000000000000000000000000000001ded4a9e3b2e378e3ca0000000000000000000000000000000000000000000000056bbffbbddb0968000000000000000000000000000000000000000000000001ded4aa5fd2aef8fbcab902f6f902f3f902f0822080b902eaf902e701828547b9010000000000400000020000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000008000000800000000000000000000100000000000000000000020000000000000000001800000000000000000080000010000000000000000000010000000000000000000000040000000000000000000000000000200000000000000020000000000000010001000000000000000000000000004000000003000000000001000000000000000000000000000000100000000020000000000000000000000000000000000000000000000000000000000000100000f901ddf89b94fe4f5145f6e09952a5ba9e956ed0c25e3fa4c7f1f863a0ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa0000000000000000000000000bbca830ee5dcabde33db24496b5524b9c5a69fe6a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000011f9013d940000000000000000000000000000000000001010f884a04dfe1bbbcf077ddc3e01291eea2d5c70c2b422b415d95645b9adcfd678cb1d63a00000000000000000000000000000000000000000000000000000000000001010a0000000000000000000000000bbca830ee5dcabde33db24496b5524b9c5a69fe6a0000000000000000000000000c26880a0af2ea0c7e8130e6ec47af756465452e8b8a000000000000000000000000000000000000000000000000000007c1fcb8018000000000000000000000000000000000000000000000000056bc077dda68980000000000000000000000000000000000000000000000001ded4a9e3b2e378e3ca0000000000000000000000000000000000000000000000056bbffbbddb0968000000000000000000000000000000000000000000000001ded4aa5fd2aef8fbca82008080`
	assert.Equal(t, expectedPayload, hexutil.Encode(payload))
}

func TestClient_DepositEtherFor(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
	assert.NoError(t, err)

	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	hash, err := client.DepositEtherFor(context.Background(), big.NewInt(10000), privateKey)
	assert.NoError(t, err)
	t.Log("txHash", hash.String())
}

func TestClient_IsCheckPointed(t *testing.T) {
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	checkPointed, err := client.IsCheckPointed(context.Background(), common.HexToHash("0xc55da852f91aad02018e92870cc440928c7ef4693e3fc5dcf8b31df58ae97f94"))
	assert.NoError(t, err)
	assert.Equal(t, checkPointed, true)
}