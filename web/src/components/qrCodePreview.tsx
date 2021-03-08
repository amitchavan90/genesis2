import * as React from "react"
import { useStyletron } from "baseui"
import { useQuery } from "@apollo/react-hooks"
import { graphql } from "../graphql"
import { Settings } from "../types/types"
import { LoadingSimple } from "../components/loading"
import { QRCode } from "react-qrcode-logo"

export const QRPreview = (props: { item?: any }) => {
	const { item } = props

	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)

	const [css] = useStyletron()
	const qrItemsContainer = css({
		display: "flex",
		justifyContent: "space-evenly",
		paddingTop: "5px",
	})
	const qrItem = css({
		display: "flex",
		textAlign: "center",
		justifyContent: "center",
		flexDirection: "column",
		alignItems: "center",
		padding: "10px",
	})
	const qrLink = css({
		fontSize: "9px",
		lineHeight: "12px",
		color: "grey",
	})

	if (settingsQuery.loading) return <LoadingSimple />
	if (settingsQuery.error) return <p>Error: {settingsQuery.error.message}</p>
	if (!settingsQuery.data) return <p>Failed to fetch settings</p>

	const viewLink = item ? `${settingsQuery.data.settings.adminHost}/q/${item.id}` : ""

	return (
		<div className={qrItemsContainer}>
			<div className={qrItem}>
				<QRCode value={viewLink} size={250} quietZone={2} />
				<a href={viewLink} className={qrLink} target="_blank">
					{viewLink}
				</a>
			</div>
		</div>
	)
}
