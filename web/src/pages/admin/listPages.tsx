import * as React from "react"
import { graphql } from "../../graphql"
import { FilterOption, SortDir, Perm, ObjectType } from "../../types/enums"
import { ItemListPage } from "../../components/itemList"
import { UserContainer } from "../../controllers/user"
import { ActionItemSet } from "../../types/actions"
import { UserActivity, Transaction, Settings } from "../../types/types"
import { SmallItemLink } from "../../components/smallItemLink"
import { GetItemIcon } from "../../themeOverrides"
import { HashButton } from "../../components/hashButton"
import { ManifestButton } from "../../components/manifestButton"
import { LoadingSimple } from "../../components/loading"
import { useQuery } from "@apollo/react-hooks"
import { LatestTrackActionColumn } from "../../components/common"

const Products = () => {
	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)
	if (settingsQuery.loading) {
		return <LoadingSimple />
	}
	if (settingsQuery.error) {
		return <p>Error: {settingsQuery.error.message}</p>
	}
	if (!settingsQuery.data) {
		return <p>No settings returned</p>
	}
	return (
		<ItemListPage
			itemName="product"
			query={graphql.query.PRODUCTS}
			batchActionMutation={graphql.mutation.BATCH_ACTION_PRODUCT}
			columns={[
				{
					label: "Description",
					value: "description",
					maxWidth: "200px",
				},
				{
					label: "Last Track Action",
					value: "latestTrackAction",
					filterable: true,
					itemName: "trackAction",
					resolver: row => <LatestTrackActionColumn value={row.latestTrackAction} />,
				},
				{
					label: "Date Created",
					value: "createdAt",
					dateTime: true,
				},
			]}
			extraFilterOptions={[
				{ label: "Not in Carton", id: FilterOption.ProductWithoutCarton },
				{ label: "Not in Order", id: FilterOption.ProductWithoutOrder },
				{ label: "No SKU", id: FilterOption.ProductWithoutSKU },
			]}
			itemLinks={["order", "sku", "contract", "distributor", "carton"]}
			actions={ActionItemSet.Products}
			createPermission={Perm.ProductCreate}
			readPermission={Perm.ProductRead}
			showQRCodesToggle
			qrLink={`${settingsQuery.data.settings.consumerHost}/api/steak/view?productID=`}
		/>
	)
}

const SKUs = () => (
	<ItemListPage
		itemName="sku"
		header="SKU"
		query={graphql.query.SKUS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_SKU}
		imageValue="masterPlan"
		firstColumnValue="name"
		firstColumnSubValue="code"
		columns={[
			{ label: "Products Amount", value: "productCount" },
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		actions={ActionItemSet.Archive}
		createPermission={Perm.SKUCreate}
		readPermission={Perm.SKURead}
	/>
)

const Contracts = () => (
	<ItemListPage
		itemName="contract"
		header="Livestock Specification"
		query={graphql.query.CONTRACTS}
		firstColumnValue="name"
		firstColumnSubValue="code"
		columns={[
			{ label: "Supplier", value: "supplierName" },
			{ label: "Date Signed", value: "dateSigned", dateTime: true },
		]}
		actions={ActionItemSet.Archive}
		createPermission={Perm.ContractCreate}
		readPermission={Perm.ContractRead}
	/>
)

const Orders = () => (
	<ItemListPage
		itemName="order"
		query={graphql.query.ORDERS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_ORDER}
		hash="products"
		columns={[
			{ label: "Products Amount", value: "productCount" },
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		itemLinks={["sku"]}
		itemLinksNoHash
		actions={ActionItemSet.Archive}
		createPermission={Perm.OrderCreate}
		readPermission={Perm.OrderRead}
	/>
)

const Cartons = () => (
	<ItemListPage
		itemName="carton"
		query={graphql.query.CARTONS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_CARTON}
		hash="products"
		itemLinks={["pallet"]}
		columns={[
			{ label: "Meat Type", value: "meatType" },
			{ label: "Weight", value: "weight" },
			{ label: "Product Amount", value: "productCount" },
			{
				label: "Last Track Action",
				value: "latestTrackAction",
				filterable: true,
				itemName: "trackAction",
				resolver: row => <LatestTrackActionColumn value={row.latestTrackAction} />,
			},
			{
				label: "Description",
				value: "description",
				maxWidth: "200px",
			},
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		extraFilterOptions={[{ label: "Not in Pallet", id: FilterOption.CartonWithoutPallet }]}
		actions={ActionItemSet.Cartons}
		createPermission={Perm.CartonCreate}
		readPermission={Perm.CartonRead}
		showQRCodesToggle
		pluralizeNew
	/>
)

const Pallets = () => (
	<ItemListPage
		itemName="pallet"
		query={graphql.query.PALLETS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_PALLET}
		hash="cartons"
		itemLinks={["container"]}
		columns={[
			{ label: "Cartons Amount", value: "cartonCount" },
			{
				label: "Last Track Action",
				value: "latestTrackAction",
				filterable: true,
				itemName: "trackAction",
				resolver: row => <LatestTrackActionColumn value={row.latestTrackAction} />,
			},
			{
				label: "Description",
				value: "description",
				maxWidth: "200px",
			},
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		extraFilterOptions={[{ label: "Not in Container", id: FilterOption.PalletWithoutContainer }]}
		actions={ActionItemSet.Pallets}
		createPermission={Perm.PalletCreate}
		readPermission={Perm.PalletRead}
		showQRCodesToggle
		pluralizeNew
	/>
)

const Containers = () => (
	<ItemListPage
		itemName="container"
		query={graphql.query.CONTAINERS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_CONTAINER}
		hash="pallets"
		columns={[
			{ label: "Pallets Amount", value: "palletCount" },
			{
				label: "Description",
				value: "description",
				maxWidth: "200px",
			},
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		actions={ActionItemSet.Containers}
		createPermission={Perm.ContainerCreate}
		readPermission={Perm.ContainerRead}
		showQRCodesToggle
		pluralizeNew
	/>
)

const Distributors = () => (
	<ItemListPage
		itemName="distributor"
		query={graphql.query.DISTRIBUTORS}
		firstColumnValue="name"
		firstColumnSubValue="code"
		columns={[{ label: "Date Created", value: "createdAt", dateTime: true }]}
		actions={ActionItemSet.Archive}
		createPermission={Perm.DistributorCreate}
		readPermission={Perm.DistributorRead}
	/>
)

const Users = () => {
	const { user } = UserContainer.useContainer()

	return (
		<ItemListPage
			itemName="user"
			query={graphql.query.USERS}
			identifier="email"
			firstColumnValue="lastName"
			firstColumnValueAlt="wechatID"
			firstColumnSubValue="email"
			columns={[
				{ label: "Role", value: "role", subValues: ["name"] },
				{ label: "Affiliate Organization", value: "affiliateOrg" },
				{
					label: "Date Created",
					value: "createdAt",
					dateTime: true,
				},
			]}
			disableCheck={{
				value: "role",
				subValue: "tier",
				valueCheck: (value: any) => {
					if (!user) return false
					return user.role.tier >= value
				},
			}}
			filter={{ sortDir: SortDir.Ascending }}
			createPermission={Perm.UserCreate}
			readPermission={Perm.UserRead}
		/>
	)
}

const Consumers = () => (
	<ItemListPage
		itemName="consumer"
		query={graphql.query.CONSUMERS}
		identifier="wechatID"
		firstColumnValue="wechatID"
		firstColumnSubValue="email"
		columns={[
			{ label: "Loyalty Points", value: "loyaltyPoints" },
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		filter={{ sortDir: SortDir.Ascending }}
		readPermission={Perm.UserRead}
	/>
)

const Roles = () => {
	const { user } = UserContainer.useContainer()

	return (
		<ItemListPage
			itemName="role"
			query={graphql.query.ROLES}
			identifier="name"
			columns={[
				{
					label: "Date Created",
					value: "createdAt",
					dateTime: true,
				},
			]}
			disableCheck={{
				value: "tier",
				valueCheck: (value: any) => {
					if (!user) return false
					return user.role.tier >= value
				},
			}}
			filter={{ sortDir: SortDir.Ascending }}
			createPermission={Perm.RoleCreate}
			readPermission={Perm.RoleRead}
		/>
	)
}

const TrackActions = () => (
	<ItemListPage
		itemName="trackAction"
		header="Track Action"
		query={graphql.query.TRACK_ACTIONS}
		identifier="code"
		firstColumnValue="name"
		firstColumnSubValue="nameChinese"
		columns={[
			{ label: "Private", value: "private" },
			{ label: "System", value: "system" },
			{ label: "Blockchain", value: "blockchain" },
			{ label: "Photos Required", value: "requirePhotos" },
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		extraFilterOptions={[
			{ label: "System Actions", id: FilterOption.System },
			{ label: "Blockchain", id: FilterOption.Blockchain },
		]}
		filter={{ sortDir: SortDir.Ascending }}
		createPermission={Perm.TrackActionCreate}
		readPermission={Perm.TrackActionRead}
	/>
)

const Transactions = () => (
	<ItemListPage
		pageHeader="Blockchain"
		itemName="transaction"
		query={graphql.query.ALL_TRANSACTIONS}
		columns={[
			{
				label: "Action",
				value: "action",
				subValues: ["name"],
				filterable: true,
				itemName: "trackAction",
				noWrap: true,
			},
			{
				label: "Carton",
				value: "carton",
				subValues: ["code"],
				itemLink: { itemName: "carton", IDSubValues: ["code"] },
				filterable: true,
			},
			{
				label: "Product",
				value: "product",
				subValues: ["code"],
				itemLink: { itemName: "product", IDSubValues: ["code"] },
				filterable: true,
			},
			{
				label: "User",
				value: "createdBy",
				subValues: ["email"],
				subValuesAlt: ["wechatID"],
				resolver: (row: Transaction) => {
					if (row.createdBy) {
						const itemCode = row.createdBy.email || row.createdBy.wechatID
						const itemName = !!row.createdBy.email ? "user" : "consumer"
						const columnIcon = GetItemIcon(itemName, true)
						return (
							<SmallItemLink title={itemCode.toString()} code={itemCode} link={`/portal/${itemName}`} icon={columnIcon.icon} iconLight={columnIcon.light} />
						)
					}
					return <div>{row.createdByName}</div>
				},
				filterable: true,
			},
			{
				label: "Location",
				value: "locationName",
			},
			{
				label: "Date/Time",
				value: "scannedAt",
				dateTime: true,
			},
			{
				label: "Hash",
				value: "transactionHash",
				resolver: (value: Transaction) => <HashButton hash={value.transactionHash} isLoading={value.transactionPending} />,
			},
			{
				label: "Manifest Line",
				value: "manifestID",
				resolver: (value: Transaction) => <ManifestButton manifestLineSha256={value.manifestLineSha256} isLoading={value.transactionPending} />,
			},
		]}
		extraFilterOptions={[{ label: "Pending", id: FilterOption.Pending }]}
		hideMainColumn
		hideEditButton
		disableOnClick
	/>
)

const UserActivity = () => (
	<ItemListPage
		pageHeader="User Activity"
		itemName="userActivity"
		query={graphql.query.USER_ACTIVITY}
		queryName="userActivities"
		columns={[
			{
				label: "User",
				value: "user",
				subValues: ["email"],
				subValuesAlt: ["wechatID"],
				itemLink: { itemName: "user", IDSubValues: ["email"] },
				itemLinkAlt: { itemName: "consumer", IDSubValues: ["wechatID"] },
				filterable: true,
			},
			{
				label: "Action",
				value: "action",
			},
			{
				label: "Object",
				resolver: (value: UserActivity) => {
					if (!value.objectID) return <div>{value.objectType}</div>

					const objID = value.objectID.length < 8 ? value.objectID : value.objectID.substr(0, 8) + "..."

					if (!value.objectCode || !value.objectType) return <div>{objID}</div>

					const objType = value.objectType.toLowerCase()
					const icon = GetItemIcon(objType, true)
					return (
						<SmallItemLink
							code={value.objectType == ObjectType.Blob ? "File" : value.objectCode}
							link={value.objectType == ObjectType.Blob ? `${value.objectCode}&view=true` : `/portal/${objType}`}
							fullLink={value.objectType == ObjectType.Blob}
							icon={icon.icon}
							iconLight={icon.light}
						/>
					)
				},
			},
			{
				label: "Date/Time",
				value: "createdAt",
				dateTime: true,
			},
		]}
		hideMainColumn
		hideEditButton
		disableOnClick
	/>
)

const TasksList = () => (
	<ItemListPage
		itemName="task"
		header="Task"
		query={graphql.query.TASKS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_TASK}
		firstColumnValue="title"
		firstColumnSubValue="finishDate"
		columns={[
			{ label: "Products Amount", value: "loyaltyPoints" },
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		actions={ActionItemSet.Archive}
		createPermission={Perm.TaskCreate}
		readPermission={Perm.TaskRead}
	/>
)

const ReferralList = () => (
	<ItemListPage
		itemName="referral"
		header="Referral"
		query={graphql.query.REFERRALS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_SKU}
		firstColumnValue="isRedemmed"
		firstColumnSubValue="createdAt"
		columns={[
			{ label: "UserFristName", value: "user" ,subValues: ["firstName"]},
			{ label: "UserLastName", value: "user" ,subValues: ["lastName"]},
			{ label: "RefralCode", value: "referee" ,subValues: ["referralCode"]},
			{ label: "RefereeFristName", value: "referee" ,subValues: ["firstName"]},
			{ label: "RefereeLastName", value: "referee" ,subValues: ["lastName"]},
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		actions={ActionItemSet.Archive}
		//createPermission={Perm.TaskCreate}
		//readPermission={Perm.TaskRead}
	/>
)

const UserPurchaseActivityList = () => (
	<ItemListPage
		itemName="userPurchaseActivity"
		header="UserPurchaseActivity"
		query={graphql.query.REFERRALS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_SKU}
		firstColumnValue="isRedemmed"
		firstColumnSubValue="createdAt"
		columns={[
			{ label: "UserFristName", value: "user" ,subValues: ["firstName"]},
			{ label: "UserLastName", value: "user" ,subValues: ["lastName"]},
			{ label: "RefralCode", value: "referee" ,subValues: ["referralCode"]},
			{ label: "RefereeFristName", value: "referee" ,subValues: ["firstName"]},
			{ label: "RefereeLastName", value: "referee" ,subValues: ["lastName"]},
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		actions={ActionItemSet.Archive}
	/>
)

const UserTaskList = () => (
	<ItemListPage
		itemName="userTask"
		header="UserTask"
		query={graphql.query.REFERRALS}
		batchActionMutation={graphql.mutation.BATCH_ACTION_SKU}
		firstColumnValue="isRedemmed"
		firstColumnSubValue="createdAt"
		columns={[
			{ label: "UserFristName", value: "user" ,subValues: ["firstName"]},
			{ label: "UserLastName", value: "user" ,subValues: ["lastName"]},
			{ label: "RefralCode", value: "referee" ,subValues: ["referralCode"]},
			{ label: "RefereeFristName", value: "referee" ,subValues: ["firstName"]},
			{ label: "RefereeLastName", value: "referee" ,subValues: ["lastName"]},
			{
				label: "Date Created",
				value: "createdAt",
				dateTime: true,
			},
		]}
		actions={ActionItemSet.Archive}
	/>
)
export const ListPage = {
	Products,
	SKUs,
	Contracts,
	Orders,
	Cartons,
	Pallets,
	Containers,
	Distributors,

	Users,
	Consumers,
	Roles,
	TrackActions,

	Transactions,
	UserActivity,
	TasksList,
	ReferralList,
	UserPurchaseActivityList,
	UserTaskList
}
