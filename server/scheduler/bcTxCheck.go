package scheduler

import (
	"context"
	"encoding/hex"
	"fmt"
	"genesis"
	"genesis/blockchain"
	"genesis/db"
	"genesis/store"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/ninja-software/terror"
	"github.com/prometheus/common/log"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// BlockchainTransactionConfirmCheck checks if the blockchain transaction is confirmed on the network
type BlockchainTransactionConfirmCheck struct {
	Conn             *sqlx.DB
	TransactionStore *store.Transaction
	ManifestStore    *store.Manifest
	Blk              *blockchain.Service
}

// Runner will check unconfirmed blockchain transaction, return hows many successfully confirmed
func (bcTxCk *BlockchainTransactionConfirmCheck) Runner() (int, error) {
	var err error

	// sanity check
	if bcTxCk.TransactionStore == nil {
		return 0, terror.New(fmt.Errorf("txStore is nil"), "")
	}
	if bcTxCk.ManifestStore == nil {
		return 0, terror.New(fmt.Errorf("maniStore is nil"), "")
	}
	if bcTxCk.Blk == nil {
		return 0, terror.New(fmt.Errorf("blk is nil"), "")
	}

	// scrub transactions table clean and make data valid
	err = scrubTransactions(bcTxCk)
	if err != nil {
		return 0, terror.New(err, "")
	}

	// recover any pending manifest if there is any
	err = recoverPendingManifest(bcTxCk)
	if err != nil {
		return 0, terror.New(err, "")
	}

	// shorthand
	manifestStore := bcTxCk.ManifestStore
	blk := bcTxCk.Blk
	ctx := context.Background()

	// Get unconfirmed manifest, which shouldnt be. or wait until previous manifest is in the blockchain
	unconfirmedManifests, err := manifestStore.AllUnconfirmed()
	if err != nil {
		return 0, terror.New(err, "get unconfirmed manifests")
	}
	if len(unconfirmedManifests) == 0 {
		return 0, nil
	}

	count := 0
	for _, m := range unconfirmedManifests {
		txh := m.TransactionHash
		if !txh.Valid {
			continue
		}
		_, isPending, err := blk.Client.TransactionByHash(ctx, common.HexToHash(txh.String))
		if err != nil {
			log.Errorln(err)
			genesis.SentrySend(ctx, nil, nil, terror.New(err, ""), "get blockchain transaction detail failed")
			continue
		}

		if !isPending {
			m.Confirmed = true
			_, err = manifestStore.Update(m)
			if err != nil {
				log.Errorln(err)
				genesis.SentrySend(ctx, nil, nil, terror.New(err, ""), "get blockchain transaction detail failed")
			}
			count++
		}
	}

	return count, nil
}

// scrub transaction table to make table data valid again,
// for some situation where the transaction has manifest id but no transaction hash, fill it
func scrubTransactions(bcTxCk *BlockchainTransactionConfirmCheck) error {
	conn := bcTxCk.Conn

	// get dirty manifest and tries to clean it
	dirtyManifests, err := db.Transactions(
		db.TransactionWhere.ManifestID.IsNotNull(),
		db.TransactionWhere.TransactionHash.IsNull(),
	).All(conn)
	if err != nil {
		return terror.New(err, "get pending manifests")
	}

	for i, t := range dirtyManifests {
		log.Infoln("cleaning transaction", i, t.ID)

		m, err := db.FindManifest(conn, t.ManifestID.String)
		if err != nil {
			log.Errorln(err)
		}

		t.TransactionHash = m.TransactionHash
		_, err = t.Update(conn, boil.Whitelist(db.TransactionColumns.TransactionHash))
		if err != nil {
			log.Errorln(err)
		}
	}

	return nil
}

func recoverPendingManifest(bcTxCk *BlockchainTransactionConfirmCheck) error {
	// shorthand
	manifestStore := bcTxCk.ManifestStore
	blk := bcTxCk.Blk
	ctx := context.Background()

	// Get pending manifest and tries to fix it
	pendingManifests, err := manifestStore.AllPending()
	if err != nil {
		return terror.New(err, "get pending manifests")
	}

	// nothing, return
	if len(pendingManifests) == 0 {
		return nil
	}

	for _, m := range pendingManifests {
		// sanity check
		if m.TransactionNonce <= 0 {
			continue
		}
		if !m.Pending || m.Confirmed {
			continue
		}
		if !m.CompiledText.Valid || len(m.CompiledText.Bytes) == 0 {
			continue
		}
		if len(m.TransactionHash.String) > 0 {
			continue
		}
		if len(m.MerkleRootSha256.String) <= 0 {
			continue
		}

		log.Infoln("manifest recovering", m.MerkleRootSha256.String)

		// Get blockchain contract
		settings, err := bcTxCk.TransactionStore.GetSettings()
		if err != nil {
			return terror.New(err, "get settings")
		}

		contract, err := blk.GetContract(common.HexToAddress(settings.SmartContractAddress))
		if err != nil {
			return terror.New(err, "get blockchain contract")
		}

		// Update Gas
		err = blk.UpdateGas(ctx)
		if err != nil {
			return terror.New(err, "update gas")
		}

		// Commit manifest merkleRoot to smart contract
		// decode hex to bytes
		bMR, err := hex.DecodeString(m.MerkleRootSha256.String)
		if err != nil {
			return terror.New(err, "decode merkle root hex")
		}
		var biMR = new(big.Int).SetBytes(bMR)
		fakeBlkAuth := *blk.Auth
		fakeBlkAuth.Nonce = new(big.Int).SetInt64(int64(m.TransactionNonce))
		bcTx, err := contract.CommitBatch(&fakeBlkAuth, biMR)
		if err != nil {
			return terror.New(err, "commit manifest to smart contract")
		}

		// save blockchain transaction hash
		bcTxHash := null.StringFrom(bcTx.Hash().Hex())
		m.TransactionHash = bcTxHash

		// unset flag
		m.Pending = false

		// Update manifest
		_, err = manifestStore.Update(m)
		if err != nil {
			return terror.New(err, "update manifest")
		}

		log.Infoln("manifest recovered", m.MerkleRootSha256.String, m.TransactionHash.String)
	}

	return terror.New(fmt.Errorf("pending manifest remaining, please wait until recovered"), "")
}
