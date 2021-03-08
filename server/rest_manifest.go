package genesis

import (
	"bytes"
	"fmt"
	"genesis/blockchain"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/ninja-software/terror"
)

var posmbaTemplate = template.Must(template.ParseFiles(`rest_manifest_tpl_list.html`))

type posmbaData struct {
	EtherscanHost string
	Transactions  []posmbaTx
}
type posmbaTx struct {
	Nonce      int
	MerkleRoot string
	Hash       string
	Confirmed  bool
}

// ProofOfSteakManifestBasicAll provide basic list of all manifests
func ProofOfSteakManifestBasicAll(
	etherscanHost string,
	TransactionStore TransactionStorer,
	ManifestStore ManifestStorer,
	TrackActionStore TrackActionStorer,
	Blk *blockchain.Service,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		hashes := []posmbaTx{}
		manifests, err := ManifestStore.AllUnarchived()
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, terror.New(err, "get confirmed manifest"))
			return
		}
		for _, manifest := range manifests {
			hashes = append(hashes, posmbaTx{
				Nonce:      manifest.TransactionNonce,
				MerkleRoot: manifest.MerkleRootSha256.String,
				Hash:       manifest.TransactionHash.String,
				Confirmed:  manifest.Confirmed,
			})
		}

		data := posmbaData{
			EtherscanHost: etherscanHost,
			Transactions:  hashes,
		}

		buf := new(bytes.Buffer)
		err = posmbaTemplate.Execute(buf, data)
		if err != nil {
			restWriteError(r.Context(), w, http.StatusInternalServerError, terror.New(err, "exec template"))
			return
		}

		w.Write(buf.Bytes())
	}

	return fn
}

// ProofOfSteakManifestByMerkleRoot provide individual manifest text
func ProofOfSteakManifestByMerkleRoot(
	etherscanHost string,
	TransactionStore TransactionStorer,
	ManifestStore ManifestStorer,
	TrackActionStore TrackActionStorer,
	Blk *blockchain.Service,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		// get blockchain transaction id
		merkleRoot := chi.URLParam(r, "merkleRoot")

		// sanity check
		if merkleRoot == "" {
			err = terror.New(fmt.Errorf("mr not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get manifest
		manifest, err := ManifestStore.GetByBlockchainByMerkleRootHash(merkleRoot)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to find manifest"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(manifest.CompiledText.Bytes)
	}

	return fn
}

// ProofOfSteakManifestByHash provide individual manifest text
func ProofOfSteakManifestByHash(
	etherscanHost string,
	TransactionStore TransactionStorer,
	ManifestStore ManifestStorer,
	TrackActionStore TrackActionStorer,
	Blk *blockchain.Service,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		// get blockchain transaction id
		bcTxID := chi.URLParam(r, "txID")

		// sanity check
		if bcTxID == "" {
			err = terror.New(fmt.Errorf("txID not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get manifest
		manifest, err := ManifestStore.GetByBlockchainByTransactionHash(bcTxID)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to find manifest"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(manifest.CompiledText.Bytes)
	}

	return fn
}

// ProofOfSteakManifestByLine provide individual manifest text, search by line
func ProofOfSteakManifestByLine(
	etherscanHost string,
	TransactionStore TransactionStorer,
	ManifestStore ManifestStorer,
	TrackActionStore TrackActionStorer,
	Blk *blockchain.Service,
) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		// get manifest line hash
		lineHash := chi.URLParam(r, "lineHash")

		// sanity check
		if lineHash == "" {
			err = terror.New(fmt.Errorf("lineHash not provided"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		// Get manifest
		manifest, err := ManifestStore.GetByLineHash(lineHash)
		if err != nil {
			err = terror.New(fmt.Errorf("failed to find manifest"), "")
			restWriteError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(manifest.CompiledText.Bytes)
	}

	return fn
}
