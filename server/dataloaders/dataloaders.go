//go:generate go run github.com/vektah/dataloaden UserLoader string *genesis/db.User
//go:generate go run github.com/vektah/dataloaden ReferralLoader string *genesis/db.Referral
//go:generate go run github.com/vektah/dataloaden TaskLoader string *genesis/db.Task
//go:generate go run github.com/vektah/dataloaden SubtaskLoader string *genesis/db.Subtask
//go:generate go run github.com/vektah/dataloaden UserTaskLoader string *genesis/db.UserTask
//go:generate go run github.com/vektah/dataloaden UserSubtaskLoader string *genesis/db.UserSubtask
//go:generate go run github.com/vektah/dataloaden RoleLoader string *genesis/db.Role
//go:generate go run github.com/vektah/dataloaden OrganisationLoader string *genesis/db.Organisation
//go:generate go run github.com/vektah/dataloaden OrganisationUsersLoader string []*genesis/db.User
//go:generate go run github.com/vektah/dataloaden SKULoader string *genesis/db.StockKeepingUnit
//go:generate go run github.com/vektah/dataloaden CategoryLoader string *genesis/db.Category
//go:generate go run github.com/vektah/dataloaden ProductCategoryLoader string *genesis/db.ProductCategory
//go:generate go run github.com/vektah/dataloaden RetailLinkLoader string *genesis/db.RetailLink
//go:generate go run github.com/vektah/dataloaden OrderLoader string *genesis/db.Order
//go:generate go run github.com/vektah/dataloaden DistributorLoader string *genesis/db.Distributor
//go:generate go run github.com/vektah/dataloaden ContainerLoader string *genesis/db.Container
//go:generate go run github.com/vektah/dataloaden PalletLoader string *genesis/db.Pallet
//go:generate go run github.com/vektah/dataloaden CartonLoader string *genesis/db.Carton
//go:generate go run github.com/vektah/dataloaden ProductLoader string *genesis/db.Product
//go:generate go run github.com/vektah/dataloaden UserPurchaseActivityLoader string *genesis/db.UserPurchaseActivity
//go:generate go run github.com/vektah/dataloaden ContractLoader string *genesis/db.Contract
//go:generate go run github.com/vektah/dataloaden TransactionLoader string *genesis/db.Transaction
//go:generate go run github.com/vektah/dataloaden ManifestLoader string *genesis/db.Manifest
//go:generate go run github.com/vektah/dataloaden TrackActionLoader string *genesis/db.TrackAction

package dataloaders
