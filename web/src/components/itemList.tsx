import * as React from "react"
import { useStyletron } from "baseui"
import { useQuery } from "@apollo/react-hooks"
import { SearchFilter, Settings } from "../types/types"
import { Block } from "baseui/block"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useHistory } from "react-router-dom"
import { PaginationBar } from "../components/paginationBar"
import { FilterOptionItem, SortByOptionItem, Perm } from "../types/enums"
import { Checkbox } from "baseui/checkbox"
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic"
import { paddingZero, GetItemIcon } from "../themeOverrides"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { SmallItemLink } from "./smallItemLink"
import { Spread } from "./spread"
import { Spaced } from "./spaced"
import { H1 } from "baseui/typography"
import { Button } from "baseui/button"
import { ActionsButton } from "./actionsButton"
import { ActionItem } from "../types/actions"
import { UserContainer } from "../controllers/user"
import { QRCode } from "react-qrcode-logo"
import { timeAgo } from "../helpers/time"
import { Value } from "baseui/select"
import { graphql } from "../graphql"
import { Modal, ModalHeader, ModalBody, ModalFooter } from "baseui/modal"
import { ItemSelectList } from "./itemSelectList"
import { LoadingSimple } from "./loading"
import { Spinner } from "baseui/spinner"

export interface ItemListColumn {
	label: string
	value?: string
	subValues?: string[]
	/** Mirrors subValues (if subValues[0] is undefined, use subValuesAlt[0]) */
	subValuesAlt?: string[]
	dateTime?: Boolean
	itemLink?: ItemLink
	/** Mirrors itemLink (used if itemLink fails to resolve) */
	itemLinkAlt?: ItemLink
	resolver?: (value: any) => JSX.Element
	maxWidth?: string
	noWrap?: boolean
	filterable?: boolean
	/** Used for filtering */
	itemName?: string
}

interface ItemLink {
	itemName: string
	IDSubValues: string[]
}

interface ItemListProps {
	query: any
	batchActionMutation?: any
	itemName: string
	/** Name of list returned from query (defaults to itemName + s) */
	queryName?: string
	link?: string
	/** Default tab to open to (eg: products) */
	hash?: string
	/** Main column header (defaults to itemName+s) */
	header?: string
	/** Page Header (defaults to itemName) */
	pageHeader?: string
	/** Added 's' at the end of the create button (eg: "New Carton" will show as "New Cartons") */
	pluralizeNew?: boolean
	/** The main value to id the item (default: "code") */
	identifier?: string
	/** The value to use for the first column (defaults to identifier) */
	firstColumnValue?: string
	/** Alternative value to use for first column if firstColumnValue is null */
	firstColumnValueAlt?: string
	/** Display smaller text under the first column's value */
	firstColumnSubValue?: string
	/** Disable rows if [value] fails the valueCheck  */
	disableCheck?: { value: string; subValue?: string; valueCheck: (value: any) => boolean }
	imageValue?: string

	filter?: SearchFilter
	limit?: number
	containerID?: string
	palletID?: string
	cartonID?: string
	productID?: string
	orderID?: string
	distributorID?: string
	skuID?: string
	contractID?: string
	userID?: string
	trackActionID?: string

	/** Additional columns to show in the table
	 *
	 * Use subValues for nested values (eg: `{ label: "SKU", value: "sku", subValues: ["name"] }` for sku.name) */
	columns?: ItemListColumn[]
	extraFilterOptions?: FilterOptionItem[]
	extraSortByOptions?: SortByOptionItem[]
	/** Names of the parent values to show as links under each item (eg: order, carton, sku) */
	itemLinks?: string[]
	itemLinksNoHash?: boolean
	/** Bulk Actions that can be done on the list items */
	actions?: ActionItem[]
	/** Called when the list is updated due to an action */
	onListUpdate?: () => void
	showQRCodesToggle?: boolean
	qrLink?: string

	hideEditButton?: boolean
	hideMainColumn?: boolean
	disableOnClick?: boolean

	createPermission?: Perm
	readPermission?: Perm
}

export const ItemListPage = (props: ItemListProps) => {
	const { pageHeader, header, itemName, createPermission } = props

	const { hasPermission } = UserContainer.useContainer()
	const history = useHistory()

	if (itemName.length <= 0) return <div />

	const icon = GetItemIcon(itemName)
	const iconAlt = GetItemIcon(itemName, true)
	const listName = `${itemName[0].toUpperCase()}${itemName.substr(1)}`

	return (
		<div style={{ margin: "20px 100px" }}>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={[icon.light ? "fal" : "fas", icon.icon]} size="3x" />
					<H1>{pageHeader || `${header || listName}s`}</H1>
				</Spaced>
				{createPermission && hasPermission(createPermission) ? (
					<Button
						onClick={() => history.push(`/portal/${itemName}/new`)}
						startEnhancer={<FontAwesomeIcon icon={[iconAlt.light ? "fal" : "fas", iconAlt.icon]} size="lg" />}
					>
						{`New ${header || listName}${props.pluralizeNew ? "s" : ""}`}
					</Button>
				) : (
					<div />
				)}
			</Spread>

			<ItemList {...props} />
		</div>
	)
}

export const ItemList = (props: ItemListProps) => {
	const {
		filter,
		limit,
		query,
		batchActionMutation,
		itemName,
		hash,
		columns,
		disableCheck,
		firstColumnValue,
		firstColumnValueAlt,
		firstColumnSubValue,
		imageValue,
		actions,
		onListUpdate,
		showQRCodesToggle,
		qrLink,
		readPermission,
		hideEditButton,
		disableOnClick,
	} = props
	const queryName = props.queryName || `${itemName}s`
	const link = props.link || `/portal/${itemName}`
	const icon = GetItemIcon(itemName)
	const identifier = props.identifier || "code"

	const { hasPermission } = UserContainer.useContainer()
	const history = useHistory()

	const [showQRCodes, setShowQRCodes] = React.useState(false)

	// Filter values
	const [containerID, setContainerID] = React.useState(props.containerID)
	const [palletID, setPalletID] = React.useState(props.palletID)
	const [cartonID, setCartonID] = React.useState(props.cartonID)
	const [productID, setProductID] = React.useState(props.productID)
	const [orderID, setOrderID] = React.useState(props.orderID)
	const [distributorID, setDistributorID] = React.useState(props.distributorID)
	const [skuID, setSkuID] = React.useState(props.skuID)
	const [contractID, setContractID] = React.useState(props.contractID)
	const [userID, setUserID] = React.useState(props.userID)
	const [trackActionID, setTrackActionID] = React.useState(props.trackActionID)

	const [showModal, setShowModal] = React.useState<boolean>()
	const [modalColumn, setModalColumn] = React.useState<number>()
	const [selectValue, setSelectValue] = React.useState<Value>()

	const getListQuery = (column: string) => {
		switch (column) {
			case "container":
				return graphql.query.CONTAINERS
			case "pallet":
				return graphql.query.PALLETS
			case "carton":
				return graphql.query.CARTONS
			case "product":
				return graphql.query.PRODUCTS
			case "order":
				return graphql.query.ORDERS
			case "sku":
				return graphql.query.SKUS
			case "contract":
				return graphql.query.CONTRACTS
			case "user":
				return graphql.query.USERS
			case "trackAction":
				return graphql.query.TRACK_ACTIONS
		}
		return graphql.query.CONTRACTS
	}
	const getListQueryIndentifier = (column: string) => {
		if (column === "contract" || column === "trackAction") return "name"
		if (column === "user") return "email"
		return undefined
	}
	const setColumnFilter = (column: string, id?: string) => {
		// clear other filters
		setContainerID(undefined)
		setPalletID(undefined)
		setCartonID(undefined)
		setProductID(undefined)
		setOrderID(undefined)
		setDistributorID(undefined)
		setSkuID(undefined)
		setContractID(undefined)
		setUserID(undefined)
		setTrackActionID(undefined)

		switch (column) {
			case "container":
				setContainerID(id)
				break
			case "pallet":
				setPalletID(id)
				break
			case "carton":
				setCartonID(id)
				break
			case "product":
				setProductID(id)
				break
			case "order":
				setOrderID(id)
				break
			case "distributor":
				setDistributorID(id)
				break
			case "sku":
				setSkuID(id)
				break
			case "contract":
				setContractID(id)
				break
			case "user":
				setUserID(id)
				break
			case "trackAction":
				setTrackActionID(id)
				break
		}
	}

	// Pagination
	const [search, setSearch] = React.useState<SearchFilter>(filter || {})
	const [offset, setOffset] = React.useState(0)

	const [items, setItems] = React.useState<any[]>([])
	const [total, setTotal] = React.useState(0)

	const pagination = (
		<PaginationBar
			offset={offset}
			total={total}
			limit={limit}
			defaultFilter={filter && filter.filter}
			defaultSortOption={filter && filter.sortBy}
			defaultSortDir={filter && filter.sortDir}
			setSearch={(value: SearchFilter) => setSearch(value)}
			setOffset={(value: number) => setOffset(value)}
			extraFilterOptions={props.extraFilterOptions}
			extraSortByOptions={props.extraSortByOptions}
		/>
	)

	// Query
	const { data, loading, error, refetch } = useQuery(query, {
		fetchPolicy: "network-only",
		variables: {
			search,
			limit: limit || 20,
			offset,

			containerID,
			palletID,
			cartonID,
			productID,
			orderID,
			distributorID,
			skuID,
			contractID,
			userID,
			trackActionID,
		},
	})

	// Checkboxes
	const hasAll = items.length > 0 && items.every(x => x.selected)
	const hasSome = items.length > 0 && items.some(x => x.selected)

	React.useEffect(() => {
		if (!data || !data[queryName]) return
		setItems(data[queryName][queryName])
		setTotal(data[queryName].total)
	}, [data, loading, error])

	// refetch every 60 secs if no activity
	React.useEffect(() => {
		const interval = setInterval(() => {
			if (!hasSome) {
				refetch()
			}
		}, 60000)
		return () => clearInterval(interval)
	}, [, hasSome])

	const toggleAll = () => {
		setItems(items =>
			items.map(row => ({
				...row,
				selected: !hasAll,
			})),
		)
	}
	const toggle = (event: any) => {
		const { name, checked } = event.currentTarget
		setItems(items =>
			items.map(row => ({
				...row,
				selected: String(row[identifier]) === name ? checked : row.selected,
			})),
		)
	}

	const itemDisabled = (i: any) =>
		disableCheck && disableCheck.valueCheck(disableCheck.subValue ? i[disableCheck.value][disableCheck.subValue] : i[disableCheck.value])

	// Styling
	const [css, theme] = useStyletron()
	const editButton = css({
		height: showQRCodes ? "222px" : "72px",
		width: "72px",
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		cursor: "pointer",
		color: "inherit",
		":hover": {
			backgroundColor: "grey",
		},
	})
	const dateStyle = css({
		display: "flex",
		flexDirection: "column",
	})
	const fullDateStyle = css({
		fontSize: "0.8rem",
		lineHeight: "0.8rem",
		color: "grey",
	})
	const columnFilterStyle = css({
		color: "grey",
		float: "right",
	})
	const flexStyle = css({
		display: "flex",
	})
	const iconStyle = css({
		marginRight: "10px",
		width: "20px",
		height: "20px",
		textAlign: "center",
	})
	const spinnerStyle = css({
		width: "100%",
		textAlign: "center",
		padding: "20px 0",
	})

	const getItemLink = (item: any, value: any, columnValue: any, itemLink: ItemLink) => {
		const columnIcon = GetItemIcon(itemLink.itemName, true)
		let itemCode = item[columnValue]
		if (itemLink.IDSubValues != undefined) itemLink.IDSubValues.forEach(s => (itemCode = itemCode[s]))

		if (itemCode === null || itemCode === undefined) return undefined

		return (
			<SmallItemLink
				key={`smallItemLink-${itemLink.itemName}-${itemCode}}`}
				title={value.toString()}
				code={itemCode}
				link={`/portal/${itemLink?.itemName}`}
				icon={columnIcon.icon}
				iconLight={columnIcon.light}
			/>
		)
	}

	return (
		<>
			<Spread>
				{actions && batchActionMutation && (
					<ActionsButton
						options={actions}
						items={items}
						batchActionMutation={batchActionMutation}
						refetch={() => {
							refetch()
							if (onListUpdate) onListUpdate()
						}}
						disabled={!hasSome}
					/>
				)}
				{showQRCodesToggle ? (
					<Button kind="secondary" onClick={() => setShowQRCodes(!showQRCodes)}>
						{showQRCodes ? "Hide QR Codes" : "Show QR Codes"}
					</Button>
				) : (
					<div />
				)}
			</Spread>

			{pagination}

			<TableBuilder
				data={items}
				onSort={(columnID: string) => {
					if (!columns) return
					// columnID = `column-${index}`
					const index = +columnID.replace("column-", "")
					setModalColumn(index)

					// set column itemName to label (lowercased) if not set
					if (!columns[index].itemName) {
						columns[index].itemName = columns[index].label.toLowerCase()
					}

					// Open filter modal
					setShowModal(true)
				}}
				overrides={{
					TableBodyRow: {
						style: {
							":hover": {
								backgroundColor: theme.colors.colorPrimary,
								color: "white",
							},
						},
					},
					TableHeadCellSortable: {
						style: {
							paddingRight: "10px",
						},
					},
					SortNoneIcon: {
						component: () => <FontAwesomeIcon icon={["fas", "filter"]} size="1x" className={columnFilterStyle} />,
					},
				}}
			>
				{/* Checkbox */}
				{actions && (
					<TableBuilderColumn
						overrides={{
							TableHeadCell: { style: { width: "1%" } },
							TableBodyCell: {
								style: {
									width: "1%",
									height: "auto",
									paddingTop: 0,
									paddingBottom: 0,
									verticalAlign: "center",
								},
							},
						}}
						header={<Checkbox checked={hasAll} isIndeterminate={!hasAll && hasSome} onChange={toggleAll} />}
					>
						{item => <Checkbox name={item[identifier]} checked={item.selected} onChange={toggle} />}
					</TableBuilderColumn>
				)}

				{/* Main Column */}
				{!props.hideMainColumn && (
					<TableBuilderColumn
						header={props.header || `${itemName[0].toUpperCase()}${itemName.substr(1)}`}
						overrides={{
							TableBodyCell: {
								style: {
									...paddingZero,
									verticalAlign: "center",
									color: "inherit",
								},
							},
						}}
					>
						{item => (
							<ItemCard
								item={item}
								itemName={itemName}
								itemLinks={props.itemLinks}
								itemLinkHash={props.itemLinksNoHash ? "" : queryName}
								icon={icon.icon}
								iconLight={icon.light}
								columnValue={firstColumnValue || identifier}
								columnValueAlt={firstColumnValueAlt}
								columnSubValue={firstColumnSubValue}
								imageValue={imageValue}
								disabled={disableOnClick || (readPermission && !hasPermission(readPermission)) || itemDisabled(item)}
								onClick={() => history.push(`${link}/${item[identifier]}${hash ? `#${hash}` : ""}`)}
								showQRCodes={showQRCodes}
								qrLink={qrLink}
							/>
						)}
					</TableBuilderColumn>
				)}

				{/* Columns */}
				{columns &&
					columns.map((c, index) => {
						return (
							<TableBuilderColumn
								key={`column-${c.label}`}
								id={`column-${index}`}
								header={c.label}
								sortable={c.filterable}
								overrides={{
									TableBodyCell: {
										style: {
											verticalAlign: "center",
											color: "inherit",
										},
									},
									TableHeadCell: {
										style: {
											whiteSpace: "normal",
										},
									},
								}}
							>
								{item => {
									// Custom resolver
									if (c.resolver !== undefined) return c.resolver(item)
									if (!c.value || !item[c.value]) return <></>

									// Get cell value
									let value = item[c.value]
									if (c.subValues !== undefined) c.subValues.forEach(s => (value = value[s]))
									if (!value && c.subValuesAlt !== undefined) {
										value = item[c.value]
										c.subValuesAlt.forEach(s => (value = value[s]))
									}

									// Date/Time?
									if (c.dateTime) {
										return (
											<div className={dateStyle}>
												<div>{timeAgo(value)}</div>
												<div className={fullDateStyle}>{`(${new Date(value).toLocaleString()})`}</div>
											</div>
										)
									}

									// Boolean (or null)?
									if (value === true) return <FontAwesomeIcon icon={["fas", "check"]} />
									if (value === false || value === null) return <></>

									// Boolean array?
									if (Array.isArray(value) && value.length > 0) {
										if (!value.some(v => v == true)) return <></>
										return (
											<div>
												{value.map((v, idx) => {
													if (v === true) return <FontAwesomeIcon key={idx} icon={["fas", "check"]} style={{ marginRight: "10px" }} />
													if (v === false) return <FontAwesomeIcon key={idx} icon={["fas", "times"]} style={{ marginRight: "10px" }} />
													return <></>
												})}
											</div>
										)
									}

									// Item Link
									if (c.itemLink) {
										const il = getItemLink(item, value, c.value, c.itemLink)
										if (!il && c.itemLinkAlt !== undefined) return getItemLink(item, value, c.value, c.itemLinkAlt)
										return il
									}

									// String
									return <div style={{ maxWidth: c.maxWidth || "unset", whiteSpace: c.noWrap ? "nowrap" : "unset" }}>{value.toString()}</div>
								}}
							</TableBuilderColumn>
						)
					})}

				{/* Edit Button */}
				{!hideEditButton && (!readPermission || hasPermission(readPermission)) && (
					<TableBuilderColumn
						overrides={{
							TableHeadCell: { style: { width: "1%" } },
							TableBodyCell: {
								style: {
									width: "1%",
									...paddingZero,
									color: "inherit",
								},
							},
						}}
					>
						{item =>
							itemDisabled(item) ? (
								<div />
							) : (
								<div className={editButton} onClick={() => history.push(`${link}/${item[identifier]}`)}>
									<FontAwesomeIcon icon={["fas", "edit"]} />
								</div>
							)
						}
					</TableBuilderColumn>
				)}
			</TableBuilder>

			{loading && (
				<div className={spinnerStyle}>
					<Spinner />
				</div>
			)}

			{items.length > 0 && pagination}

			{/* Filter by column modal */}
			<Modal isOpen={showModal} onClose={() => setShowModal(false)}>
				{columns && modalColumn !== undefined && (
					<>
						<ModalHeader>
							<div className={flexStyle}>
								<div className={iconStyle}>
									<FontAwesomeIcon icon={["fas", "filter"]} />
								</div>
								<div>{`Filter by ${columns[modalColumn].label}`}</div>
							</div>
						</ModalHeader>
						<ModalBody>
							<ItemSelectList
								itemName={columns[modalColumn].itemName || ""}
								identifier={getListQueryIndentifier(columns[modalColumn].itemName || "")}
								value={selectValue}
								setValue={setSelectValue}
								query={getListQuery(columns[modalColumn].itemName || "")}
							/>
						</ModalBody>
						<ModalFooter>
							<Button
								onClick={() => {
									const id = selectValue && selectValue.length > 0 ? (selectValue[0].id as string) : undefined
									setColumnFilter(columns[modalColumn].itemName || "", id)

									setModalColumn(undefined)
									setShowModal(false)
								}}
							>
								Filter
							</Button>
						</ModalFooter>
					</>
				)}
			</Modal>
		</>
	)
}

interface ItemCardProps {
	item: any
	itemName: string
	icon: IconName
	iconLight?: boolean
	disabled?: boolean
	onClick?: () => void
	itemLinks?: string[]
	itemLinkHash: string
	columnValue: string
	columnValueAlt?: string
	columnSubValue?: string
	imageValue?: string
	showQRCodes?: boolean
	qrLink?: string
}

export const ItemCard = (props: ItemCardProps) => {
	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)

	const { item, icon, iconLight, itemLinks, itemLinkHash, disabled, columnValue, columnValueAlt, columnSubValue, imageValue, showQRCodes, qrLink } = props

	const [css, theme] = useStyletron()
	const containerStyle = css({
		display: "flex",
		opacity: item.archived ? 0.5 : "unset",
		alignItems: "center",
		minHeight: showQRCodes ? "222px" : "72px",
		cursor: disabled ? "cursor" : "pointer",
	})
	const iconItem = css({
		display: "flex",
		width: "60px",
		textAlign: "center",
		justifyContent: "center",
	})
	const qrItem = css({
		display: "flex",
		width: "200px",
		textAlign: "center",
		justifyContent: "center",
	})
	const nameStyle = css({
		fontSize: "1.15rem",
		lineHeight: "1rem",
	})
	const nameSubStyle = css({
		fontSize: "0.8rem",
		lineHeight: "1rem",
		color: "grey",
		marginTop: "3px",
	})
	const descriptionStyle = css({
		fontSize: "0.8rem",
		lineHeight: "1rem",
		color: "grey",
		marginTop: "3px",
		maxWidth: "200px",
		whiteSpace: "nowrap",
		overflow: "hidden",
		textOverflow: "ellipsis",
	})
	const linksContainer = css({
		display: "flex",
		opacity: item.archived ? 0.5 : "unset",
		flexFlow: "wrap",
	})
	const imageStyle = css({
		height: "64px",
		maxWidth: "64px",
		objectFit: "contain",
		textAlign: "center",
		margin: "auto",
	})
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
		<div className={containerStyle} onClick={disabled ? undefined : props.onClick}>
			{showQRCodes && (
				<div className={qrItem}>
					<QRCode value={`${qrLink || `${settingsQuery.data.settings.adminHost}/q/`}${item.id}`} size={200} quietZone={2} />
				</div>
			)}
			<div className={iconItem}>
				{imageValue !== undefined && item[imageValue]?.file_url ? (
					<img className={imageStyle} src={item[imageValue].file_url} />
				) : (
					<FontAwesomeIcon icon={[iconLight ? "fal" : "fas", icon]} size="2x" />
				)}
			</div>
			<Block margin="auto 0 auto 8px" padding="10px 0">
				{columnSubValue ? (
					<Block margin="auto 0 auto 8px">
						<div className={nameStyle}>{item[columnValue] || (columnValueAlt && item[columnValueAlt])}</div>
						<div className={nameSubStyle}>{item[columnSubValue]}</div>
					</Block>
				) : (
					<div className={nameStyle}>{item[columnValue]}</div>
				)}

				<div className={descriptionStyle}>{item.archived ? "[ARCHIVED]" : ""}</div>

				{itemLinks && (
					<div className={linksContainer}>
						{itemLinks.map((itemName, index) => {
							const linkIcon = GetItemIcon(itemName, true)
							if (!item[itemName]) return <div key={`smallItemLink-${index}}`} />
							return (
								<SmallItemLink
									key={`smallItemLink-${index}}`}
									title={item[itemName]["name"]}
									code={item[itemName]["code"]}
									link={`/portal/${itemName}`}
									hash={itemLinkHash}
									icon={linkIcon.icon}
									iconLight={linkIcon.light}
								/>
							)
						})}
					</div>
				)}
			</Block>
		</div>
	)
}
