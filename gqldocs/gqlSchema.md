directive @hasPerm(p: Perm!) on FIELD_DEFINITION
directive @hasAnyPerm(p: [Perm!]) on FIELD_DEFINITION
directive @hasAllPerms(p: [Perm!]) on FIELD_DEFINITION
enum Action {
  Archive
  Unarchive
  SetSKU
  SetOrder
  SetDistributor
  SetContract
  SetCarton
  SetPallet
  SetContainer
  DetachFromSKU
  DetachFromOrder
  DetachFromDistributor
  DetachFromContract
  DetachFromCarton
  DetachFromPallet
  DetachFromContainer
  SetBonusLoyaltyPoints
  InheritCartonHistory
}

input BatchActionInput {
  str: NullString
  no: NullInt
  dateTime: NullTime
  bool: NullBool
}

type Blob {
  id: String!
  file_url: String
  file_name: String!
  mime_type: String!
  file_size_bytes: Int!
  views: Int!
}

type Carton {
  id: ID!
  code: String!
  weight: String!
  processedAt: NullTime
  description: String!
  meatType: String!
  archived: Boolean!
  createdAt: Time!
  pallet: Pallet
  sku: SKU
  order: Order
  distributor: Distributor
  transactions: [Transaction!]!
  latestTrackAction: LatestTransactionInfo
  productCount: Int!
  spreadsheetLink: String!
}

type CartonResult {
  cartons: [Carton!]!
  total: Int!
}

type ConsumersResult {
  consumers: [User!]!
  total: Int!
}

type Container {
  id: ID!
  code: String!
  archived: Boolean!
  createdAt: Time!
  description: String!
  palletCount: Int!
}

type ContainerResult {
  containers: [Container!]!
  total: Int!
}

type Contract {
  id: ID!
  code: String!
  name: String!
  description: String!
  supplierName: String!
  dateSigned: NullTime
  archived: Boolean!
  createdAt: Time!
}

type ContractResult {
  contracts: [Contract!]!
  total: Int!
  createdAt: Time!
}

input CreateCarton {
  palletID: NullString
  quantity: Int!
  description: String!
}

input CreateContainer {
  quantity: Int!
  description: String!
}

input CreateOrder {
  contractID: NullString
  skuID: NullString
  quantity: Int!
}

input CreatePallet {
  containerID: NullString
  quantity: Int!
  description: String!
}

type Distributor {
  id: ID!
  code: String!
  name: String!
  archived: Boolean!
  createdAt: Time!
}

type DistributorResult {
  distributors: [Distributor!]!
  total: Int!
  createdAt: Time!
}

enum FilterOption {
  All
  Active
  Archived
  ProductWithoutOrder
  ProductWithoutCarton
  ProductWithoutSKU
  CartonWithoutPallet
  PalletWithoutContainer
  System
  Blockchain
  Pending
}

type GetObjectResponse {
  product: Product
  carton: Carton
  pallet: Pallet
  container: Container
}

input GetObjectsRequest {
  productIDs: [String!]!
  cartonIDs: [String!]!
  palletIDs: [String!]!
  containerIDs: [String!]!
}

type GetObjectsResponse {
  products: [Product!]!
  cartons: [Carton!]!
  pallets: [Pallet!]!
  containers: [Container!]!
}

type LatestTransactionInfo {
  name: String!
  createdAt: Time!
}

type Manifest {
  id: ID!
  transactionHash: NullString
  merkleRootSha256: NullString
}

type Mutation {
  RequestToken(input: RequestToken): String!
  fileUpload(file: Upload!): Blob!
  fileUploadMultiple(files: [Upload!]!): [Blob!]!
  deploySmartContract: Settings!
  cartonCreate(input: CreateCarton!): String!
  cartonUpdate(id: ID!, input: UpdateCarton!): Carton!
  cartonArchive(id: ID!): Carton!
  cartonUnarchive(id: ID!): Carton!
  cartonBatchAction(
    ids: [ID!]!
    action: Action!
    value: BatchActionInput
  ): Boolean!
  containerCreate(input: CreateContainer!): String!
  containerUpdate(id: ID!, input: UpdateContainer!): Container!
  containerArchive(id: ID!): Container!
  containerUnarchive(id: ID!): Container!
  containerBatchAction(
    ids: [ID!]!
    action: Action!
    value: BatchActionInput
  ): Boolean!
  contractCreate(input: UpdateContract!): Contract!
  contractUpdate(id: ID!, input: UpdateContract!): Contract!
  contractArchive(id: ID!): Contract!
  contractUnarchive(id: ID!): Contract!
  distributorCreate(input: UpdateDistributor!): Distributor!
  distributorUpdate(id: ID!, input: UpdateDistributor!): Distributor!
  distributorArchive(id: ID!): Distributor!
  distributorUnarchive(id: ID!): Distributor!
  orderCreate(input: CreateOrder!): Order!
  orderUpdate(id: ID!, input: UpdateOrder!): Order!
  orderArchive(id: ID!): Order!
  orderUnarchive(id: ID!): Order!
  orderBatchAction(
    ids: [ID!]!
    action: Action!
    value: BatchActionInput
  ): Boolean!
  palletCreate(input: CreatePallet!): String!
  palletUpdate(id: ID!, input: UpdatePallet!): Pallet!
  palletArchive(id: ID!): Pallet!
  palletUnarchive(id: ID!): Pallet!
  palletBatchAction(
    ids: [ID!]!
    action: Action!
    value: BatchActionInput
  ): Boolean!
  productCreate(input: UpdateProduct!): Product!
  productUpdate(id: ID!, input: UpdateProduct!): Product!
  productArchive(id: ID!): Product!
  productUnarchive(id: ID!): Product!
  productBatchAction(
    ids: [ID!]!
    action: Action!
    value: BatchActionInput
  ): Boolean!
  roleCreate(input: UpdateRole!): Role!
  roleUpdate(id: ID!, input: UpdateRole!): Role!
  roleArchive(id: ID!): Role!
  roleUnarchive(id: ID!): Role!
  skuCreate(input: UpdateSKU!): SKU!
  skuUpdate(id: ID!, input: UpdateSKU!): SKU!
  skuArchive(id: ID!): SKU!
  skuUnarchive(id: ID!): SKU!
  skuBatchAction(
    ids: [ID!]!
    action: Action!
    value: BatchActionInput
  ): Boolean!
  trackActionCreate(input: UpdateTrackAction!): TrackAction!
  trackActionUpdate(id: ID!, input: UpdateTrackAction!): TrackAction!
  trackActionArchive(id: ID!): TrackAction!
  trackActionUnarchive(id: ID!): TrackAction!
  recordTransaction(input: RecordTransactionInput!): Boolean!
  flushPendingTransactions: Boolean!
  changePassword(oldPassword: String!, password: String!): Boolean!
  changeDetails(input: UpdateUser!): User!
  userCreate(input: UpdateUser!): User!
  userUpdate(id: ID!, input: UpdateUser!): User!
  forgotPassword(email: String!, viaSMS: Boolean): Boolean!
  resetPassword(token: String!, password: String!, email: NullString): Boolean!
  resendEmailVerification(email: String!): Boolean!
  userArchive(id: ID!): User!
  userUnarchive(id: ID!): User!
}

scalar NullBool

scalar NullInt

scalar NullString

scalar NullTime

enum ObjectType {
  Self
  User
  Role
  Sku
  Contract
  Order
  Container
  Pallet
  Carton
  Distributor
  Product
  TrackAction
  Blob
  Blockchain
}

type Order {
  id: ID!
  code: String!
  archived: Boolean!
  createdAt: Time!
  sku: SKU
  productCount: Int!
}

type OrderResult {
  orders: [Order!]!
  total: Int!
}

type Organisation {
  id: ID!
  name: String!
  users: [User!]!
}

type PageInfo {
  startCursor: ID!
  endCursor: ID!
}

type Pallet {
  id: ID!
  code: String!
  archived: Boolean!
  createdAt: Time!
  description: String!
  container: Container
  latestTrackAction: LatestTransactionInfo
  cartonCount: Int!
}

type PalletResult {
  pallets: [Pallet!]!
  total: Int!
}

enum Perm {
  UserList
  UserCreate
  UserRead
  UserUpdate
  UserArchive
  UserUnarchive
  OrganisationList
  OrganisationCreate
  OrganisationRead
  OrganisationUpdate
  OrganisationArchive
  OrganisationUnarchive
  RoleList
  RoleCreate
  RoleRead
  RoleUpdate
  RoleArchive
  RoleUnarchive
  SKUList
  SKUCreate
  SKURead
  SKUUpdate
  SKUArchive
  SKUUnarchive
  ContainerList
  ContainerRead
  ContainerCreate
  ContainerUpdate
  ContainerArchive
  ContainerUnarchive
  PalletList
  PalletRead
  PalletCreate
  PalletUpdate
  PalletArchive
  PalletUnarchive
  CartonList
  CartonRead
  CartonCreate
  CartonUpdate
  CartonArchive
  CartonUnarchive
  ProductList
  ProductRead
  ProductCreate
  ProductUpdate
  ProductArchive
  ProductUnarchive
  OrderList
  OrderRead
  OrderCreate
  OrderUpdate
  OrderArchive
  OrderUnarchive
  TrackActionList
  TrackActionRead
  TrackActionCreate
  TrackActionUpdate
  TrackActionArchive
  TrackActionUnarchive
  ContractList
  ContractRead
  ContractCreate
  ContractUpdate
  ContractArchive
  ContractUnarchive
  DistributorList
  DistributorRead
  DistributorCreate
  DistributorUpdate
  DistributorArchive
  DistributorUnarchive
  ActivityListBlockchainActivity
  ActivityListUserActivity
  UseAdvancedMode
  UseAdminPortal
}

type Product {
  id: ID!
  code: String!
  registerID: String!
  loyaltyPoints: Int!
  loyaltyPointsExpire: Time!
  archived: Boolean!
  createdAt: Time!
  sku: SKU
  description: String!
  carton: Carton
  order: Order
  contract: Contract
  distributor: Distributor
  registered: Boolean!
  registeredBy: User
  transactions: [Transaction!]!
  latestTrackAction: LatestTransactionInfo
}

type ProductResult {
  products: [Product!]!
  total: Int!
}

type Query {
  settings: Settings!
  getTickerInfo: TickerInfo!
  getObject(id: String!): GetObjectResponse!
  getObjects(input: GetObjectsRequest!): GetObjectsResponse!
  cartons(
    search: SearchFilter!
    limit: Int!
    offset: Int!
    palletID: String
    trackActionID: String
  ): CartonResult!
  carton(code: String!): Carton!
  containers(search: SearchFilter!, limit: Int!, offset: Int!): ContainerResult!
  container(code: String!): Container!
  contracts(search: SearchFilter!, limit: Int!, offset: Int!): ContractResult!
  contract(code: String!): Contract!
  distributors(
    search: SearchFilter!
    limit: Int!
    offset: Int!
  ): DistributorResult!
  distributor(code: String!): Distributor!
  getLoyaltyActivity(userID: ID!): [UserLoyaltyActivity!]!
  orders(search: SearchFilter!, limit: Int!, offset: Int!): OrderResult!
  order(code: String!): Order!
  pallets(
    search: SearchFilter!
    limit: Int!
    offset: Int!
    containerID: String
    trackActionID: String
  ): PalletResult!
  pallet(code: String!): Pallet!
  products(
    search: SearchFilter!
    limit: Int!
    offset: Int!
    cartonID: String
    orderID: String
    skuID: String
    distributorID: String
    contractID: String
    trackActionID: String
  ): ProductResult!
  product(code: String!): Product!
  productByID(id: String!): Product!
  roles(
    search: SearchFilter!
    limit: Int!
    offset: Int!
    excludeSuper: Boolean!
  ): RolesResult!
  role(name: String!): Role!
  skus(search: SearchFilter!, limit: Int!, offset: Int!): SKUResult!
  sku(code: String!): SKU!
  skuByID(id: ID!): SKU!
  skuCloneTree(id: ID!): [SKUClone!]
  trackActions(
    search: SearchFilter!
    limit: Int!
    offset: Int!
  ): TrackActionResult!
  trackAction(code: String!): TrackAction!
  transactions(
    search: SearchFilter!
    limit: Int!
    offset: Int!
    productID: String
    cartonID: String
    trackActionID: String
  ): TransactionsResult!
  pendingTransactionsCount: Int!
  ethereumAccountAddress: String!
  ethereumAccountBalance: String!
  userActivities(
    search: SearchFilter!
    limit: Int!
    offset: Int!
    userID: String
  ): UserActivityResult!
  me: User!
  organisations: [Organisation!]!
  users(search: SearchFilter!, limit: Int!, offset: Int!): UsersResult!
  user(email: String, wechatID: String): User!
  consumers(search: SearchFilter!, limit: Int!, offset: Int!): ConsumersResult!
  verifyResetToken(token: String!, email: NullString): Boolean!
}

input RecordTransactionInput {
  trackActionCode: String!
  productIDs: [String!]
  cartonIDs: [String!]
  palletIDs: [String!]
  containerIDs: [String!]
  productScanTimes: [Time!]
  cartonScanTimes: [Time!]
  palletScanTimes: [Time!]
  containerScanTimes: [Time!]
  cartonPhotoBlobIDs: [String!]
  productPhotoBlobIDs: [String!]
  memo: NullString
  locationGeohash: NullString
  locationName: NullString
}

input RequestToken {
  email: String!
  password: String!
}

type Role {
  id: String
  name: String
  tier: Int!
  archived: Boolean!
  createdAt: Time!
  permissions: [Perm!]!
  trackActions: [TrackAction!]!
}

type RolesResult {
  roles: [Role!]!
  total: Int!
}

input SearchFilter {
  search: NullString
  filter: FilterOption
  sortBy: SortByOption
  sortDir: SortDir
}

type Settings {
  consumerHost: String!
  adminHost: String!
  etherscanHost: String!
  fieldappVersion: String!
  smartContractAddress: String!
}

type SKU {
  id: ID!
  name: String!
  code: String!
  description: String!
  isBeef: Boolean!
  loyaltyPoints: Int!
  archived: Boolean!
  createdAt: Time!
  hasClones: Boolean!
  cloneParentID: NullString
  masterPlan: Blob
  video: Blob
  urls: [SKUContent!]!
  productInfo: [SKUContent!]!
  photos: [Blob!]!
  productCount: Int!
}

type SKUClone {
  sku: SKU!
  depth: Int!
}

type SKUContent {
  title: String!
  content: String!
}

input SKUContentInput {
  title: String!
  content: String!
}

type SKUResult {
  skus: [SKU!]!
  total: Int!
}

enum SortByOption {
  DateCreated
  DateUpdated
  Alphabetical
}

enum SortDir {
  Ascending
  Descending
}

type TickerInfo {
  lastTick: Time!
  tickInterval: Int!
}

scalar Time

type TrackAction {
  id: ID!
  code: String!
  requirePhotos: [Boolean!]
  name: String!
  nameChinese: String!
  private: Boolean!
  system: Boolean!
  blockchain: Boolean!
  archived: Boolean!
  createdAt: Time!
}

type TrackActionResult {
  trackActions: [TrackAction!]!
  total: Int!
}

type Transaction {
  id: ID!
  transactionHash: NullString
  transactionPending: Boolean!
  manifestLineJson: NullString
  manifestLineSha256: NullString
  manifestID: NullString
  manifest: Manifest
  locationGeohash: NullString
  locationName: NullString
  action: TrackAction!
  memo: NullString
  createdAt: Time!
  createdBy: User
  createdByName: String!
  carton: Carton
  product: Product
  photos: TransactionPhotos
  scannedAt: NullTime
}

type TransactionPhotos {
  cartonPhoto: Blob
  productPhoto: Blob
}

type TransactionsResult {
  transactions: [Transaction!]!
  total: Int!
}

input UpdateCarton {
  code: NullString
  weight: NullString
  palletID: NullString
  processedAt: NullTime
  description: NullString
  meatType: NullString
}

input UpdateContainer {
  code: NullString
  description: NullString
}

input UpdateContract {
  name: NullString
  description: NullString
  supplierName: NullString
  dateSigned: NullTime
}

input UpdateDistributor {
  name: NullString
  code: NullString
}

input UpdateLoyaltyPoints {
  productIDs: [String!]
  cartonID: NullString
  palletID: NullString
  containerID: NullString
  loyaltyPoints: Int!
  loyaltyPointsExpire: Time!
}

input UpdateOrder {
  code: NullString
}

input UpdatePallet {
  code: NullString
  description: NullString
  containerID: NullString
}

input UpdateProduct {
  code: NullString
  cartonID: NullString
  orderID: NullString
  skuID: NullString
  contractID: NullString
  distributorID: NullString
  loyaltyPoints: NullInt
  loyaltyPointsExpire: NullTime
  inheritCartonHistory: NullBool
  description: NullString
}

input UpdateRole {
  name: NullString
  permissions: [String!]
  trackActionIDs: [String!]
}

input UpdateSKU {
  name: NullString
  code: NullString
  description: NullString
  isBeef: NullBool
  loyaltyPoints: NullInt
  masterPlanBlobID: NullString
  videoBlobID: NullString
  urls: [SKUContentInput!]
  productInfo: [SKUContentInput!]
  photoBlobIDs: [String!]
  cloneParentID: NullString
}

input UpdateTrackAction {
  name: NullString
  requirePhotos: [Boolean!]
  nameChinese: NullString
  private: NullBool
  blockchain: NullBool
}

input UpdateUser {
  email: NullString
  firstName: NullString
  lastName: NullString
  roleID: NullString
  password: NullString
  affiliateOrg: NullString
  mobilePhone: NullString
}

scalar Upload

type User {
  id: ID!
  firstName: NullString
  lastName: NullString
  email: NullString
  organisation: Organisation
  verified: Boolean!
  role: Role!
  archived: Boolean!
  mobilePhone: NullString
  mobileVerified: Boolean!
  wechatID: NullString
  loyaltyPoints: Int!
  affiliateOrg: NullString
  createdAt: Time!
}

type UserActivity {
  id: ID!
  user: User!
  action: String!
  objectID: NullString
  objectCode: NullString
  objectType: ObjectType!
  createdAt: Time!
}

type UserActivityResult {
  userActivities: [UserActivity!]!
  total: Int!
}

type UserLoyaltyActivity {
  id: ID!
  user: User!
  product: Product
  amount: Int!
  bonus: Int!
  message: String!
  transactionHash: NullString
  createdAt: Time!
}

type UsersResult {
  users: [User!]!
  total: Int!
}

