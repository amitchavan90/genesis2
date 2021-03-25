import { Perm, FilterOption, SortByOption, SortDir, ObjectType } from "./enums"

export interface UserError {
	message: string
	field: string[]
}

export interface Blob {
	id: string
	file_url: string
}

export interface Onboard {
	prospect: Prospect
}

export interface Prospect {
	id: string
	email: string
	firstName?: string
	lastName?: string
	onboardingComplete: boolean
}

export interface OnboardingInput {
	email: string
	firstName?: string
	lastName?: string
}

export interface Login {
	email: string
	password: string
}

export interface Role {
	id: string
	name: string
	tier: number
	archived: boolean
	permissions: Perm[]
	trackActions: { id: string }[]
}

export interface User {
	id: string
	email: string
	firstName: string
	lastName: string
	verified: boolean
	role: Role
	affiliateOrg: string
	wechatID: string
	mobilePhone: string
	mobileVerified: boolean
	archived: boolean

	loyaltyPoints: number
}

export interface SKUContent {
	title: string
	content: string
}

export interface SKUCategory {
	name: string
}

export interface SKU {
	id: string
	name: string
	brand :string
	price: number
	purchasePoints: number
	weight: number
	ingredients: string
	categories: SKUCategory[]
	productCategories: SKUCategory[]
	code: string
	description: string
	isBeef: boolean
	isAppSku: boolean
	isPointSku: boolean
	loyaltyPoints: number
	archived: boolean

	cloneParentID?: string
	hasClones: boolean

	masterPlan: Blob
	video: Blob
	photos: Blob[]
	urls: SKUContent[]
	productInfo: SKUContent[]
}
export interface SKUClone {
	sku: SKU
	depth: number
}

export interface SubTask {
	title: string
	description: string
}

export interface Task{
	id:string
    title:string
	code: string
    description:String
    loyaltyPoints:number
    isTimeBound:boolean
    isPeopleBound:boolean
    isProductRelevant: boolean
    finishDate:Date
    maximumPeople:number
	sku: SKU
	subtasks:SubTask[]
}

export interface Container {
	id: string
	code: string
	archived: boolean
	palletCount: number
	description: string
}

export interface Pallet {
	id: string
	code: string
	archived: boolean
	cartonCount: number
	description: string

	container: Container
	latestTrackAction?: LatestTransactionInfo
}

export interface Carton {
	id: string
	code: string
	weight: string
	processedAt: string
	archived: boolean
	productCount: number
	description: string
	meatType: string
	spreadsheetLink: string

	pallet: Pallet
	order: Order
	sku: SKU
	distributor: Distributor

	transactions: Transaction[]
	latestTrackAction?: LatestTransactionInfo
	createdAt: Date
}

export interface CartonOption {
	carton: Carton
	selected: boolean
}

export interface Product {
	id: string
	code: string
	description: string
	registerID: string
	loyaltyPoints: number
	loyaltyPointsExpire: string
	archived: boolean
	createdAt: string

	sku: SKU
	carton: Carton
	order: Order
	contract: Contract
	distributor: Distributor

	transactions: Transaction[]
	latestTrackAction?: LatestTransactionInfo

	registered: boolean
	registeredBy?: User
}

export interface Distributor {
	id: string
	code: string
	name: string

	archived: boolean
	createdAt: string
}

export interface Settings {
	consumerHost: string
	adminHost: string
	etherscanHost: string
	fieldappVersion: string
	smartContractAddress: string
}

export interface TrackAction {
	id: string
	name: string
	nameChinese: string
	private: boolean
	blockchain: boolean
	system: boolean
	requirePhotos: boolean[]
	archived: boolean
}

export interface Transaction {
	id: string
	transactionHash?: string
	manifestID?: string
	manifest?: Manifest
	manifestLineJson?: string
	manifestLineSha256?: string
	transactionPending?: boolean
	locationGeohash?: string
	locationName?: string
	memo?: string
	action: TrackAction
	createdAt: string
	createdBy: User
	createdByName: string

	carton?: Carton
	product?: Product

	photos: {
		cartonPhoto: Blob | null
		productPhoto: Blob | null
	}

	scannedAt?: string
}

export interface Manifest {
	id: string
	transactionHash?: string
	merkleRootSha256?: string
}

export interface Contract {
	id: string
	code: string
	name: string
	description: string
	latitude:number
	longitude: number
	supplierName: string
	dateSigned: string
	archived: boolean
	createdAt: string
}

export interface Order {
	id: string
	code: string
	archived: boolean
	productCount: number

	sku: SKU
}

export interface ErrorMap {
	[key: string]: string
}

export interface SearchFilter {
	search?: string
	filter?: FilterOption
	sortBy?: SortByOption
	sortDir?: SortDir
}

export interface UserActivity {
	id: string
	user: User
	action: string
	objectID: string
	objectCode: string
	objectType: ObjectType
	createdAt: string
}

export interface UserLoyaltyActivity {
	id: string
	amount: number
	bonus: number
	message: string
	transactionHash?: string
	createdAt: string
	product?: Product
}

export interface TickerInfo {
	lastTick: string
	tickInterval: number
}

export interface LatestTransactionInfo {
	name: string
	createdAt: string
}
