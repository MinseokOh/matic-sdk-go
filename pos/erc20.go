package pos

import (
	"context"
	"crypto/ecdsa"
	"github.com/MinseokOh/matic-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	ether "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type ERC20 struct {
	*Client

	address common.Address
}

func newERC20(address common.Address, client *Client) *ERC20 {
	return &ERC20{
		Client:  client,
		address: address,
	}
}

func (erc20 *ERC20) Approve(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) {

}

func (erc20 *ERC20) Deposit(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) {

}

func (erc20 *ERC20) Withdraw(ctx context.Context, amount *big.Int, privateKey *ecdsa.PrivateKey) error {
	childClient := erc20.Client.Child

	address := crypto.PubkeyToAddress(privateKey.Public().(ecdsa.PublicKey))
	chainId, err := childClient.ChainID(ctx)
	if err != nil {
		return err
	}

	nonce, err := childClient.PendingNonceAt(ctx, address)
	if err != nil {
		return err
	}

	gasTipCap, err := childClient.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		Gas:       5e6,
		Nonce:     nonce,
		To:        &erc20.address,
		Value:     big.NewInt(0),
	})
	if err != nil {
		return err
	}

	err = childClient.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	return nil

}

func (erc20 *ERC20) Exit(ctx context.Context, txHash common.Hash, privateKey *ecdsa.PrivateKey) error {
	rootClient := erc20.Client.Root

	address := crypto.PubkeyToAddress(privateKey.Public().(ecdsa.PublicKey))
	chainId, err := rootClient.ChainID(ctx)
	if err != nil {
		return err
	}

	nonce, err := rootClient.PendingNonceAt(ctx, address)
	if err != nil {
		return err
	}

	gasTipCap, err := rootClient.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	data, err := erc20.Client.BuildPayloadForExit(ctx, txHash, types.ERC20Transfer)
	if err != nil {
		return err
	}

	signer := ether.NewLondonSigner(chainId)
	tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: gasTipCap,
		Gas:       1e7,
		Nonce:     nonce,
		To:        &rootClient.config.RootChainManager,
		Value:     big.NewInt(0),
		Data:      data,
	})
	if err != nil {
		return err
	}

	err = rootClient.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}
