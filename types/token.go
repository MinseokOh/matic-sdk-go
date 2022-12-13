package types

type TokenType int

const (
	ERC20   = TokenType(1)
	ERC721  = TokenType(2)
	ERC1155 = TokenType(3)
)

func (token TokenType) String() string {
	switch token {
	case ERC20:
		return "erc20"
	case ERC721:
		return "erc721"
	case ERC1155:
		return "erc1155"
	}

	return ""
}
