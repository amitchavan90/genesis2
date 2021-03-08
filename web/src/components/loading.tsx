import * as React from "react"
import { Spinner } from "baseui/spinner"
import { useStyletron } from "baseui"
import { H6 } from "baseui/typography"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"

export const Loading = () => {
	const [css, theme] = useStyletron()
	const containerStyle: string = css({
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		height: "100vh",
	})
	const contentStyle: string = css({
		display: "flex",
		flexDirection: "column",
		justifyContent: "center",
		alignItems: "center",
	})
	const balloonStyle: string = css({
		paddingTop: "20px",
		paddingBottom: "20px",
	})

	return (
		<div className={containerStyle}>
			<div className={contentStyle}>
				<div className={balloonStyle}>
					<FontAwesomeIcon icon={["fas", "steak"]} size="8x" />
				</div>
				<Spinner />
				<H6>Loading...</H6>
			</div>
		</div>
	)
}

export const LoadingSimple = () => {
	const [css, theme] = useStyletron()
	const containerStyle: string = css({
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		height: "100vh",
	})
	const contentStyle: string = css({
		display: "flex",
		flexDirection: "column",
		justifyContent: "center",
		alignItems: "center",
	})

	return (
		<div className={containerStyle}>
			<div className={contentStyle}>
				<Spinner />
				<H6>Loading...</H6>
			</div>
		</div>
	)
}
