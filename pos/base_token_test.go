package pos

import (
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseToken_getPredicateAddress(t *testing.T) {
	client, err := NewClient(types.NewDefaultConfig(types.TestNet))
	assert.NoError(t, err)

	erc20Predicate := client.ERC20(RootDummyERC20, types.Root).PredicateAddress()
	assert.Equal(t, erc20Predicate, common.HexToAddress("0xdD6596F2029e6233DEFfaCa316e6A95217d4Dc34"))

	erc721Predicate := client.ERC721(RootDummyERC721, types.Root).PredicateAddress()
	assert.Equal(t, erc721Predicate, common.HexToAddress("0x56E14C4C1748a818a5564D33cF774c59EB3eDF59"))
}
