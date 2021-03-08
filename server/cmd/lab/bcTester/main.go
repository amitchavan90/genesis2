package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"genesis/blockchain"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ninja-software/terror"
)

func main() {
	// arguments:
	// makeContract
	// balance
	// send

	blk := &blockchain.Service{}

	b, err := base64.StdEncoding.DecodeString(cfgBlockchain.PrivateKeyBytes)
	if err != nil {
		panic("decode base64: " + err.Error())
	}
	key, err := keystore.DecryptKey(b, cfgBlockchain.PrivateKeyPassword)
	if err != nil {
		panic("decrypt key: " + err.Error())
	}

	privateKey = key.PrivateKey

	ctx := context.Background()

	blk, err = blockchain.New(ctx, cfgBlockchain.EthereumHost, privateKey, false)
	if err != nil {
		panic("failed to init blockchain cfg. " + err.Error())
	}

	args := os.Args[1:]
	if len(args) == 0 {
		panic("no argument given")
	}
	arg := args[0]

	switch arg {
	case "makeContract":
		err = makeContract(ctx, blk)
	case "balance":
		err = balance(ctx, blk)
	case "send":
		err = sendEther(ctx, blk)
	default:
		err = fmt.Errorf("unknown option. " + arg)
	}

	if err != nil {
		terror.Echo(err)
	}
}
