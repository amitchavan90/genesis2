package main

import (
	"context"
	"genesis/blockchain"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ninja-software/terror"
)

/* references
https://goethereumbook.org/transfer-eth/
https://goethereumbook.org/en/client-simulated/
*/

// send ether to an address
func sendEther(ctx context.Context, blk *blockchain.Service) error {
	chainID, err := blk.Client.ChainID(ctx)
	if err != nil {
		return terror.New(err, "")
	}

	err = blk.UpdateGas(ctx)
	if err != nil {
		return terror.New(err, "")
	}

	// get gas price
	// typically 1 gwei, 1000000000
	gasPrice, err := blk.Chain.SuggestGasPrice(ctx)
	if err != nil {
		return terror.New(err, "")
	}
	spew.Dump("gasPrice:   ", gasPrice)
	// panic("gasPrice")
	auth := *blk.Auth
	auth.Context = nil
	spew.Dump("auth:    ", auth)

	fromAddress := common.HexToAddress(strAddressFrom)
	nonce, err := blk.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return terror.New(err, "")
	}

	value := big.NewInt(1000000000) // in 1 gwei (0.000000001 eth)
	gasLimit := uint64(42000)       // 2x of 21000 min value

	var data []byte //nil
	toAddress := common.HexToAddress(strAddressTo)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, blk.Auth.GasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return terror.New(err, "")
	}

	// Final Step
	err = blk.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return terror.New(err, "")
	}

	spew.Dump("chainID: ", chainID)
	spew.Dump("tx:      ", tx)
	spew.Dump("signedTx:", signedTx)

	return nil
}
