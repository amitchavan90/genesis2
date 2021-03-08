import * as React from "react"
import { useStyletron } from "baseui"

const Logo = require("../../assets/images/Genesis_white.png")

const Home = () => {
	const [css, theme] = useStyletron()
	const container = css({
		minHeight: "100vh",
		width: "100%",
		backgroundColor: theme.colors.colorPrimary,
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		flexDirection: "column",
	})
	const logoStyle: string = css({
		width: "70%",
		maxWidth: "530px",
		marginBottom: "40px",
	})

	return (
		<div className={container}>
			<img src={Logo} alt="Genesis" className={logoStyle} />
		</div>
	)
}

export default Home
