package genesis

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"genesis/blockchain"
	"genesis/db"
	"genesis/graphql"
	"genesis/helpers"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

const trackContractCreated = "TRACK000"
const trackRegistered = "TRACK001"
const trackMovedToCarton = "TRACK002"
const trackMovedToPallet = "TRACK003"
const trackMovedToContainer = "TRACK004"
const trackRemovedFromCarton = "TRACK005"
const trackRemovedFromPallet = "TRACK006"
const trackRemovedFromContainer = "TRACK007"

////////////////
//  Resolver  //
////////////////

// Transaction resolver
func (r *Resolver) Transaction() graphql.TransactionResolver {
	return &transactionResolver{r}
}

type transactionResolver struct{ *Resolver }

func (r *transactionResolver) TransactionPending(ctx context.Context, obj *db.Transaction) (bool, error) {
	if !obj.TransactionHash.Valid || obj.TransactionHash.String == "-" {
		return false, nil
	}

	_, isPending, err := r.Blk.TransactionReader.TransactionByHash(ctx, common.HexToHash(obj.TransactionHash.String))
	if err != nil {
		fmt.Println(fmt.Errorf("get transaction pending: %w (%s)", err, obj.ID))
		return false, nil
	}
	return isPending, nil
}

func (r *transactionResolver) Action(ctx context.Context, obj *db.Transaction) (*db.TrackAction, error) {
	uuid, err := uuid.FromString(obj.TrackActionID)
	if err != nil {
		return nil, terror.New(err, "get transaction action")
	}
	result, err := TrackActionLoaderFromContext(ctx, uuid)
	if err != nil {
		return nil, terror.New(err, "get transaction action")
	}
	return result, nil
}

func (r *transactionResolver) CreatedBy(ctx context.Context, obj *db.Transaction) (*db.User, error) {
	if !obj.CreatedByID.Valid {
		return nil, nil
	}

	createdByUUID, err := uuid.FromString(obj.CreatedByID.String)
	if err != nil {
		return nil, terror.New(err, "get transaction createdBy")
	}
	result, err := UserLoaderFromContext(ctx, createdByUUID)
	if err != nil {
		return nil, terror.New(err, "get transaction createdBy")
	}

	// Omit user email if not logged in
	userID, _ := r.Auther.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return &db.User{
			ID:           result.ID,
			AffiliateOrg: result.AffiliateOrg,
		}, nil
	}

	return result, nil
}

func (r *transactionResolver) Carton(ctx context.Context, obj *db.Transaction) (*db.Carton, error) {
	if !obj.CartonID.Valid {
		return nil, nil
	}

	cartonUUID, err := uuid.FromString(obj.CartonID.String)
	if err != nil {
		return nil, terror.New(err, "get carton from item transaction")
	}
	carton, err := CartonLoaderFromContext(ctx, cartonUUID)
	if err != nil {
		return nil, terror.New(err, "get carton from item transaction")
	}
	return carton, nil
}
func (r *transactionResolver) Product(ctx context.Context, obj *db.Transaction) (*db.Product, error) {
	if !obj.ProductID.Valid {
		return nil, nil
	}

	productUUID, err := uuid.FromString(obj.ProductID.String)
	if err != nil {
		return nil, terror.New(err, "get product from item transaction")
	}
	product, err := ProductLoaderFromContext(ctx, productUUID)
	if err != nil {
		return nil, terror.New(err, "get product from item transaction")
	}
	return product, nil
}

func (r *transactionResolver) Photos(ctx context.Context, obj *db.Transaction) (*graphql.TransactionPhotos, error) {
	result := &graphql.TransactionPhotos{}

	// only return photos if admin
	userID, _ := r.Auther.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, nil
	}

	// only return photos if action took photo
	actionUUID, err := uuid.FromString(obj.TrackActionID)
	if err != nil {
		return nil, terror.New(err, "get transaction action")
	}
	action, err := TrackActionLoaderFromContext(ctx, actionUUID)
	if err != nil {
		return nil, terror.New(err, "get transaction action")
	}

	requiresPhoto := false
	for _, b := range action.RequirePhotos {
		if b {
			requiresPhoto = true
			break
		}
	}

	if !requiresPhoto {
		return nil, nil
	}

	// get carton photo blob
	if obj.CartonPhotoBlobID.Valid {
		cartonBlobUUID, err := uuid.FromString(obj.CartonPhotoBlobID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}
		result.CartonPhoto, err = r.BlobStore.Get(cartonBlobUUID)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "failed to get carton 'carton photo blob'")
		}
	}

	// get product photo blob
	if obj.ProductPhotoBlobID.Valid {
		productBlobUUID, err := uuid.FromString(obj.ProductPhotoBlobID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}
		result.ProductPhoto, err = r.BlobStore.Get(productBlobUUID)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "failed to get carton 'product photo blob'")
		}
	}

	return result, nil
}

func (r *transactionResolver) Manifest(ctx context.Context, obj *db.Transaction) (*db.Manifest, error) {
	if !obj.ManifestID.Valid {
		return nil, nil
	}
	manifestBlobUUID, err := uuid.FromString(obj.ManifestID.String)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "get transaction manifest")
	}
	manifest, err := r.ManifestStore.Get(manifestBlobUUID)
	if err != nil {
		return nil, terror.New(err, "get transaction manifest")
	}
	return manifest, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) PendingTransactionsCount(ctx context.Context) (int, error) {
	count, err := r.TransactionStore.AllPendingCount()
	if err != nil {
		return 0, terror.New(err, "count pending transactions")
	}
	return int(count), nil
}

func (r *queryResolver) EthereumAccountAddress(ctx context.Context) (string, error) {
	return r.Blk.Auth.From.Hex(), nil
}

func (r *queryResolver) EthereumAccountBalance(ctx context.Context) (string, error) {
	balance, err := r.Blk.GetBalance(ctx)
	if err != nil {
		return "", terror.New(err, "get blockchain balance")
	}

	return balance.String(), nil
}

func (r *queryResolver) Transactions(
	ctx context.Context,
	search graphql.SearchFilter,
	limit int,
	offset int,
	productID *string,
	cartonID *string,
	trackActionID *string,
) (*graphql.TransactionsResult, error) {
	total, results, err := r.TransactionStore.SearchSelect(
		search,
		limit,
		offset,
		null.StringFromPtr(productID),
		null.StringFromPtr(cartonID),
		null.StringFromPtr(trackActionID),
	)
	if err != nil {
		return nil, terror.New(err, "list transactions")
	}

	result := &graphql.TransactionsResult{
		Transactions: results,
		Total:        int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// RecordTransaction adds a contract transaction to the following contract id(s)
func (r *mutationResolver) RecordTransaction(ctx context.Context, input graphql.RecordTransactionInput) (bool, error) {
	// Get Track Action
	action, err := r.TrackActionStore.GetByCode(input.TrackActionCode)
	if err != nil {
		return false, terror.New(err, "invalid track action code")
	}

	// Get user
	user, err := r.Auther.UserFromContext(ctx)
	if err != nil {
		return false, terror.New(terror.ErrBadContext, "")
	}
	createdByName := user.AffiliateOrg.String
	if createdByName == "" {
		createdByName = user.LastName.String
	}

	// Get product IDs
	var cartons db.CartonSlice
	var products db.ProductSlice
	cartonScanTimes := []*time.Time{}
	productScanTimes := []*time.Time{}

	if len(input.ProductIDs) > 0 {
		_products, err := r.ProductStore.GetMany(input.ProductIDs)
		if err != nil {
			return false, terror.New(err[0], "on track action")
		}

		products = _products
		productScanTimes = input.ProductScanTimes
	}
	if len(input.CartonIDs) > 0 {
		_cartons, err := r.CartonStore.GetMany(input.CartonIDs)
		if err != nil {
			return false, terror.New(err[0], "on track action")
		}

		cartons = _cartons
		cartonScanTimes = input.CartonScanTimes
	}

	for pi, palletID := range input.PalletIDs {
		palletUUID, err := uuid.FromString(palletID)
		if err != nil {
			return false, terror.New(err, "on track action")
		}

		_cartons, err := r.CartonStore.GetManyByPalletID(palletUUID)
		if err != nil {
			return false, terror.New(err, "on track action")
		}

		cartons = append(cartons, _cartons...)
		for i := 0; i < len(_cartons); i++ {
			if input.PalletScanTimes == nil || input.PalletScanTimes[pi] == nil {
				continue
			}
			cartonScanTimes = append(cartonScanTimes, input.PalletScanTimes[pi])
		}
	}
	for ci, containerID := range input.ContainerIDs {
		containerUUID, err := uuid.FromString(containerID)
		if err != nil {
			return false, terror.New(err, "on track action")
		}

		_cartons, err := r.CartonStore.GetManyByContainerID(containerUUID)
		if err != nil {
			return false, terror.New(err, "on track action")
		}

		cartons = append(cartons, _cartons...)
		for i := 0; i < len(_cartons); i++ {
			if input.ContainerScanTimes == nil || input.ContainerScanTimes[ci] == nil {
				continue
			}
			cartonScanTimes = append(cartonScanTimes, input.ContainerScanTimes[ci])
		}
	}

	// Setup transaction for extended info
	t := &db.Transaction{}
	if input.Memo != nil {
		t.Memo = *input.Memo
	}
	if input.LocationGeohash != nil {
		t.LocationGeohash = *input.LocationGeohash
	}
	if input.LocationName != nil {
		t.LocationName = *input.LocationName
	}

	// start transaction
	// all tx are the same, so just using any store's begintx to avoid reimplement (being lazy)
	tx, err := r.ManifestStore.BeginTransaction()
	if err != nil {
		return false, terror.New(err, "begin db tx")
	}
	defer tx.Rollback()

	// Attach transaction to cartons/products
	updatedCartons := make(map[string]bool)
	updatedProducts := make(map[string]bool)
	for pi, product := range products {
		// prevent products being updated multiple times in one action
		_, updated := updatedProducts[product.ID]
		if updated {
			continue
		}
		updatedProducts[product.ID] = true

		scanTime := time.Now()
		if pi < len(productScanTimes) {
			scanTime = *productScanTimes[pi]
		}

		// Create transaction using setup transaction above
		tt := *t
		tt.ScannedAt = null.TimeFrom(scanTime)
		_, err = r.TransactionStore.InsertByProduct(product, action, user, createdByName, &tt, tx)
		if err != nil {
			return false, terror.New(err, "on track action")
		}
	}

	for ci, carton := range cartons {
		// prevent cartons being updated multiple times in one action
		_, updated := updatedCartons[carton.ID]
		if updated {
			continue
		}
		updatedCartons[carton.ID] = true

		// get scan time
		scanTime := time.Now()
		if ci < len(cartonScanTimes) {
			scanTime = *cartonScanTimes[ci]
		}

		// get photos
		cartonPhotoBlobID := null.NewString("", false)
		productPhotoBlobID := null.NewString("", false)
		if input.CartonPhotoBlobIDs != nil && len(input.CartonPhotoBlobIDs) > ci && input.CartonPhotoBlobIDs[ci] != "" {
			cartonPhotoBlobID = null.StringFrom(input.CartonPhotoBlobIDs[ci])
		}
		if input.ProductPhotoBlobIDs != nil && len(input.ProductPhotoBlobIDs) > ci && input.ProductPhotoBlobIDs[ci] != "" {
			productPhotoBlobID = null.StringFrom(input.ProductPhotoBlobIDs[ci])
		}

		// Attach transaction by carton (skipped if same user already did the same action)
		// todo, what is it mean skipped? *above
		tt := *t // dup details
		_, err := r.TransactionStore.InsertByCarton(
			carton,
			action,
			user,
			createdByName,
			scanTime,
			cartonPhotoBlobID,
			productPhotoBlobID,
			&tt,
			tx,
		)
		if err != nil {
			return false, terror.New(err, "on track action")
		}

		cartonUUID, err := uuid.FromString(carton.ID)
		if err != nil {
			return false, terror.New(err, "on track action")
		}

		// Attach transaction to products in carton
		_products, err := r.ProductStore.GetManyByCartonID(cartonUUID)
		if err != nil {
			return false, terror.New(err, "on track action")
		}
		for pi, product := range _products {
			// prevent products being updated multiple times in one action
			_, updated := updatedProducts[product.ID]
			if updated {
				continue
			}
			updatedProducts[product.ID] = true

			scanTime := time.Now()
			if pi < len(productScanTimes) {
				scanTime = *productScanTimes[pi]
			}

			t := &db.Transaction{
				ScannedAt: null.TimeFrom(scanTime),
			}
			_, err = r.TransactionStore.InsertByProduct(product, action, user, createdByName, t)
			if err != nil {
				return false, terror.New(err, "on track action")
			}
		}
	}

	// commit to db
	err = tx.Commit()
	if err != nil {
		return false, terror.New(err, "commit create order")
	}

	return true, nil
}

func (r *mutationResolver) CopyCartonTransactionsToProduct(ctx context.Context, cartonUUID uuid.UUID, product *db.Product) error {
	// start transaction
	// all tx are the same, so just using any store's begintx to avoid reimplement (being lazy)
	tx, err := r.ManifestStore.BeginTransaction()
	if err != nil {
		return terror.New(err, "begin db tx")
	}
	defer tx.Rollback()

	transactions, err := r.TransactionStore.GetByCartonID(cartonUUID)
	if err != nil {
		return terror.New(err, "copy transactions to product")
	}

	err = r.TransactionStore.AttachManyToProduct(transactions, product, tx)
	if err != nil {
		return terror.New(err, "copy transactions to product")
	}

	// commit to db
	err = tx.Commit()
	if err != nil {
		return terror.New(err, "commit CopyCartonTransactionsToProduct")
	}

	return nil
}

func (r *mutationResolver) FlushPendingTransactions(ctx context.Context) (bool, error) {
	success, err := FlushPendingTransactions(ctx, r.TransactionStore, r.ManifestStore, r.Blk)
	if err != nil {
		return false, terror.New(err, "")
	}

	// Track user activity
	if success {
		r.RecordUserActivity(ctx, "Blockchain Commit", graphql.ObjectTypeBlockchain, nil, nil)
	}

	// Reset timer
	r.SystemTicker.Reset(ctx)

	// Check balance
	err = r.Blk.CheckBalance(ctx)
	if err != nil {
		fmt.Println("Failed to check current ETH account balance")
		terror.Echo(err)
	}

	return success, nil
}

// first line of the manifest, use for json data
type manifestHeader struct {
	ContractAddress  string `json:"contract"`         // smart contract address
	TransactionNonce int    `json:"txNonce"`          // blockchain transaction nonce
	PreviousManifest string `json:"previousManifest"` // previous manifest it refer to
}

// FlushPendingTransactions commits all pending transactions to the blockchain smart contract
func FlushPendingTransactions(ctx context.Context, TransactionStore TransactionStorer, ManifestStore ManifestStorer, Blk *blockchain.Service) (bool, error) {
	var err error

	// sanity check
	if ctx == nil {
		return false, terror.New(terror.ErrDataBlank, "ctx is nil")
	}
	if TransactionStore == nil {
		return false, terror.New(terror.ErrDataBlank, "txStore is nil")
	}
	if ManifestStore == nil {
		return false, terror.New(terror.ErrDataBlank, "maniStore is nil")
	}
	if Blk == nil {
		return false, terror.New(terror.ErrDataBlank, "blk is nil")
	}

	// lock, so no multiple flush running
	Blk.Mux.Lock()
	defer Blk.Mux.Unlock()

	// Get unfinished manifest, which shouldnt be. or wait until previous manifest is in the blockchain
	unfinishedManifests, err := ManifestStore.AllUnfinished()
	if err != nil {
		return false, terror.New(err, "get unfinished manifests")
	}
	if len(unfinishedManifests) > 0 {
		return false, terror.New(fmt.Errorf("unfinished manifest (busy with previous commit or still not confirmed on blockchain)"), "")
	}

	// Get pending transactions
	pendingTransactions, err := TransactionStore.AllPending()
	if err != nil {
		return false, terror.New(err, "get pending transactions")
	}
	if len(pendingTransactions) == 0 {
		// no pending transactions, stop
		return true, nil
	}

	// Get blockchain contract
	settings, err := TransactionStore.GetSettings()
	if err != nil {
		return false, terror.New(err, "get settings (2)")
	}

	contract, err := Blk.GetContract(common.HexToAddress(settings.SmartContractAddress))
	if err != nil {
		return false, terror.New(err, "get blockchain contract")
	}

	// uint64
	nonce, err := Blk.Chain.PendingNonceAt(ctx, Blk.Auth.From)
	if err != nil {
		return false, terror.New(err, "get nonce")
	}

	// start transaction
	tx, err := ManifestStore.BeginTransaction()
	if err != nil {
		return false, terror.New(err, "get transaction")
	}
	defer tx.Rollback()

	// find last manifest
	prevManifest, err := db.Manifests(
		db.ManifestWhere.Archived.EQ(false),
		qm.OrderBy(db.ManifestColumns.TransactionNonce+" DESC"),
	).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return false, terror.New(err, "find last manifest")
	}
	// sanity check
	var prevMerkleRootSha256 string // first one will be nil
	if prevManifest == nil {
		// its ok, can be blank if its first one
	} else {
		// exist
		if !prevManifest.MerkleRootSha256.Valid {
			return false, terror.New(fmt.Errorf("previous manifest is not valid"), "")
		}
		if prevManifest.MerkleRootSha256.IsZero() {
			return false, terror.New(fmt.Errorf("previous manifest is zero"), "")
		}
		if !prevManifest.Confirmed {
			return false, terror.New(fmt.Errorf("previous manifest is not confirmed"), "")
		}
		prevMerkleRootSha256 = prevManifest.MerkleRootSha256.String
	}

	// insert temporary manifest
	manifest := &db.Manifest{
		ContractAddress:  settings.SmartContractAddress,
		TransactionNonce: int(nonce),
		Pending:          true,
		Confirmed:        false,
	}
	err = manifest.Insert(tx, boil.Infer())
	if err != nil {
		return false, terror.New(err, "manifest insert")
	}

	// generate merkle root, sha256( byte( []sha256 ) )
	bsPtxSha256s := []byte{}
	// lines for compiled manifest text
	manifestLines := []string{}
	// add manifest header
	mheader := &manifestHeader{
		ContractAddress:  settings.SmartContractAddress,
		TransactionNonce: int(nonce),
		PreviousManifest: prevMerkleRootSha256,
	}
	jheader, err := json.Marshal(mheader)
	if err != nil {
		return false, terror.New(err, "")
	}
	jheaderSum := sha256.Sum256(jheader)
	manifestLines = append(manifestLines, fmt.Sprintf("%064x %s", jheaderSum, jheader))
	// add manifest header sha256
	bsPtxSha256s = append(bsPtxSha256s, jheaderSum[:]...)
	// add each pending transaction manifest line and data
	for _, pendingTransaction := range pendingTransactions {
		b := []byte{}
		// calc and save sha256 if not calc for manifest line json
		if pendingTransaction.ManifestLineSha256.IsZero() {
			b32 := sha256.Sum256([]byte(pendingTransaction.ManifestLineJSON.String))
			b = b32[:]
			// add and update pending transaction
			pendingTransaction.ManifestLineSha256 = null.StringFrom(fmt.Sprintf("%064x", b))
			_, err = pendingTransaction.Update(tx, boil.Infer())
			if err != nil {
				return false, terror.New(err, "add transaction json sha256")
			}
		}

		// copy sha256 in bytes
		b, err = helpers.HexToBytes(pendingTransaction.ManifestLineSha256.String)
		if err != nil {
			return false, terror.New(err, "")
		}
		// append sha256 bytes
		bsPtxSha256s = append(bsPtxSha256s, b...)

		// add manifest line
		line := fmt.Sprintf("%064x", b)
		manifestLines = append(manifestLines, line)

		// add manifest id
		pendingTransaction.ManifestID = null.StringFrom(manifest.ID)
		_, err = pendingTransaction.Update(tx, boil.Infer())
		if err != nil {
			return false, terror.New(err, "")
		}
	}
	// sha256 the sha256s
	merkleRoot := sha256.Sum256(bsPtxSha256s)

	// record merkle root
	manifest.MerkleRootSha256 = null.StringFrom(fmt.Sprintf("%x", merkleRoot[:]))

	// build manifest text
	manifest.CompiledText = null.BytesFrom([]byte(strings.Join(manifestLines, "\n")))

	// save data into db
	_, err = manifest.Update(tx, boil.Infer())
	if err != nil {
		return false, terror.New(err, "manifest insert/update")
	}

	// finish transaction
	err = tx.Commit()
	if err != nil {
		return false, terror.New(err, "transaction commit")
	}

	// Update Gas
	err = Blk.UpdateGas(ctx)
	if err != nil {
		return false, terror.New(err, "update gas")
	}

	// Commit manifest merkleRoot to smart contract
	var biMR = new(big.Int).SetBytes(merkleRoot[:])
	bcTx, err := contract.CommitBatch(Blk.Auth, biMR)
	if err != nil {
		return false, terror.New(err, "commit manifest to smart contract")
	}

	// start transaction 2nd time
	tx, err = ManifestStore.BeginTransaction()
	if err != nil {
		return false, terror.New(err, "get transaction 2")
	}

	// save blockchain transaction hash
	bcTxHash := null.StringFrom(bcTx.Hash().Hex())
	manifest.TransactionHash = bcTxHash

	// unset flag
	manifest.Pending = false

	// Update manifest
	_, err = ManifestStore.Update(manifest, tx)
	if err != nil {
		return false, terror.New(err, "update manifest")
	}

	// update transaction again
	for _, pendingTransaction := range pendingTransactions {
		pendingTransaction.TransactionHash = bcTxHash
		_, err = pendingTransaction.Update(tx, boil.Whitelist(db.TransactionColumns.TransactionHash))
		if err != nil {
			return false, terror.New(err, "update transactions")
		}
	}

	// finish transaction
	err = tx.Commit()
	if err != nil {
		return false, terror.New(err, "transaction commit")
	}

	return true, nil
}
