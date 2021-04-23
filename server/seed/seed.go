package seed

import (
	"fmt"
	"genesis/crypto"
	"genesis/db"
	"genesis/graphql"
	"genesis/helpers"
	"genesis/store"
	"math/rand"
	"time"

	"github.com/volatiletech/null"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/types"
)

var roleSuperAdmin = "SUPERADMIN"
var roleOrganisationAdmin = "ORGANISATIONADMIN"
var roleMember = "MEMBER"
var roleConsumer = "CONSUMER"

// Run seed for database spinup
func Run(conn *sqlx.DB) error {
	var err error

	fmt.Println("Seeding settings")
	err = Settings(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding track actions")
	err = TrackActions(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding roles")
	err = Roles(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding organisations")
	err = Organisations(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding users")
	err = Users(conn, false)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding inventory")
	err = Inventory(conn, false)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seed complete")
	return nil
}

// RunProd seed for production database spinup
func RunProd(conn *sqlx.DB) error {
	var err error

	fmt.Println("Seeding settings")
	err = Settings(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding track actions")
	err = TrackActions(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding roles")
	err = Roles(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding organisations")
	err = Organisations(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding users")
	err = Users(conn, true)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seeding inventory")
	err = Inventory(conn, true)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Seed complete")
	return nil
}

// Settings seed for database spinup
func Settings(conn *sqlx.DB) error {
	exists, err := db.Settings().Exists(conn)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Settings already initialized")
	}

	setting := &db.Setting{}

	return setting.Insert(conn, boil.Infer())
}

// Organisations seed for database spinup
func Organisations(conn *sqlx.DB) error {
	ninjaSoftware := &db.Organisation{Name: "Ninja Software"}
	err := ninjaSoftware.Insert(conn, boil.Infer())
	if err != nil {
		return err
	}

	l28 := &db.Organisation{Name: "Latitude28 Produce"}
	err = l28.Insert(conn, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

// TrackActions seed for database spinup
func TrackActions(conn *sqlx.DB) error {
	trackActionsDefault := []*db.TrackAction{
		// System actions (actions logged by system)
		{
			Name:       "Contract Created",
			Private:    true,
			System:     true,
			Blockchain: true,
		},
		{
			Name:       "Registered",
			Private:    false,
			System:     true,
			Blockchain: true,
		},
		{
			Name:        "Moved to Carton",
			NameChinese: "搬到纸箱",
			Private:     true,
			System:      true,
		},
		{
			Name:        "Moved to Pallet",
			NameChinese: "移至棧板",
			Private:     true,
			System:      true,
		},
		{
			Name:        "Moved to Container",
			NameChinese: "移至容器",
			Private:     true,
			System:      true,
		},
		{
			Name:        "Removed from Carton",
			NameChinese: "从纸箱中取出",
			Private:     true,
			System:      true,
		},
		{
			Name:        "Removed from Pallet",
			NameChinese: "从棧板中取出",
			Private:     true,
			System:      true,
		},
		{
			Name:        "Removed from Container",
			NameChinese: "从容器中取出",
			Private:     true,
			System:      true,
		},

		// Public actions (shown under tracking information on product view page)
		{Name: "Picked up from Farm", Blockchain: true},
		{Name: "Received for Processing", Blockchain: true},
		{Name: "Processed and Packaged", Blockchain: true, RequirePhotos: []bool{true, true}},
		{Name: "Dispatched", Blockchain: true},
		{Name: "Entered Cold Storage", NameChinese: "貨物已入冷庫", Blockchain: true},
		{Name: "Hand Selected", Blockchain: true},
	}

	for i, a := range trackActionsDefault {
		a.Code = fmt.Sprintf("TRACK%03d", i)
		err := a.Insert(conn, boil.Infer())
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// Roles seed for database spinup
func Roles(conn *sqlx.DB) error {
	// Get track actions
	trackActions, err := db.TrackActions(db.TrackActionWhere.System.EQ(false)).All(conn)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Super Admin
	allPerms := types.StringArray{}
	for _, perm := range graphql.AllPerm {
		allPerms = append(allPerms, string(perm))
	}
	r := &db.Role{
		Name:        roleSuperAdmin,
		Permissions: allPerms,
		Tier:        1,
	}
	err = r.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = r.AddTrackActions(conn, false, trackActions...)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Organisation Admin
	r2 := &db.Role{
		Name: roleOrganisationAdmin,
		Permissions: types.StringArray{
			string(graphql.PermUserList),
			string(graphql.PermUserCreate),
			string(graphql.PermUserRead),
			string(graphql.PermUserUpdate),
			string(graphql.PermUserArchive),
			string(graphql.PermUserUnarchive),

			string(graphql.PermRoleList),
			string(graphql.PermRoleCreate),
			string(graphql.PermRoleRead),
			string(graphql.PermRoleUpdate),
			string(graphql.PermRoleArchive),
			string(graphql.PermRoleUnarchive),

			string(graphql.PermOrganisationRead),

			string(graphql.PermReferralList),
			string(graphql.PermReferralRead),

			string(graphql.PermTaskList),
			string(graphql.PermTaskCreate),
			string(graphql.PermTaskRead),
			string(graphql.PermTaskUpdate),
			string(graphql.PermTaskArchive),
			string(graphql.PermTaskUnarchive),

			string(graphql.PermUserTaskList),
			string(graphql.PermUserTaskCreate),
			string(graphql.PermUserTaskRead),
			string(graphql.PermUserTaskUpdate),
			string(graphql.PermUserTaskArchive),
			string(graphql.PermUserTaskUnarchive),

			string(graphql.PermCategoryList),
			string(graphql.PermCategoryCreate),
			string(graphql.PermCategoryRead),
			string(graphql.PermCategoryUpdate),
			string(graphql.PermCategoryArchive),
			string(graphql.PermCategoryUnarchive),

			string(graphql.PermProductCategoryList),
			string(graphql.PermProductCategoryCreate),
			string(graphql.PermProductCategoryRead),
			string(graphql.PermProductCategoryUpdate),
			string(graphql.PermProductCategoryArchive),
			string(graphql.PermProductCategoryUnarchive),

			string(graphql.PermSKUList),
			string(graphql.PermSKURead),
			string(graphql.PermSKUCreate),
			string(graphql.PermSKUUpdate),
			string(graphql.PermSKUArchive),
			string(graphql.PermSKUUnarchive),

			string(graphql.PermContainerList),
			string(graphql.PermContainerRead),
			string(graphql.PermContainerCreate),
			string(graphql.PermContainerUpdate),
			string(graphql.PermContainerArchive),
			string(graphql.PermContainerUnarchive),

			string(graphql.PermPalletList),
			string(graphql.PermPalletRead),
			string(graphql.PermPalletCreate),
			string(graphql.PermPalletUpdate),
			string(graphql.PermPalletArchive),
			string(graphql.PermPalletUnarchive),

			string(graphql.PermCartonList),
			string(graphql.PermCartonRead),
			string(graphql.PermCartonCreate),
			string(graphql.PermCartonUpdate),
			string(graphql.PermCartonArchive),
			string(graphql.PermCartonUnarchive),

			string(graphql.PermProductList),
			string(graphql.PermProductRead),
			string(graphql.PermProductCreate),
			string(graphql.PermProductUpdate),
			string(graphql.PermProductArchive),
			string(graphql.PermProductUnarchive),

			string(graphql.PermOrderList),
			string(graphql.PermOrderRead),
			string(graphql.PermOrderCreate),
			string(graphql.PermOrderUpdate),
			string(graphql.PermOrderArchive),
			string(graphql.PermOrderUnarchive),

			string(graphql.PermTrackActionList),
			string(graphql.PermTrackActionRead),
			string(graphql.PermTrackActionCreate),
			string(graphql.PermTrackActionUpdate),
			string(graphql.PermTrackActionArchive),
			string(graphql.PermTrackActionUnarchive),

			string(graphql.PermContractList),
			string(graphql.PermContractRead),
			string(graphql.PermContractCreate),
			string(graphql.PermContractUpdate),
			string(graphql.PermContractArchive),
			string(graphql.PermContractUnarchive),

			string(graphql.PermDistributorList),
			string(graphql.PermDistributorRead),
			string(graphql.PermDistributorCreate),
			string(graphql.PermDistributorUpdate),
			string(graphql.PermDistributorArchive),
			string(graphql.PermDistributorUnarchive),

			string(graphql.PermUserPurchaseActivityList),
			string(graphql.PermUserPurchaseActivityRead),
			string(graphql.PermUserPurchaseActivityCreate),
			string(graphql.PermUserPurchaseActivityUpdate),

			string(graphql.PermActivityListBlockchainActivity),
			string(graphql.PermActivityListUserActivity),
			string(graphql.PermUseAdvancedMode),
			string(graphql.PermUseAdminPortal),
		},
		Tier: 2,
	}
	err = r2.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = r2.AddTrackActions(conn, false, trackActions...)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Member
	r3 := &db.Role{
		Name: roleMember,
		Permissions: types.StringArray{
			string(graphql.PermReferralList),
			string(graphql.PermReferralRead),

			string(graphql.PermTaskList),
			string(graphql.PermTaskCreate),
			string(graphql.PermTaskRead),
			string(graphql.PermTaskUpdate),
			string(graphql.PermTaskArchive),
			string(graphql.PermTaskUnarchive),

			string(graphql.PermUserTaskList),
			string(graphql.PermUserTaskCreate),
			string(graphql.PermUserTaskRead),
			string(graphql.PermUserTaskUpdate),
			string(graphql.PermUserTaskArchive),
			string(graphql.PermUserTaskUnarchive),

			string(graphql.PermCategoryList),
			string(graphql.PermCategoryCreate),
			string(graphql.PermCategoryRead),
			string(graphql.PermCategoryUpdate),
			string(graphql.PermCategoryArchive),
			string(graphql.PermCategoryUnarchive),

			string(graphql.PermProductCategoryList),
			string(graphql.PermProductCategoryCreate),
			string(graphql.PermProductCategoryRead),
			string(graphql.PermProductCategoryUpdate),
			string(graphql.PermProductCategoryArchive),
			string(graphql.PermProductCategoryUnarchive),

			string(graphql.PermSKUList),
			string(graphql.PermSKURead),
			string(graphql.PermSKUCreate),
			string(graphql.PermSKUUpdate),
			string(graphql.PermSKUArchive),
			string(graphql.PermSKUUnarchive),

			string(graphql.PermContainerList),
			string(graphql.PermContainerRead),
			string(graphql.PermContainerCreate),
			string(graphql.PermContainerUpdate),
			string(graphql.PermContainerArchive),
			string(graphql.PermContainerUnarchive),

			string(graphql.PermPalletList),
			string(graphql.PermPalletRead),
			string(graphql.PermPalletCreate),
			string(graphql.PermPalletUpdate),
			string(graphql.PermPalletArchive),
			string(graphql.PermPalletUnarchive),

			string(graphql.PermCartonList),
			string(graphql.PermCartonRead),
			string(graphql.PermCartonCreate),
			string(graphql.PermCartonUpdate),
			string(graphql.PermCartonArchive),
			string(graphql.PermCartonUnarchive),

			string(graphql.PermProductList),
			string(graphql.PermProductRead),
			string(graphql.PermProductCreate),
			string(graphql.PermProductUpdate),
			string(graphql.PermProductArchive),
			string(graphql.PermProductUnarchive),

			string(graphql.PermUserPurchaseActivityList),
			string(graphql.PermUserPurchaseActivityRead),
			string(graphql.PermUserPurchaseActivityCreate),
			string(graphql.PermUserPurchaseActivityUpdate),
		},
	}
	err = r3.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = r3.AddTrackActions(conn, false, trackActions...)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Consumer
	roleConsumer := &db.Role{
		Name: roleConsumer,
		Permissions: types.StringArray{
			string(graphql.PermReferralList),
			string(graphql.PermReferralRead),

			string(graphql.PermUserRead),
			string(graphql.PermUserUpdate),

			string(graphql.PermTaskList),
			string(graphql.PermTaskRead),

			string(graphql.PermUserTaskList),
			string(graphql.PermUserTaskCreate),
			string(graphql.PermUserTaskRead),
			string(graphql.PermUserTaskUpdate),
			string(graphql.PermUserTaskArchive),

			string(graphql.PermCategoryList),
			string(graphql.PermCategoryRead),
			string(graphql.PermCategoryUpdate),

			string(graphql.PermProductCategoryList),
			string(graphql.PermProductCategoryRead),

			string(graphql.PermSKUList),
			string(graphql.PermSKURead),

			string(graphql.PermProductList),
			string(graphql.PermProductRead),

			string(graphql.PermUserPurchaseActivityList),
			string(graphql.PermUserPurchaseActivityRead),
			string(graphql.PermUserPurchaseActivityCreate),
			string(graphql.PermUserPurchaseActivityUpdate),
		},
		Tier: 10,
	}
	err = roleConsumer.Insert(conn, boil.Infer())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Users seed for database spinup
func Users(conn *sqlx.DB, isProduction bool) error {
	orgs, err := db.Organisations().All(conn)
	if err != nil {
		return err
	}

	roleStore := store.NewRoleStore(conn)
	roleSuperAdmin, err := roleStore.GetByName(roleSuperAdmin)
	if err != nil {
		return err
	}
	roleOrgAdmin, err := roleStore.GetByName(roleOrganisationAdmin)
	if err != nil {
		return err
	}
	roleMember, err := roleStore.GetByName(roleMember)
	if err != nil {
		return err
	}

	for i, org := range orgs {
		if i == 0 {
			// Ninja Software
			u := store.UserFactory()
			u.Email = null.StringFrom("superadmin@example.com")
			u.FirstName = null.StringFrom("Super")
			u.LastName = null.StringFrom("Admin")
			u.PasswordHash = crypto.HashPassword("******")
			u.OrganisationID = null.StringFrom(org.ID)
			u.RoleID = roleSuperAdmin.ID
			u.Verified = true
			err = u.Insert(conn, boil.Infer())
			if err != nil {
				return err
			}

			u2 := store.UserFactory()
			u2.Email = null.StringFrom("john+orgadmin@ninjasoftware.com.au")
			u.FirstName = null.StringFrom("John")
			u.LastName = null.StringFrom("Nguyen")
			u2.PasswordHash = crypto.HashPassword("******")
			u2.OrganisationID = null.StringFrom(org.ID)
			u2.RoleID = roleOrgAdmin.ID
			u2.Verified = true
			err = u2.Insert(conn, boil.Infer())
			if err != nil {
				return err
			}

			u3 := store.UserFactory()
			u3.Email = null.StringFrom("john@ninjasoftware.com.au")
			u.FirstName = null.StringFrom("John")
			u.LastName = null.StringFrom("Nguyen")
			u3.PasswordHash = crypto.HashPassword("******")
			u3.OrganisationID = null.StringFrom(org.ID)
			u3.RoleID = roleMember.ID
			u3.Verified = true
			err = u3.Insert(conn, boil.Infer())
			if err != nil {
				return err
			}
		} else if i == 1 {
			// Latitude28 Produce
			u := store.UserFactory()
			u.FirstName = null.StringFrom("James")
			u.LastName = null.StringFrom("Williamson")
			u.Email = null.StringFrom("james@latitude28produce.com")
			u.PasswordHash = crypto.HashPassword("password")
			if isProduction {
				u.PasswordHash = crypto.HashPassword("4f45c5f9-56f1-4194-95b7-7bc915688f7d")
			}
			u.OrganisationID = null.StringFrom(org.ID)
			u.RoleID = roleSuperAdmin.ID
			u.Verified = true
			err = u.Insert(conn, boil.Infer())
			if err != nil {
				return err
			}

			u2 := store.UserFactory()
			u.FirstName = null.StringFrom("Rhys")
			u.LastName = null.StringFrom("Williamson")
			u2.Email = null.StringFrom("rhys@latitude28produce.com")
			u2.PasswordHash = crypto.HashPassword("password")
			if isProduction {
				u.PasswordHash = crypto.HashPassword("9a486a51-e209-4c18-b21f-b8a9f7c927e7")
			}
			u2.OrganisationID = null.StringFrom(org.ID)
			u2.RoleID = roleSuperAdmin.ID
			u2.Verified = true
			err = u2.Insert(conn, boil.Infer())
			if err != nil {
				return err
			}
		}
	}

	// stop if production
	if isProduction {
		return nil
	}

	for i := 0; i < 5; i++ {
		u := store.UserFactory()
		u.RoleID = roleMember.ID
		u.PasswordHash = crypto.HashPassword("password")
		u.Verified = true
		err := u.Insert(conn, boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

// Inventory seed for database spinup
func Inventory(conn *sqlx.DB, isProduction bool) error {
	user, err := db.Users().One(conn)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	palletCount := 0
	cartonCount := 0

	// Contracts
	contracts := 3
	if isProduction {
		contracts = 1
	}
	fmt.Printf("  - contracts (%v)\n", contracts)
	for i := 0; i < contracts; i++ {
		contract := store.ContractFactory()
		contract.Name = "Placeholder contract"
		contract.CreatedByID = user.ID
		contract.Code = fmt.Sprintf("CONTRACT%05d", i)
		err = contract.Insert(conn, boil.Infer())
		if err != nil {
			return err
		}
	}

	// SKU
	sku := store.SKUFactory()
	sku.Code = fmt.Sprintf("L28%05d", 1)
	sku.CreatedByID = user.ID
	sku.LoyaltyPoints = 10
	// Seed last SKU
	sku.Name = "28°South Latitude Shahao Steak"
	sku.Description =
		`Portholes steak is taken from the outside back of the beef, which is characterized by its meatiness.
		Tenderness, with a layer of Pupu fat on the top to increase the mouth sense. The flat arched pan is fired under the fire to lock in the delicious gravy.
		All L28 products are complete in the Australian packaging committee,
		It is certified and supervised by Australian export workers until
		Before you open the package with your home, the product will remain in place Han Zhou Benju's original appearance.`
	sku.IsAppBound = true
	sku.IsPointBound = false
	sku.Price = 100
	sku.Currency = "AUD"
	sku.Weight = 1000
	sku.WeightUnit = "gram"
	sku.PurchasePoints = 0

	err = sku.Insert(conn, boil.Infer())
	if err != nil {
		return err
	}

	skuContent := []*db.StockKeepingUnitContent{
		{Title: "Produce", Content: "Australia", ContentType: db.ContentTypeINFO},
		{Title: "Processing", Content: "Australia", ContentType: db.ContentTypeINFO},
		{Title: "Package", Content: "Australia", ContentType: db.ContentTypeINFO},
		{Title: "Pure natural Australian beef", Content: "Quality assurance", ContentType: db.ContentTypeINFO},
		{Title: "Pasture farming", Content: "Australia Environmental Ecological Farming", ContentType: db.ContentTypeINFO},
		{Title: "Halal", Content: "Verified", ContentType: db.ContentTypeINFO},

		{Title: "28° Layin on South Street", Content: "https://www.latitude28produce.com", ContentType: db.ContentTypeURL},
		{Title: "28° South Latitude Official Website", Content: "https://www.latitude28produce.com", ContentType: db.ContentTypeURL},
		{Title: "Australian government certification", Content: "https://www.latitude28produce.com", ContentType: db.ContentTypeURL},
	}

	for _, content := range skuContent {

		err := sku.AddSkuStockKeepingUnitContents(conn, true, content)
		if err != nil {
			return err
		}
	}

	// Containers, Pallets and Cartons
	containers := 4
	if isProduction {
		containers = 1
	}
	fmt.Printf("  - containers (%v)\n", containers)
	for i := 0; i < containers; i++ {
		// Containers
		container := store.ContainerFactory()
		container.Code = fmt.Sprintf("CON%05d", i)
		container.CreatedByID = user.ID
		err = container.Insert(conn, boil.Infer())
		if err != nil {
			return err
		}

		// Pallets
		pallets := 1 + rand.Intn(3)
		if isProduction {
			pallets = 1
		}
		fmt.Printf("    - pallets (%v)\n", pallets)
		for c := 0; c < pallets; c++ {
			pallet := store.PalletFactory()
			pallet.Code = fmt.Sprintf("PAL%05d", palletCount)
			pallet.CreatedByID = user.ID
			pallet.ContainerID = null.StringFrom(container.ID)
			err = pallet.Insert(conn, boil.Infer())
			if err != nil {
				return err
			}

			palletCount++

			// Cartons
			cartons := 5 + rand.Intn(10)
			if isProduction {
				cartons = 5
			}
			ssLink := helpers.GetCartonSpreadsheetLink(fmt.Sprintf("CAR%05d", cartonCount), cartons, cartonCount)
			fmt.Printf("      - cartons (%v)\n", cartons)
			for c := 0; c < cartons; c++ {

				carton := store.CartonFactory()
				carton.Code = fmt.Sprintf("CAR%05d", cartonCount)
				carton.CreatedByID = user.ID
				carton.PalletID = null.StringFrom(pallet.ID)
				carton.SpreadsheetLink = ssLink

				err = carton.Insert(conn, boil.Infer())
				if err != nil {
					return err
				}
				cartonCount++
			}
		}
	}

	return nil
}
