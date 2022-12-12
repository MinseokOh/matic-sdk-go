package pos

import (
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	RootDummyERC20  = common.HexToAddress("0x655f2166b0709cd575202630952d71e2bb0d61af")
	ChildDummyERC20 = common.HexToAddress("0xfe4f5145f6e09952a5ba9e956ed0c25e3fa4c7f1")

	RootDummyERC721  = common.HexToAddress("0x16F7EF3774c59264C46E5063b1111bCFd6e7A72f")
	ChildDummyERC721 = common.HexToAddress("0xbD88C3A7c0e242156a46Fbdf87141Aa6D0c0c649")

	RootDummyERC1155  = common.HexToAddress("0x2e3Ef7931F2d0e4a7da3dea950FF3F19269d9063")
	ChildDummyERC1155 = common.HexToAddress("0xA07e45A987F19E25176c877d98388878622623FA")

	ChildWETH = common.HexToAddress("0xA6FA4fB5f76172d178d61B04b0ecd319C5d1C0aa")
	Matic     = common.HexToAddress(types.MaticAddress)

	TestPrivateKey, _ = crypto.HexToECDSA("1c28edecd1cdfbdb2e32c38d8e06ed042f3e31fb05d9884e5322376cce4706d4")
)

func TestBaseToken_getPredicateAddress(t *testing.T) {
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20Predicate := client.ERC20(RootDummyERC20, types.Root).PredicateAddress()
	assert.Equal(t, erc20Predicate, common.HexToAddress("0xdD6596F2029e6233DEFfaCa316e6A95217d4Dc34"))

	erc721Predicate := client.ERC721(RootDummyERC721, types.Root).PredicateAddress()
	assert.Equal(t, erc721Predicate, common.HexToAddress("0x56E14C4C1748a818a5564D33cF774c59EB3eDF59"))
}
