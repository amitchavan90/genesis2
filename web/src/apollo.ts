import * as React from "react"
import { ApolloLink } from "apollo-link"
import { WebSocketLink } from "apollo-link-ws"
import { SubscriptionClient } from "subscriptions-transport-ws"
import { createUploadLink } from "apollo-upload-client"
import { onError } from "apollo-link-error"
import { getMainDefinition } from "apollo-utilities"

/** Additional HTTP Error Handling outside of GraphQL Response */
const errorHttpLink = onError(({ networkError }) => {
	// Fix up JSON parse error...
	if (networkError && "statusCode" in networkError && networkError.statusCode == 413) {
		networkError.message = "413 Request Entity Too Large"
	}
})

/** Apollo Link Splitter context for accessing the auth features. */
export const ApolloLinkSplitterContext = React.createContext<ApolloLinkSplitter | null>(null)

/** Apollo Link Splitter to handle auth and file upload. */
export class ApolloLinkSplitter {
	link: ApolloLink
	wsClient: SubscriptionClient

	/**
	 * Setup Apollo Link Splitter.
	 *
	 * @param endPoint
	 * @param isSecure
	 */
	constructor(endPoint: string, isSecure: boolean) {
		let httpEndpoint: string
		let wsEndpoint: string

		if (isSecure) {
			httpEndpoint = "https://" + endPoint
			wsEndpoint = "wss://" + endPoint
		} else {
			httpEndpoint = "http://" + endPoint
			wsEndpoint = "ws://" + endPoint
		}

		// Setup Apollo HTTP Link with Uploading Support
		const uploadLink = ApolloLink.from([errorHttpLink, createUploadLink({ uri: httpEndpoint })])

		// Setup Apollo Split Link
		this.wsClient = new SubscriptionClient(wsEndpoint, {
			reconnect: true,
		})
		const wsLink = new WebSocketLink(this.wsClient)

		this.link = ApolloLink.split(
			({ query }) => {
				const definition = getMainDefinition(query)
				return definition.kind === "OperationDefinition" && definition.operation === "subscription"
			},
			wsLink,
			uploadLink,
		)
	}

	/** Reset Websocket Link. */
	resetLink = () => this.wsClient.close(true, true)

	/**
	 * Return a Apollo Link with Splitting features.
	 *
	 * @returns
	 *     An ApolloLink object for ApolloClient
	 */
	getLink = () => this.link
}

/** Deletes queries from the cache to force a refetch on components that use it.
 *
 * Used when creating or archiving.
 *
 * `update: (cache: any) => invalidateListQueries(cache, "products"),`
 */
export const invalidateListQueries = (cache: any, queryName: string) => {
	Object.keys(cache.data.data).forEach(key => key.startsWith("$ROOT_QUERY." + queryName) && cache.data.delete(key))
}
