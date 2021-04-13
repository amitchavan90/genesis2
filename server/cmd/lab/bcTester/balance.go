package main

import (
	"context"
	"genesis/blockchain"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninja-software/terror"
)

/*
 * get current ether wallet balance
 */
func balance(ctx context.Context, blk *blockchain.Service) error {
	balance, err := blk.GetBalance(ctx)
	if err != nil {
		return terror.New(err, "")
	}

	spew.Dump("balance: ", balance)
	return nil
}
