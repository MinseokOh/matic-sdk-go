## Matic SDK for Go

> **Warning**
>
> Initial development is in progress, but there has not yet been a stable.


This repository contains the matic go client library. converted [maticjs](https://github.com/maticnetwork/matic.js) in go and compatible [go-ethereum](https://github.com/ethereum/go-ethereum)

This library will help developers to move assets from Ethereum chain to Matic chain, and withdraw from Matic to Ethereum using fraud proofs.


> **Note**
>
> Requires Go [1.19+](https://go.dev/dl/)


---


### Setup Client

```go
import (
    "github.com/MinseokOh/matic-sdk-go/pos"
    "github.com/MinseokOh/matic-sdk-go/types"
)

posClient, err := pos.NewClient(types.NewDefaultConfig(types.TestNet))
if err != nil {
	panic(err)
}
```


---


### Ether Deposit and Withdraw Guide

#### Deposit ETH

- Make the depositEtherFor call on the RootChainManager and send the ether asset.

```go
hash, err := posClient.DepositEtherFor(context.Background(), big.NewInt(10000), &types.TxOption{
	PrivateKey: privateKey
})
if err != nil {
    // handle error
}
fmt.Println(hash)
```

#### Withdraw ETH

1. ***Burn*** tokens on Polygon chain.

```go
hash, err := posClient.ERC20(WETHAddress, types.Child).Withdraw(context.Background(), big.NewInt(10000), &types.TxOption{
PrivateKey: privateKey
})
if err != nil {
	// handle error
}
fmt.Println(hash)
```

2. Call exit function on **RootChainManager** to submit proof of burn transaction. This call can be made ***after checkpoint*** is submitted for the block containing burn transaction.

> **Note**
>
> The Withdraw transaction must be checkpointed in order to exit the withdraw.

```go
// token address can be null for native tokens like ethereum or matic
hash, err := posClient.ERC20(common.Address{}, types.Root).Exit(context.Background(), burnTxHash, &types.TxOption{
PrivateKey: privateKey
})
if err != nil {
	// handle error
}
fmt.Println(hash)
```


---


### ERC20 Deposit and Withdraw Guide

#### Deposit ERC20

1. ***Approve ERC20Predicate*** contract to spend the tokens that have to be deposited.

```go
hash, err := posClient.ERC20(rootTokenAddress, types.Root).Approve(context.Background(), big.NewInt(10000), &types.TxOption{
PrivateKey: privateKey
})
if err != nil {
    // handle error
}
fmt.Println(hash)
```

2. Make ***depositFor*** call on ***RootChainManager***.

```go
hash, err := posClient.ERC20(rootTokenAddress, types.Root).DepositFor(context.Background(), big.NewInt(10000), &types.TxOption{
PrivateKey: privateKey
})
if err != nil {
    // handle error
}
fmt.Println(hash)
```

#### Withdraw ERC20

1. ***Burn tokens*** on the Polygon chain.

```go
hash, err := posClient.ERC20(childTokenAddress, types.Child).Withdraw(context.Background(), big.NewInt(10000), &types.TxOption{
PrivateKey: privateKey
})
if err != nil {
    // handle error
}
fmt.Println(hash)
```

2. Call the `exit()` function on ***RootChainManager*** to submit proof of burn transaction. This call can be made after the checkpoint is submitted for the block containing the burn transaction.


> **Note**
>
> The Withdraw transaction must be checkpointed in order to exit the withdraw.


```go
hash, err := posClient.ERC20(rootTokenAddress, types.Root).Exit(context.Background(), burnTxHash, &types.TxOption{
PrivateKey: privateKey
})
if err != nil {
	// handle error
}
fmt.Println(hash)
```


---


### Exit With Raw CallData

```go
chainId, err := posClient.Root.ChainID(ctx)
if err != nil {
    // handle error
}

nonce, err := posClient.Root.PendingNonceAt(ctx, address)
if err != nil {
    // handle error
}

gasTipCap, err := posClient.Root.SuggestGasTipCap(ctx)
if err != nil {
    // handle error
}

payload, err := posClient.BuildPayloadForExit(ctx, txHash, types.ERC20Transfer)
if err != nil {
    // handle error
}

data, err := maticabi.RootChainManager.Pack("exit", payload)
if err != nil {
    // handle error
}

signer := ether.NewLondonSigner(chainId)
tx, err := ether.SignNewTx(privateKey, signer, &ether.DynamicFeeTx{
    ChainID:   chainId,
    GasTipCap: gasTipCap,
    GasFeeCap: gasTipCap,
    Gas:       1e6,
    Nonce:     nonce,
    To:        rootChainManagerAddress,
    Value:     big.NewInt(0),
    Data:      data,
})
if err != nil {
    // handle error
}

err = posClient.Root.SendTransaction(ctx, tx)
if err != nil {
    // handle error
}
```
