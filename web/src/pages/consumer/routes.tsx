import * as React from "react"
import MetaTags from "react-meta-tags"
import { AuthContainer } from "../../controllers/auth"
import { BrowserRouter as Router, Route } from "react-router-dom"
import { useStyletron } from "baseui"
import "slick-carousel/slick/slick.css"
import "slick-carousel/slick/slick-theme.css"
const Home = React.lazy(() => import("./home"))
const ProductLandingPage = React.lazy(() => import("./landingPage"))

const ConsumerRoutes = () => {
	const [css, theme] = useStyletron()
	const routeStyle: string = css({
		width: "100%",
		minHeight: "100vh",
	})
	return (
		<AuthContainer.Provider>
			<MetaTags>
				<title>L28 | 100% Authenticated Beef Blockchain</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<meta id="meta-description" name="description" content="100% Authenticated Beef Blockchain" />
				<meta id="og-title" property="og:title" content="L28" />
			</MetaTags>

			<div className={routeStyle}>
				<Router>
					<Route path="/" exact component={Home} />
					<Route path={"/view"} component={ProductLandingPage} />
				</Router>
			</div>
		</AuthContainer.Provider>
	)
}

export default ConsumerRoutes
