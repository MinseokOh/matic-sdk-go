package types

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ether "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"math/big"
)

const (
	LegacyTxType = iota
	AccessListTxType
	DynamicFeeTxType
)

var (
	EmptyPrivateKey = fmt.Errorf("empty private key")
	EmptyTxOption   = fmt.Errorf("tx option is nil")
)

type TxOption struct {
	// PrivateKey required
	PrivateKey *ecdsa.PrivateKey

	// Optional Parameters
	TxType int

	// GasLimit : gas limit for tx
	GasLimit uint64

	// GasPrice : gas price for LegacyTxType
	GasPrice *big.Int

	// GasTipCap : maxPriorityFeePerGas for DynamicFeeTxType
	GasTipCap *big.Int

	// GasFeeCap : maxFeePerGas for DynamicFeeTxType
	GasFeeCap *big.Int

	ChainId *big.Int
	Nonce   uint64

	data  []byte
	value *big.Int
	to    common.Address
}

func ValidateTxOption(txOption *TxOption) error {
	if txOption == nil {
		return EmptyTxOption
	}

	if err := txOption.Validate(); err != nil {
		return err
	}

	return nil
}

func (txOption *TxOption) Validate() error {
	if txOption.PrivateKey == nil {
		return EmptyPrivateKey
	}

	return nil
}

func (txOption *TxOption) SetTxData(to common.Address, data []byte, value *big.Int) *TxOption {
	txOption.to = to
	txOption.data = data
	txOption.value = value

	return txOption
}

func (txOption *TxOption) From() common.Address {
	return crypto.PubkeyToAddress(*txOption.PrivateKey.Public().(*ecdsa.PublicKey))
}

func (txOption *TxOption) Build(ctx context.Context, client IClient) (*ether.Transaction, error) {
	var err error
	if txOption.PrivateKey == nil {
		return nil, fmt.Errorf("empty private key")
	}

	if txOption.ChainId == nil {
		txOption.ChainId, err = client.ChainID(ctx)
		if err != nil {
			return nil, err
		}
	}

	if txOption.GasLimit == 0 {
		txOption.GasLimit = 1e6
	}

	if txOption.Nonce == 0 {
		txOption.Nonce, err = client.NonceAt(ctx, txOption.From(), nil)
		if err != nil {
			return nil, err
		}
	}

	switch txOption.TxType {
	case LegacyTxType:
		if txOption.GasPrice == nil {
			txOption.GasPrice, err = client.SuggestGasPrice(ctx)
			if err != nil {
				return nil, err
			}
		}

		client.Logger().Debug("Sign Transaction", log.Fields{
			"@type":    "LegacyTxType",
			"nonce":    txOption.Nonce,
			"gasPrice": txOption.GasPrice,
			"gas":      txOption.GasLimit,
			"to":       txOption.to,
			"value":    txOption.value,
			"data":     hexutil.Encode(txOption.data),
		})
		return ether.SignNewTx(txOption.PrivateKey, ether.LatestSignerForChainID(txOption.ChainId), &ether.LegacyTx{
			Nonce:    txOption.Nonce,
			GasPrice: txOption.GasPrice,
			Gas:      txOption.GasLimit,
			To:       &txOption.to,
			Value:    txOption.value,
			Data:     txOption.data,
		})
	case DynamicFeeTxType:
		if txOption.GasTipCap == nil {
			txOption.GasTipCap, err = client.SuggestGasTipCap(ctx)
			if err != nil {
				return nil, err
			}
		}

		if txOption.GasFeeCap == nil {
			txOption.GasFeeCap = txOption.GasTipCap
		}

		client.Logger().Debug("Sign Transaction", log.Fields{
			"@type":     "DynamicFeeTxType",
			"nonce":     txOption.Nonce,
			"gasTipCap": txOption.GasTipCap,
			"gasFeeCap": txOption.GasFeeCap,
			"gas":       txOption.GasLimit,
			"to":        txOption.to,
			"value":     txOption.value,
			"data":      hexutil.Encode(txOption.data),
		})
		return ether.SignNewTx(txOption.PrivateKey, ether.LatestSignerForChainID(txOption.ChainId), &ether.DynamicFeeTx{
			ChainID:   txOption.ChainId,
			Nonce:     txOption.Nonce,
			GasTipCap: txOption.GasTipCap,
			GasFeeCap: txOption.GasFeeCap,
			Gas:       txOption.GasLimit,
			To:        &txOption.to,
			Value:     txOption.value,
			Data:      txOption.data,
		})
	}
	return nil, fmt.Errorf("invalid tx type: %d", txOption.TxType)
}
