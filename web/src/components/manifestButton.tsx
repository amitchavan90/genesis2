import * as React from "react"
import { useStyletron } from "baseui"
import { graphql } from "../graphql"
import { useQuery } from "@apollo/react-hooks"
import { Spinner } from "baseui/spinner"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

export const ManifestButton = (props: {
	manifestLineSha256?: string
	isLoading?: boolean
	isAddress?: boolean
	manifestLineSha256Length?: number
}) => {
	const { manifestLineSha256, isLoading, isAddress } = props
	const manifestLineSha256Length = props.manifestLineSha256Length || 8

	const settingsQuery = useQuery<{ settings: { etherscanHost: string } }>(graphql.query.SETTINGS)

	const truncate = (input: string) => {
		if (input.length > manifestLineSha256Length) return input.substring(0, manifestLineSha256Length) + "..."
		return input
	}

	// Styling
	const [css, theme] = useStyletron()
	const transactionLink = css({
		cursor: "pointer",
		fontFamily: "monospace",
		marginLeft: "auto",
		display: "flex",
		":hover": {
			color: theme.colors.primary,
			backgroundColor: "white",
		},
		padding: "2px 4px",
		borderRadius: "2px",
		width: "fit-content",
		float: "left",
	})
	const transactionIcon = css({
		marginRight: "0.5rem",
	})
	const pendingStyle = css({
		fontFamily: "monospace",
		marginLeft: "auto",
		display: "flex",
		padding: "2px 4px",
		borderRadius: "2px",
		width: "fit-content",
		float: "left",
	})

	if (!manifestLineSha256) {
		// pending transaction
		return <div className={pendingStyle}>Pending</div>
	}

	return (
		<div
			className={transactionLink}
			onClick={() => {
				if (!settingsQuery.data) return
				window.open(`/api/manifest/line/${manifestLineSha256}`, "_blank")
			}}>
			<div className={transactionIcon}>
				{isLoading ? <Spinner size="15px" /> : <FontAwesomeIcon icon={["fas", "external-link"]} />}
			</div>
			{truncate(manifestLineSha256)}
		</div>
	)
}
