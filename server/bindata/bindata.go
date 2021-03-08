// Code generated by go-bindata. DO NOT EDIT.
// sources:
// migrations/20190805142119_initial_migration.up.sql (17.054kB)

package bindata

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __20190805142119_initial_migrationUpSql = []byte(`BEGIN;
DROP TYPE IF EXISTS CONTENT_TYPE;
CREATE TYPE CONTENT_TYPE AS ENUM (
    'INFO',
    'URL'
);
-- Settings
CREATE TABLE settings (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    -- blockchain smart contract address
    smart_contract_address text UNIQUE NOT NULL DEFAULT ''
);
-- Blobs
CREATE TABLE blobs (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    file_name text NOT NULL,
    mime_type text NOT NULL,
    file_size_bytes bigint NOT NULL,
    extension TEXT NOT NULL,
    file bytea NOT NULL,
    views integer NOT NULL DEFAULT 0,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
--  Users
CREATE TABLE organisations (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    name text NOT NULL,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE prospects (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    email text UNIQUE NOT NULL,
    first_name text,
    last_name text,
    onboarding_complete boolean NOT NULL DEFAULT FALSE
);
CREATE TABLE roles (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    name text UNIQUE NOT NULL,
    permissions text[] NOT NULL,
    tier integer NOT NULL DEFAULT 3, -- users can never edit another user with a tier <= to their own (SUPERADMIN = 1, ORGADMIN = 2)
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    organisation_id uuid REFERENCES organisations (id),
    role_id uuid NOT NULL REFERENCES roles (id),
    email text UNIQUE,
    first_name text,
    last_name text,
    affiliate_org text,
    referral_code text,
    mobile_phone text UNIQUE,
    mobile_verified boolean NOT NULL DEFAULT FALSE,
    wechat_id text UNIQUE,
    verified boolean NOT NULL DEFAULT FALSE,
    verify_token text NOT NULL DEFAULT gen_random_uuid (),
    verify_token_expires timestamptz NOT NULL DEFAULT NOW(),
    reset_token text NOT NULL DEFAULT gen_random_uuid (),
    reset_token_expires timestamptz NOT NULL DEFAULT NOW(),
    password_hash text NOT NULL,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE referrals (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id),
    referred_by_id uuid REFERENCES users (id),
    is_redemmed boolean NOT NULL DEFAULT FALSE,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE issued_tokens (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id),
    device text NOT NULL,
    token_created timestamptz NOT NULL DEFAULT NOW(),
    token_expires timestamptz NOT NULL,
    blacklisted boolean NOT NULL DEFAULT FALSE
);
-- SKUS
CREATE TABLE stock_keeping_units (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    name text NOT NULL,
    code text UNIQUE NOT NULL,
    description text NOT NULL,
    weight int NOT NULL DEFAULT 0,
    price int NOT NULL DEFAULT 0,
    currency text NOT NULL,
    is_beef boolean NOT NULL DEFAULT FALSE,
    is_point_sku boolean NOT NULL DEFAULT FALSE,
    is_app_sku boolean NOT NULL DEFAULT FALSE,
    master_plan_blob_id uuid REFERENCES blobs (id),
    video_blob_id uuid REFERENCES blobs (id),
    loyalty_points int NOT NULL DEFAULT 0,
    clone_parent_id uuid REFERENCES stock_keeping_units (id),
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
CREATE TABLE stock_keeping_unit_content (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    content_type CONTENT_TYPE NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE stock_keeping_unit_photos (
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    photo_id uuid NOT NULL REFERENCES blobs (id),
    sort_index integer NOT NULL DEFAULT 0,
    PRIMARY KEY (sku_id, photo_id)
);
-- Categories
CREATE TABLE categories (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    name text NOT NULL UNIQUE,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE product_categories (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    name text NOT NULL UNIQUE,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Tasks Specifications
CREATE TABLE tasks (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    title text NOT NULL,
    description text NOT NULL,
    loyalty_points int NOT NULL DEFAULT 0,
    is_time_bound boolean NOT NULL DEFAULT FALSE,
    is_people_bound boolean NOT NULL DEFAULT FALSE,
    is_product_relevant boolean NOT NULL DEFAULT FALSE,
    finish_date timestamptz,
    maximum_people int NOT NULL DEFAULT 0,
    sku_id uuid REFERENCES stock_keeping_units (id),
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE task_steps (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    task_id uuid REFERENCES tasks (id),
    name text NOT NULL,
    description text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE user_tasks (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    task_id uuid REFERENCES tasks (id),
    user_id uuid NOT NULL REFERENCES users (id),
    status text NOT NULL,
    is_complete boolean NOT NULL DEFAULT FALSE,
    is_active boolean NOT NULL DEFAULT FALSE,
    loyalty_points int NOT NULL DEFAULT 0,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE user_task_steps (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    user_task_id uuid REFERENCES user_tasks (id),
    name text NOT NULL,
    description text NOT NULL,
    status text NOT NULL,
    is_complete boolean NOT NULL DEFAULT FALSE,
    is_active boolean NOT NULL DEFAULT FALSE,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Procurement Contracts / Livestock Specification
CREATE TABLE contracts (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    supplier_name text NOT NULL,
    date_signed timestamptz,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
-- "Orders" used for package generation
CREATE TABLE orders (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
-- Distributors - attached to cartons/products (distributor code is a parameter on QR2 scan)
CREATE TABLE distributors (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    name text NOT NULL,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
-- Inventory
CREATE TABLE containers (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    description text NOT NULL DEFAULT '',
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
CREATE TABLE pallets (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    description text NOT NULL DEFAULT '',
    container_id uuid REFERENCES containers (id),
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
CREATE TABLE cartons (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    weight text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    meat_type text NOT NULL DEFAULT '',
    spreadsheet_link text NOT NULL DEFAULT '',
    processed_at timestamptz,
    pallet_id uuid REFERENCES pallets (id),
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
-- Products
CREATE TABLE products (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    description text NOT NULL DEFAULT '',
    weight int NOT NULL DEFAULT 0,
    price int NOT NULL DEFAULT 0,
    currency text NOT NULL,
    register_id uuid NOT NULL DEFAULT gen_random_uuid (),
    is_beef boolean NOT NULL DEFAULT FALSE,
    is_point_product boolean NOT NULL DEFAULT FALSE,
    is_app_product boolean NOT NULL DEFAULT FALSE,
    loyalty_points int NOT NULL DEFAULT 0,
    loyalty_points_expire timestamptz NOT NULL DEFAULT NOW(),
    sku_id uuid REFERENCES stock_keeping_units (id),
    carton_id uuid REFERENCES cartons (id),
    order_id uuid REFERENCES orders (id),
    contract_id uuid REFERENCES contracts (id),
    distributor_id uuid REFERENCES distributors (id),
    -- require this to close register, regenerate each time user scan QR2
    close_register_id uuid,
    transaction_hash text NOT NULL,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL REFERENCES users (id)
);
-- Tracking
CREATE TABLE track_actions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    require_photos boolean[] NOT NULL DEFAULT ARRAY[FALSE, FALSE],
    code text UNIQUE NOT NULL,
    name text NOT NULL DEFAULT '',
    name_chinese text NOT NULL DEFAULT '',
    private boolean NOT NULL DEFAULT FALSE, -- whether consumers can see this action on the product view page
    system boolean NOT NULL DEFAULT FALSE, -- default track action that is only logged by system (eg: moved to carton, moved to pallet, etc)
    blockchain boolean NOT NULL DEFAULT FALSE, -- whether or not transactions made with this track action are commited to the blockchain
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid REFERENCES users (id)
);
-- Track actions a role is allowed to do
CREATE TABLE role_track_actions (
    role_id uuid NOT NULL REFERENCES roles (id),
    track_action_id uuid NOT NULL REFERENCES track_actions (id),
    PRIMARY KEY (role_id, track_action_id)
);
-- manifest equate to a single block to blockchain, implementing blockchain in a blockchain
CREATE TABLE manifests (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    -- all refer to blockchain
    contract_address text NOT NULL,
    transaction_nonce int NOT NULL, -- blockchain transaction nonce
    transaction_hash text UNIQUE,
    confirmed boolean NOT NULL DEFAULT FALSE, -- confirmed it is in the blockchain
    merkle_root_sha256 text, -- uint256 sha256 data for blockchain. sha256( byte(sha256({"contract": "address", "nonce": "nonce"})) + byte( []transactions.manifest_sha256 ) ), refer to manifest example file
    -- end of blockchain

    compiled_text bytea, -- compiled manifest text file, do not change after compiled, because it is published to the world
    pending boolean NOT NULL DEFAULT TRUE, -- is it being processed
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Transactions, when record is created, it is set in stone, only can update few sections
CREATE TABLE transactions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    -- write once
    track_action_id uuid NOT NULL REFERENCES track_actions (id),
    memo text,
    product_id uuid REFERENCES products (id),
    carton_id uuid REFERENCES cartons (id),
    scanned_at timestamptz,
    location_geohash text,
    location_name text,
    -- manifest data, if blank, do not change after inputed
    manifest_line_json text, -- individual manifest line, not using JSONB type because it needs to be the same permanently after set
    manifest_line_sha256 text, -- in hex representation. indexed so can search
    manifest_id uuid REFERENCES manifests (id), -- if filled, it means it is been registered to manifests table
    transaction_hash text, -- blockchain smart contract transaction hash, updated last async. if pre-set to "-", it will not be processed and published
    product_photo_blob_id uuid REFERENCES blobs (id), -- can change without problem, as long as not save into manifest
    carton_photo_blob_id uuid REFERENCES blobs (id), -- ditto
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_by_id uuid REFERENCES users (id),
    created_by_name text NOT NULL -- name we want to capture in full, if created_by_id is not available
);
-- prevent accidental creating multiple duplicate transactions when product inherit carton history
ALTER TABLE transactions
    ADD CONSTRAINT unique_transaction UNIQUE (track_action_id, product_id, created_at);
-- make it easier to find latest transaction for product
CREATE OR REPLACE VIEW product_latest_transactions AS SELECT DISTINCT ON (product_id)
    t.*
FROM
    transactions t
WHERE
    product_id IS NOT NULL
    AND track_action_id IS NOT NULL
ORDER BY
    product_id,
    created_at DESC;
-- make it easier to find latest transaction for carton
CREATE OR REPLACE VIEW carton_latest_transactions AS SELECT DISTINCT ON (carton_id)
    t.*
FROM
    transactions t
WHERE
    carton_id IS NOT NULL
    AND track_action_id IS NOT NULL
ORDER BY
    carton_id,
    created_at DESC;
-- make it easier to find latest transaction for pallet
CREATE OR REPLACE VIEW pallet_latest_transactions AS SELECT DISTINCT ON (c2.pallet_id)
    t.*,
    p2.id AS pallet_id
FROM
    carton_latest_transactions t
    INNER JOIN LATERAL (
        SELECT
            c1.id,
            pallet_id
        FROM
            cartons c1) c2 ON t.carton_id = c2.id
    INNER JOIN LATERAL (
        SELECT
            p1.id
        FROM
            pallets p1) p2 ON c2.pallet_id = p2.id
WHERE
    carton_id IS NOT NULL
ORDER BY
    c2.pallet_id,
    created_at DESC;
-- User Activity Tracking
CREATE TABLE user_activities (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id),
    action text NOT NULL,
    object_id text, -- uuid
    object_code text, -- used for user activity list (for ease of reading and links)
    object_type text NOT NULL, -- enum defined in gql
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Loyalty
CREATE TABLE user_loyalty_activities (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id),
    product_id uuid REFERENCES products (id),
    amount int NOT NULL, -- sku points + bonus points
    bonus int NOT NULL DEFAULT 0, -- note only, not added in total calculation
    message text NOT NULL DEFAULT '',
    transaction_hash text,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
COMMIT;

`)

func _20190805142119_initial_migrationUpSqlBytes() ([]byte, error) {
	return __20190805142119_initial_migrationUpSql, nil
}

func _20190805142119_initial_migrationUpSql() (*asset, error) {
	bytes, err := _20190805142119_initial_migrationUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190805142119_initial_migration.up.sql", size: 17054, mode: os.FileMode(0644), modTime: time.Unix(1615011067, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x29, 0x1e, 0x8f, 0x7c, 0x72, 0x88, 0x7f, 0x27, 0x2c, 0x80, 0xc2, 0x7f, 0x2d, 0x7f, 0x2, 0x73, 0xba, 0x1f, 0x7b, 0x65, 0xc6, 0x11, 0x7d, 0xa7, 0xc6, 0xcf, 0x13, 0xd9, 0x12, 0x97, 0x6d, 0x4a}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"20190805142119_initial_migration.up.sql": _20190805142119_initial_migrationUpSql,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"20190805142119_initial_migration.up.sql": {_20190805142119_initial_migrationUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
