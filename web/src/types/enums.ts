export enum Enviroment {
	Admin = "Admin",
	Consumer = "Consumer",
}

export enum Perm {
	UserList = "UserList",
	UserCreate = "UserCreate",
	UserRead = "UserRead",
	UserUpdate = "UserUpdate",
	UserArchive = "UserArchive",
	UserUnarchive = "UserUnarchive",
	OrganisationList = "OrganisationList",
	OrganisationCreate = "OrganisationCreate",
	OrganisationRead = "OrganisationRead",
	OrganisationUpdate = "OrganisationUpdate",
	OrganisationArchive = "OrganisationArchive",
	OrganisationUnarchive = "OrganisationUnarchive",

	RoleList = "RoleList",
	RoleCreate = "RoleCreate",
	RoleRead = "RoleRead",
	RoleUpdate = "RoleUpdate",
	RoleArchive = "RoleArchive",
	RoleUnarchive = "RoleUnarchive",

	SKUList = "SKUList",
	SKUCreate = "SKUCreate",
	SKURead = "SKURead",
	SKUUpdate = "SKUUpdate",
	SKUArchive = "SKUArchive",
	SKUUnarchive = "SKUUnarchive",

	ContainerList = "ContainerList",
	ContainerCreate = "ContainerCreate",
	ContainerRead = "ContainerRead",
	ContainerUpdate = "ContainerUpdate",
	ContainerArchive = "ContainerArchive",
	ContainerUnarchive = "ContainerUnarchive",

	PalletList = "PalletList",
	PalletCreate = "PalletCreate",
	PalletRead = "PalletRead",
	PalletUpdate = "PalletUpdate",
	PalletArchive = "PalletArchive",
	PalletUnarchive = "PalletUnarchive",

	CartonList = "CartonList",
	CartonCreate = "CartonCreate",
	CartonRead = "CartonRead",
	CartonUpdate = "CartonUpdate",
	CartonArchive = "CartonArchive",
	CartonUnarchive = "CartonUnarchive",

	ProductList = "ProductList",
	ProductCreate = "ProductCreate",
	ProductRead = "ProductRead",
	ProductUpdate = "ProductUpdate",
	ProductArchive = "ProductArchive",
	ProductUnarchive = "ProductUnarchive",

	OrderList = "OrderList",
	OrderCreate = "OrderCreate",
	OrderRead = "OrderRead",
	OrderUpdate = "OrderUpdate",
	OrderArchive = "OrderArchive",
	OrderUnarchive = "OrderUnarchive",

	TrackActionList = "TrackActionList",
	TrackActionRead = "TrackActionRead",
	TrackActionCreate = "TrackActionCreate",
	TrackActionUpdate = "TrackActionUpdate",
	TrackActionArchive = "TrackActionArchive",
	TrackActionUnarchive = "TrackActionUnarchive",

	TaskList = "TaskList",
	TaskRead = "TaskRead",
	TaskCreate = "TaskCreate",

	ReferralList = "ReferralList",

	ContractList = "ContractList",
	ContractRead = "ContractRead",
	ContractCreate = "ContractCreate",
	ContractUpdate = "ContractUpdate",
	ContractArchive = "ContractArchive",
	ContractUnarchive = "ContractUnarchive",

	DistributorList = "DistributorList",
	DistributorRead = "DistributorRead",
	DistributorCreate = "DistributorCreate",
	DistributorUpdate = "DistributorUpdate",
	DistributorArchive = "DistributorArchive",
	DistributorUnarchive = "DistributorUnarchive",

	ActivityListBlockchainActivity = "ActivityListBlockchainActivity",
	ActivityListUserActivity = "ActivityListUserActivity",

	UseAdvancedMode = "UseAdvancedMode",
	UseAdminPortal = "UseAdminPortal",
}

export const AffiliateOrgs = [
	{ id: "Latitude 28° Produce" },
	{ id: "NCMC - Australian Certified Processor" },
	{ id: "Virgin Airways" },
	{ id: "PIL Logistics - Mainfreight" },
]

export const enum ObjectType {
	Self = "Self",
	User = "User",
	Role = "Role",
	Sku = "Sku",
	Contract = "Contract",
	Distributor = "Distributor",
	Order = "Order",
	Container = "Container",
	Pallet = "Pallet",
	Carton = "Carton",
	Product = "Product",
	TrackAction = "TrackAction",
	Blob = "Blob",
}

export enum FilterOption {
	All = "All",
	Active = "Active",
	Archived = "Archived",

	ProductWithoutOrder = "ProductWithoutOrder",
	ProductWithoutCarton = "ProductWithoutCarton",
	ProductWithoutSKU = "ProductWithoutSKU",
	CartonWithoutPallet = "CartonWithoutPallet",
	PalletWithoutContainer = "PalletWithoutContainer",

	System = "System",
	Blockchain = "Blockchain",

	Pending = "Pending",
}
export enum SortByOption {
	DateCreated = "DateCreated",
	DateUpdated = "DateUpdated",
	Alphabetical = "Alphabetical",
}
export enum SortDir {
	Ascending = "Ascending",
	Descending = "Descending",
}

export interface FilterOptionItem {
	label: string
	id?: FilterOption
}
export interface SortByOptionItem {
	label: string
	id: SortByOption
}
