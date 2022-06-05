package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mangoim/goeth-demo/contracts/counter"
	"log"
	"math/big"
)

func main() {
	// use bsc testnet
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545")
	if err != nil {
		log.Fatal(err)
	}

	// put in your test private key, make sure it has bsc testnet BNB
	privateKey, err := crypto.HexToECDSA("private key")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// contract address
	address := common.HexToAddress("0x392698b12759AB3CB5dBADBFC1601a24d3d4b7aD")
	instance, err := counter.NewCounter(address, client)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(97))
	if err != nil {
		log.Fatal(err)
	}

	transactOpts := &bind.TransactOpts{
		From:     fromAddress,
		Nonce:    auth.Nonce,
		Signer:   auth.Signer,
		Value:    big.NewInt(0),
		GasPrice: auth.GasPrice,
		GasLimit: auth.GasLimit,
		Context:  auth.Context,
		NoSend:   false,
	}

	// send tx
	tx, err := instance.Inc(transactOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tx sent: %s", tx.Hash().Hex())

	// get count
	count, err := instance.Get(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("count: %d", count)
}
