import * as React from "react"
import { styled, useStyletron } from "baseui"
import { LatestTransactionInfo } from "../types/types"

export const CenteredPage = styled("div", {
	width: "95%",
	maxWidth: "1200px",
	margin: "0 auto 1rem",
})

export const LatestTrackActionColumn = ({ value }: { value: LatestTransactionInfo }) => {
	const [css] = useStyletron()
	const fullDateStyle = css({
		fontSize: "0.8rem",
		lineHeight: "0.8rem",
		color: "grey",
	})

	if (!value) return null

	return (
		<div>
			<div>{value.name}</div>
			<div className={fullDateStyle}>{`(${new Date(value.createdAt).toLocaleString()})`}</div>
		</div>
	)
}
