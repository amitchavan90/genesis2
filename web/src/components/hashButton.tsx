import * as React from "react"
import { useStyletron } from "baseui"
import { graphql } from "../graphql"
import { useQuery } from "@apollo/react-hooks"
import { Spinner } from "baseui/spinner"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

export const HashButton = (props: { hash?: string; isLoading?: boolean; isAddress?: boolean; hashLength?: number; alignLeft?: boolean }) => {
	const { hash, isLoading, isAddress, alignLeft } = props
	const hashLength = props.hashLength || 8

	const settingsQuery = useQuery<{ settings: { etherscanHost: string } }>(graphql.query.SETTINGS)

	const truncate = (input: string) => {
		if (input.length > hashLength) return input.substring(0, hashLength) + "..."
		return input
	}

	// Styling
	const [css, theme] = useStyletron()
	const transactionLink = css({
		cursor: "pointer",
		fontFamily: "monospace",
		marginLeft: alignLeft ? "auto" : "unset",
		display: "flex",
		":hover": {
			color: theme.colors.primary,
			backgroundColor: "white",
		},
		padding: "2px 4px",
		borderRadius: "2px",
		width: "fit-content",
		float: "left",
		textDecoration: "none",
		":visited": {
			color: theme.colors.primary,
		},
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

	if (hash === "-") {
		// no blockchain hash
		return <div className={pendingStyle}>-</div>
	}

	if (!hash) {
		// pending transaction
		return <div className={pendingStyle}>Pending</div>
	}

	return (
		<a
			className={transactionLink}
			href={settingsQuery.data?.settings?.etherscanHost ? `${settingsQuery.data.settings.etherscanHost}/${isAddress ? "address" : "tx"}/${hash}` : undefined}
			target="_blank"
		>
			<div className={transactionIcon}>{isLoading ? <Spinner size="15px" /> : <FontAwesomeIcon icon={["fas", "external-link"]} />}</div>
			{truncate(hash)}
		</a>
	)
}
