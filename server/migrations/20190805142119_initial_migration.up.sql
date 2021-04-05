BEGIN;
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
    wallet_points integer NOT NULL DEFAULT 0,
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
    code text UNIQUE NOT NULL,
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
    brand text NOT NULL,
    ingredients text NOT NULL,
    description text NOT NULL,
    weight float NOT NULL DEFAULT 0,
    weight_unit text NOT NULL,
    price float NOT NULL DEFAULT 0,
    purchase_points int NOT NULL DEFAULT 0,
    loyalty_points int NOT NULL DEFAULT 0,
    currency text NOT NULL,
    is_beef boolean NOT NULL DEFAULT FALSE,
    is_point_bound boolean NOT NULL DEFAULT FALSE,
    is_app_bound boolean NOT NULL DEFAULT FALSE,
    is_app_sku boolean NOT NULL DEFAULT FALSE,
    brand_logo_blob_id uuid REFERENCES blobs (id),
    gif_blob_id uuid REFERENCES blobs (id),
    master_plan_blob_id uuid REFERENCES blobs (id),
    video_blob_id uuid REFERENCES blobs (id),
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
-- Retail Links
CREATE TABLE retail_links (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    name text NOT NULL,
    url text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Categories
CREATE TABLE categories (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE product_categories (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    sku_id uuid NOT NULL REFERENCES stock_keeping_units (id),
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Tasks Specifications
CREATE TABLE tasks (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    loyalty_points int NOT NULL DEFAULT 0,
    is_time_bound boolean NOT NULL DEFAULT FALSE,
    is_people_bound boolean NOT NULL DEFAULT FALSE,
    is_product_relevant boolean NOT NULL DEFAULT FALSE,
    finish_date timestamptz,
    maximum_people int NOT NULL DEFAULT 0,
    sku_id uuid REFERENCES stock_keeping_units (id),
    banner_photo_blob_id uuid REFERENCES blobs (id),
    brand_logo_blob_id uuid REFERENCES blobs (id),
    is_final boolean NOT NULL DEFAULT FALSE,
    archived boolean NOT NULL DEFAULT FALSE,
    archived_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE subtasks (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    task_id uuid REFERENCES tasks (id),
    title text NOT NULL,
    description text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE user_tasks (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    task_id uuid REFERENCES tasks (id),
    user_id uuid NOT NULL REFERENCES users (id),
    status text NOT NULL,
    is_complete boolean NOT NULL DEFAULT FALSE,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE user_subtasks (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    subtask_id uuid REFERENCES subtasks (id),
    user_task_id uuid REFERENCES user_tasks (id),
    status text NOT NULL,
    is_complete boolean NOT NULL DEFAULT FALSE,
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
    latitude float NOT NULL,
    longitude float NOT NULL,
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
    is_point_bound boolean NOT NULL DEFAULT FALSE,
    is_app_bound boolean NOT NULL DEFAULT FALSE,
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
    register_id uuid NOT NULL DEFAULT gen_random_uuid (),
    is_beef boolean NOT NULL DEFAULT FALSE,
    is_point_bound boolean NOT NULL DEFAULT FALSE,
    is_app_bound boolean NOT NULL DEFAULT FALSE,
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
-- Purchase
CREATE TABLE user_purchase_activities (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    code text UNIQUE NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id),
    product_id uuid REFERENCES products (id),
    loyalty_points int NOT NULL, -- sku points + bonus points
    message text NOT NULL DEFAULT '',
    transaction_hash text,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
-- Wallet Transaction
CREATE TABLE wallet_transaction (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL REFERENCES users (id),
    loyalty_points int NOT NULL,
    message text NOT NULL DEFAULT '',
    is_credit boolean NOT NULL DEFAULT FALSE,
    transaction_hash text,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
COMMIT;