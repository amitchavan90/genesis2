import * as React from "react"
import MetaTags from "react-meta-tags"
import { ApolloProvider } from "@apollo/react-hooks"
import { ApolloClient, InMemoryCache } from "apollo-boost"
import { Client as Styletron } from "styletron-engine-atomic"
import { Provider as StyletronProvider } from "styletron-react"
import { BaseProvider } from "baseui"
import { LightTheme } from "./themeOverrides"
import { loadIcons } from "./helpers/loadicons"
import { ApolloLinkSplitter, ApolloLinkSplitterContext } from "./apollo"
import { Enviroment } from "./types/enums"
import { Loading } from "./components/loading"
import "react-vertical-timeline-component/style.min.css"
const AdminRoutes = React.lazy(() => import("./pages/admin/routes"))
const ConsumerRoutes = React.lazy(() => import("./pages/consumer/routes"))

const engine = new Styletron()

loadIcons()

var loc = window.location
const GRAPHQL_ENDPOINT = "//" + loc.host + "/api/gql/query"
const apolloLinkSplitter = new ApolloLinkSplitter(GRAPHQL_ENDPOINT, loc.protocol === "https:")
const apolloClient = new ApolloClient({
	cache: new InMemoryCache(),
	link: apolloLinkSplitter.getLink(),
})

// /** For QR links */
// export const API_URL = (() => {
// 	if (loc.host.indexOf("staging.theninja.life") != -1) return "https://genesis.staging.theninja.life/api"
// 	return `http://${loc.host}/api`
// })()
// /** The url for the consumer view site (:8082 or consumer.genesis) */
// export const CONSUMER_URL = (() => {
// 	if (loc.host.indexOf("staging.theninja.life") != -1) return "https://consumer.genesis.staging.theninja.life"
// 	return `http://${loc.host.substring(0, loc.host.indexOf(":"))}:8082`
// })()

export const resetClient = () => {
	apolloClient.resetStore()
	apolloLinkSplitter.resetLink()
}

const App = (props: { enviroment: Enviroment }) => {
	return (
		<StyletronProvider value={engine}>
			<BaseProvider theme={LightTheme}>
				<ApolloProvider client={apolloClient}>
					<ApolloLinkSplitterContext.Provider value={apolloLinkSplitter}>
						<React.Suspense fallback={<Loading />}>
							{props.enviroment == Enviroment.Admin && <AdminRoutes />}
							{props.enviroment == Enviroment.Consumer && <ConsumerRoutes />}
						</React.Suspense>
					</ApolloLinkSplitterContext.Provider>
				</ApolloProvider>
			</BaseProvider>
		</StyletronProvider>
	)
}

export { App }
