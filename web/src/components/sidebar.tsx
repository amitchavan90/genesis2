import * as React from "react"
import { useStyletron } from "baseui"
import { useHistory } from "react-router-dom"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { H6 } from "baseui/typography"
import { UserContainer } from "../controllers/user"
import { Perm } from "../types/enums"

export const SideBar = () => {
	const { hasPermission } = UserContainer.useContainer()

	const history = useHistory()

	const onInventoryPage =
		history.location.pathname.startsWith("/portal/container") ||
		history.location.pathname.startsWith("/portal/pallet") ||
		history.location.pathname.startsWith("/portal/carton") ||
		history.location.pathname.startsWith("/portal/product")
	const onActivityPage = history.location.pathname.startsWith("/portal/activity/blockchain") || history.location.pathname.startsWith("/portal/activity/user")
	const onSystemPage =
		history.location.pathname.startsWith("/portal/user") ||
		history.location.pathname.startsWith("/portal/consumer") ||
		history.location.pathname.startsWith("/portal/role") ||
		history.location.pathname.startsWith("/portal/trackAction")
	const onAppPage = 
		history.location.pathname.startsWith("/portal/tasks") || 
		history.location.pathname.startsWith("/portal/referrals")||
		history.location.pathname.startsWith("/portal/purchesActivity")
	const [expandInventoryPanel, setExpandInventoryPanel] = React.useState(onInventoryPage)
	const [expandActivityPanel, setExpandActivityPanel] = React.useState(onActivityPage)
	const [expandSystemPanel, setExpandSystemPanel] = React.useState(onSystemPage)
	const [expandAppPanel, setExpandAppPanel] = React.useState(onAppPage)

	const [css, theme] = useStyletron()

	const containerStyle: string = css({
		display: "flex",
		flexDirection: "column",
		minHeight: "100vh",
		maxHeight: "100vh",
		width: "280px",
		flexShrink: 0,
		backgroundColor: theme.colors.primary,
		boxShadow: "0px 3px 6px #00000029",
		overflow: "auto",
	})

	const logoContainerStyle: string = css({
		width: "100%",
		height: "150px",
		display: "flex",
		justifyContent: "center",
		alignItems: "center",

		paddingTop: "20px",
		paddingBottom: "20px",

		"@media screen and (max-height: 1080px)": {
			height: "80px",
		},
	})

	const accordionPanelInventory: string = css({
		transformOrigin: "top",
		transform: expandInventoryPanel ? "scaleY(1)" : "scaleY(0)",
		maxHeight: expandInventoryPanel ? "300px" : "0px",
		transition: "0.25s",
	})
	const accordionPanelActivity: string = css({
		transformOrigin: "top",
		transform: expandActivityPanel ? "scaleY(1)" : "scaleY(0)",
		maxHeight: expandActivityPanel ? "150px" : "0px",
		transition: "0.25s",
	})
	const accordionPanelSystem: string = css({
		transformOrigin: "top",
		transform: expandSystemPanel ? "scaleY(1)" : "scaleY(0)",
		maxHeight: expandSystemPanel ? "150px" : "0px",
		transition: "0.25s",
	})
	const accordionPanelApp: string = css({
		transformOrigin: "top",
		transform: expandAppPanel ? "scaleY(1)" : "scaleY(0)",
		maxHeight: expandAppPanel ? "300px" : "0px",
		transition: "0.25s",
	})

	const showInventoryPanel = hasPermission(Perm.ContainerList) || hasPermission(Perm.CartonList) || hasPermission(Perm.ProductList)
	const showActivityPanel = hasPermission(Perm.ActivityListBlockchainActivity) || hasPermission(Perm.ActivityListUserActivity)
	const showSystemPanel = hasPermission(Perm.UserList) || hasPermission(Perm.RoleList) || hasPermission(Perm.TrackActionList)
	const showAppPanel = hasPermission(Perm.TaskList) || hasPermission(Perm.ReferralList) || hasPermission(Perm.UserPurchaseActivityList)
	React.useEffect(() => {
		if (expandSystemPanel != onSystemPage) setExpandSystemPanel(onSystemPage)
		if (expandActivityPanel != onActivityPage) setExpandActivityPanel(onActivityPage)
		if (expandInventoryPanel != onInventoryPage) setExpandInventoryPanel(onInventoryPage)
	}, [history.location.pathname])

	var index = 0

	return (
		<div className={containerStyle}>
			<div className={logoContainerStyle}>
				<FontAwesomeIcon icon={["fas", "steak"]} size="8x" color="white" style={{ height: "100%" }} />
			</div>
			<div>
				<SideMenuButton index={index++} icon="chart-pie" label="Dashboard" url="/portal" strictSelection />

				{hasPermission(Perm.OrderList) && <SideMenuButton index={index++} icon="shopping-cart" label="Orders" url="/portal/orders" altURL="/portal/order" />}

				{showInventoryPanel && (
					<>
						<SideMenuButton
							index={index++}
							icon="inventory"
							label="Inventory"
							selected={onInventoryPage}
							onClick={() => {
								if (onInventoryPage) return
								setExpandInventoryPanel(!expandInventoryPanel)
								if (!onInventoryPage) history.push("/portal/containers")
							}}
						/>
						<div className={accordionPanelInventory}>
							{hasPermission(Perm.ContainerList) && (
								<SideMenuButton
									index={index++}
									subMenu
									icon="container-storage"
									iconLight
									label="Containers"
									url="/portal/containers"
									altURL="/portal/container"
								/>
							)}
							{hasPermission(Perm.PalletList) && (
								<SideMenuButton index={index++} subMenu icon="pallet-alt" iconLight label="Pallets" url="/portal/pallets" altURL="/portal/pallet" />
							)}
							{hasPermission(Perm.CartonList) && (
								<SideMenuButton index={index++} subMenu icon="box" iconLight label="Cartons" url="/portal/cartons" altURL="/portal/carton" />
							)}
							{hasPermission(Perm.ProductList) && (
								<SideMenuButton index={index++} subMenu icon="steak" iconLight label="Products" url="/portal/products" altURL="/portal/product" />
							)}
						</div>
					</>
				)}

				{hasPermission(Perm.SKUList) && <SideMenuButton index={index++} icon="barcode-alt" iconLight label="SKUs" url="/portal/skus" altURL="/portal/sku" />}

				{hasPermission(Perm.DistributorList) && (
					<SideMenuButton index={index++} icon="shopping-basket" label="Distributors" url="/portal/distributors" altURL="/portal/distributor" />
				)}

				{hasPermission(Perm.ContractList) && (
					<SideMenuButton
						index={index++}
						icon="file-contract"
						label="ORIGIN"
						url="/portal/contracts"
						altURL="/portal/contract"
						fontSize="16px"
					/>
				)}

				{showActivityPanel && (
					<>
						<SideMenuButton
							index={index++}
							icon="chart-line"
							label="Activity"
							selected={onActivityPage}
							onClick={() => {
								if (onActivityPage) return
								setExpandActivityPanel(!expandActivityPanel)
								if (!onActivityPage) history.push("/portal/activity/blockchain")
							}}
						/>
						<div className={accordionPanelActivity}>
							{hasPermission(Perm.ActivityListBlockchainActivity) && (
								<SideMenuButton index={index++} subMenu icon="cubes" label="Blockchain" url="/portal/activity/blockchain" />
							)}
							{hasPermission(Perm.ActivityListUserActivity) && (
								<SideMenuButton index={index++} subMenu icon="user-chart" label="User Activity" url="/portal/activity/user" />
							)}
						</div>
					</>
				)}
				{showAppPanel && (
					<>
						<SideMenuButton
							index={index++}
							icon="chart-line"
							label="App"
							selected={onAppPage}
							onClick={() => {
								if (onAppPage) return
								setExpandAppPanel(!expandAppPanel)
								if (!onAppPage) history.push("/portal/tasks")
							}}
						/>
						<div className={accordionPanelApp}>
							{hasPermission(Perm.TaskList) && (
								<SideMenuButton index={index++} subMenu icon="file-contract" label="Tasks" url="/portal/tasks" />
							)}
							{hasPermission(Perm.ReferralList) && (
								<SideMenuButton index={index++} subMenu icon="user-chart" label="Referrals" url="/portal/referrals" />
							)}
							{hasPermission(Perm.UserPurchaseActivityList) && (
								<SideMenuButton index={index++} subMenu icon="steak" label="Purchase Activity" url="/portal/purchesActivity" />
							)}
							{hasPermission(Perm.UserTaskList) && ( 
							<SideMenuButton index={index++} subMenu icon="users" label="User Tasks" url="/portal/userTasks"/>
							)}
						</div>
					</>
				)}
				{showSystemPanel && (
					<>
						<SideMenuButton
							index={index++}
							icon="cogs"
							label="System"
							selected={onSystemPage}
							onClick={() => {
								if (onSystemPage) return
								setExpandSystemPanel(!expandSystemPanel)
								if (!onSystemPage) history.push("/portal/users")
							}}
						/>
						<div className={accordionPanelSystem}>
							{hasPermission(Perm.UserList) && <SideMenuButton index={index++} subMenu icon="users" label="Users" url="/portal/users" altURL="/portal/user" />}
							{hasPermission(Perm.UserList) && (
								<SideMenuButton index={index++} subMenu icon="smile" label="Consumers" url="/portal/consumers" altURL="/portal/consumer" />
							)}
							{hasPermission(Perm.RoleList) && (
								<SideMenuButton index={index++} subMenu icon="user-tag" label="Roles" url="/portal/roles" altURL="/portal/role" />
							)}
							{hasPermission(Perm.TrackActionList) && (
								<SideMenuButton index={index++} subMenu icon="truck-moving" label="Track Actions" url="/portal/trackActions" altURL="/portal/trackAction" />
							)}
						</div>
					</>
				)}
			</div>
		</div>
	)
}

interface ButtonProps {
	index: number
	label: string
	icon: IconName
	iconLight?: boolean
	/** URL to go to + highlight button if current pathname contains this url  */
	url?: string
	/** Highlight button if current pathname contains this url */
	altURL?: string
	/** Highlight button ONLY if the current pathname equals the url */
	strictSelection?: boolean
	/** Highlight the button override */
	selected?: boolean
	onClick?: () => void
	/** Part of a submenu (will indent) */
	subMenu?: boolean
	fontSize?: string
}

const SideMenuButton = (props: ButtonProps) => {
	const [css, theme] = useStyletron()
	const { index, label, icon, url, altURL, strictSelection, subMenu, fontSize } = props

	const history = useHistory()

	const selected =
		props.selected ||
		(url && !strictSelection && history.location.pathname.startsWith(url)) ||
		(url && strictSelection && history.location.pathname == url) ||
		(altURL && history.location.pathname.startsWith(altURL))

	const buttonStyle: string = css({
		height: "75px",
		width: subMenu ? "calc(100% - 3rem)" : "calc(100% - 1.5rem)",
		backgroundColor: selected ? "rgba(0, 0, 0, 0.75)" : index % 2 == 0 ? "rgba(0, 0, 0, 0.2)" : "transparent",
		display: "flex",
		paddingLeft: subMenu ? "3rem" : "1.5rem",
		alignItems: "center",
		color: "white",
		textAlign: "center",
		":hover": {
			backgroundColor: selected ? "rgba(0, 0, 0, 0.75)" : "rgba(0, 0, 0, 0.5)",
		},
		cursor: "pointer",
	})

	const buttonInnerStyle: string = css({
		display: "flex",
	})

	const iconStyle: string = css({
		width: "40px",
		color: selected ? "white" : "grey",
	})

	return (
		<div className={buttonStyle} onClick={url ? () => history.push(url) : props.onClick}>
			<div className={buttonInnerStyle}>
				<div className={iconStyle}>
					<FontAwesomeIcon icon={[props.iconLight ? "fal" : "fas", icon]} size="2x" />
				</div>
				<H6
					color={selected ? "white" : "grey"}
					overrides={{
						Block: {
							style: {
								marginTop: 0,
								marginBottom: 0,
								marginLeft: "0.75rem",
								lineHeight: "32px",
								...(fontSize ? { fontSize: fontSize } : {}),
							},
						},
					}}
				>
					{label}
				</H6>
			</div>
		</div>
	)
}
