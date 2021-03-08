import * as React from "react"
import MetaTags from "react-meta-tags"
import { AuthContainer } from "../../controllers/auth"
import { BrowserRouter as Router, Route } from "react-router-dom"
import { useStyletron } from "baseui"

const Home = React.lazy(() => import("./home"))
const Portal = React.lazy(() => import("./portal"))
const EmailVerify = React.lazy(() => import("../emailVerify"))
const ResetPassword = React.lazy(() => import("../resetPassword"))
const FieldappDownload = React.lazy(() => import("./fieldappDownload"))

const AdminRoutes = () => {
	const [css, theme] = useStyletron()
	const routeStyle: string = css({
		width: "100%",
		minHeight: "100vh",
	})
	return (
		<AuthContainer.Provider>
			<MetaTags>
				<title>Genesis</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<meta id="meta-description" name="description" content="Steaks as a Service" />
				<meta id="og-title" property="og:title" content="Genesis" />
			</MetaTags>

			<div className={routeStyle}>
				<Router>
					<Route path="/" exact component={Home} />
					<Route path="/verify/:code" exact component={EmailVerify} />
					<Route path="/verify" exact component={EmailVerify} />
					<Route path="/reset/:token" exact component={ResetPassword} />
					<Route path="/download" exact component={FieldappDownload} />
					<Route path="/portal" component={Portal} />
				</Router>
			</div>
		</AuthContainer.Provider>
	)
}

export default AdminRoutes
