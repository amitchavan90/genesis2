import * as React from "react"
import { useStyletron } from "baseui"
import { TopBar } from "../../components/topBar"
import { SideBar } from "../../components/sidebar"
import { Switch, Route, Redirect } from "react-router"
import { UserContainer } from "../../controllers/user"
import { AuthContainer } from "../../controllers/auth"
import { VerificationComplete } from "../verificationComplete"
import { NotLoggedIn } from "../notLogedIn"
import { Loading } from "../../components/loading"
import { Perm } from "../../types/enums"

import { ListPage } from "./listPages"
import { UserEdit, ConsumerEdit } from "./system/userEdit"
const SKUEdit = React.lazy(() => import("./skuEdit"))
const TaskEdit = React.lazy(() => import("./inventory/taskEdit"))
const RoleEdit = React.lazy(() => import("./system/roleEdit"))
const ContainerEdit = React.lazy(() => import("./inventory/containerEdit"))
const PalletEdit = React.lazy(() => import("./inventory/palletEdit"))
const CartonEdit = React.lazy(() => import("./inventory/cartonEdit"))
const ProductEdit = React.lazy(() => import("./inventory/productEdit"))
const OrderEdit = React.lazy(() => import("./inventory/orderEdit"))
const TrackActionEdit = React.lazy(() => import("./system/trackActionEdit"))
const ContractEdit = React.lazy(() => import("./contractEdit"))
const DistributorEdit = React.lazy(() => import("./distributorEdit"))
const Dashboard = React.lazy(() => import("./dashboard"))
const Settings = React.lazy(() => import("./settings"))

const PortalInner = () => {
	const [css, theme] = useStyletron()

	const containerStyle: string = css({
		...theme.typography.font300,
		color: theme.colors.colorPrimary,
		display: "flex",
		minHeight: "100vh",
		width: "100%",
	})

	const mainStyle: string = css({
		display: "flex",
		width: "100%",
		flexDirection: "column",
	})

	const contentStyle: string = css({
		width: "100%",
		height: "calc(100vh - 71px)",
		overflowX: "auto",
	})

	const { user, loading, useFetchUser, hasPermission } = UserContainer.useContainer()
	const { logoutRedirect, setLogoutRedirect, showVerifyComplete } = AuthContainer.useContainer()
	useFetchUser()

	// Redirect triggered by logout
	if (logoutRedirect) {
		setLogoutRedirect(false)
		return <Redirect to={"/"} />
	}

	if (!loading && !user) {
		return <NotLoggedIn />
	}

	if (!user) {
		return <Loading />
	}

	return (
		<React.Fragment>
			<div className={containerStyle}>
				<SideBar />
				<div className={mainStyle}>
					<TopBar />
					<div className={contentStyle}>
						<Switch>
							<Route path={"/portal/settings"} component={Settings} />
							{hasPermission(Perm.ContractList) && <Route path={"/portal/contracts"} component={ListPage.Contracts} />}
							{hasPermission(Perm.ContractRead) && <Route path={"/portal/contract/:code"} component={ContractEdit} />}

							{hasPermission(Perm.OrderList) && <Route path={"/portal/orders"} component={ListPage.Orders} />}
							{hasPermission(Perm.OrderRead) && <Route path={"/portal/order/:code"} component={OrderEdit} />}

							{hasPermission(Perm.ContainerList) && <Route path={"/portal/containers"} component={ListPage.Containers} />}
							{hasPermission(Perm.ContainerRead) && <Route path={"/portal/container/:code"} component={ContainerEdit} />}

							{hasPermission(Perm.PalletList) && <Route path={"/portal/pallets"} component={ListPage.Pallets} />}
							{hasPermission(Perm.PalletRead) && <Route path={"/portal/pallet/:code"} component={PalletEdit} />}

							{hasPermission(Perm.CartonList) && <Route path={"/portal/cartons"} component={ListPage.Cartons} />}
							{hasPermission(Perm.CartonRead) && <Route path={"/portal/carton/:code"} component={CartonEdit} />}

							{hasPermission(Perm.ProductList) && <Route path={"/portal/products"} component={ListPage.Products} />}
							{hasPermission(Perm.ProductRead) && <Route path={"/portal/product/:code"} component={ProductEdit} />}

							{hasPermission(Perm.ActivityListBlockchainActivity) && <Route path={"/portal/activity/blockchain"} component={ListPage.Transactions} />}

							{hasPermission(Perm.ActivityListUserActivity) && <Route path={"/portal/activity/user"} component={ListPage.UserActivity} />}

							{hasPermission(Perm.SKUUpdate) && <Route path={"/portal/sku/:code"} component={SKUEdit} />}
							{hasPermission(Perm.SKUList) && <Route path={"/portal/skus"} component={ListPage.SKUs} />}

							{hasPermission(Perm.DistributorUpdate) && <Route path={"/portal/distributor/:code"} component={DistributorEdit} />}
							{hasPermission(Perm.DistributorList) && <Route path={"/portal/distributors"} component={ListPage.Distributors} />}

							{hasPermission(Perm.UserUpdate) && <Route path={"/portal/user/:email"} component={UserEdit} />}
							{hasPermission(Perm.UserList) && <Route path={"/portal/users"} component={ListPage.Users} />}

							{hasPermission(Perm.UserList) && <Route path={"/portal/consumers"} component={ListPage.Consumers} />}
							{hasPermission(Perm.UserUpdate) && <Route path={"/portal/consumer/:wechatID"} component={ConsumerEdit} />}

							{hasPermission(Perm.RoleUpdate) && <Route path={"/portal/role/:name"} component={RoleEdit} />}
							{hasPermission(Perm.RoleList) && <Route path={"/portal/roles"} component={ListPage.Roles} />}

							{hasPermission(Perm.TrackActionUpdate) && <Route path={"/portal/trackAction/:code"} component={TrackActionEdit} />}
							{hasPermission(Perm.TrackActionList) && <Route path={"/portal/trackActions"} component={ListPage.TrackActions} />}

							{hasPermission(Perm.TaskCreate) && <Route path={"/portal/task/:id"} component={TaskEdit} />}
							{hasPermission(Perm.TaskList) && <Route path={"/portal/tasks"} component={ListPage.TasksList} />}

							{hasPermission(Perm.ReferralList) && <Route path={"/portal/referrals"} component={ListPage.ReferralList} />}
							{hasPermission(Perm.UserPurchaseActivityList) && <Route path={"/portal/purchesActivity"} component={ListPage.ReferralList} />}
							{hasPermission(Perm.UserTaskList) && <Route path={"/portal/userTask"} component={ListPage.ReferralList} />}
							<Route path={"/portal/"} component={Dashboard} />
						</Switch>
					</div>
				</div>
			</div>
			{showVerifyComplete && <VerificationComplete />}
		</React.Fragment>
	)
}

const Portal = () => {
	return (
		<UserContainer.Provider>
			<PortalInner />
		</UserContainer.Provider>
	)
}

export default Portal
