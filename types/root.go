package types

import "math/big"

type RootBlockInfo struct {
	HeaderBlockNumber *big.Int
	Start             *big.Int
	End               *big.Int
}
