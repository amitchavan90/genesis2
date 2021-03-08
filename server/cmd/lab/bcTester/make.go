package main

import (
	"context"
	"genesis/blockchain"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninja-software/terror"
)

func makeContract(ctx context.Context, blk *blockchain.Service) error {
	address, tx, contract, err := blk.DeploySmartContract(ctx, "Proof of Steak v3")
	if err != nil {
		return terror.New(err, "")
	}

	spew.Dump("address: ", address)
	spew.Dump("tx:      ", tx)
	spew.Dump("contract:", contract)

	return nil
}
