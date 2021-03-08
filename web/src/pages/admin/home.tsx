import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { Login } from "../../components/login"
import { Spaced } from "../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
const { version } = require("../../../package.json")
const { build_start } = require("../../../stats.json")

const Logo = require("../../assets/images/Genesis_white.png")

interface IProps extends RouteComponentProps {}

const Home = (props: IProps) => {
	const history = useHistory()

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
	const homeStyle = css({
		padding: "2rem",
		borderWidth: "1px",
		borderStyle: "solid",
		borderColor: theme.colors.primary,
		backgroundColor: "white",
		width: "70%",
		maxWidth: "600px",
	})
	const logoStyle: string = css({
		width: "70%",
		maxWidth: "530px",
		marginBottom: "40px",
	})

	const versionStyle: string = css({
		...theme.typography.LabelSmall,
		position: "absolute",
		right: "6px",
		bottom: "10px",
		color: "white",
		textAlign: "right",
	})
	const buildTimeStyle: string = css({
		color: "grey",
	})

	const downloadLinkStyle: string = css({
		...theme.typography.LabelMedium,
		position: "absolute",
		left: "6px",
		bottom: "10px",
		color: "white",
		cursor: "pointer",
		":hover": {
			color: "#84aaff",
		},
	})

	return (
		<div className={container}>
			<img src={Logo} alt="Genesis" className={logoStyle} />
			<div className={homeStyle}>
				<Login hideCancel />
			</div>

			<div className={downloadLinkStyle} onClick={() => history.push("/download")}>
				<Spaced>
					<FontAwesomeIcon icon={["fas", "file-download"]} />
					<div>Download Fieldapp</div>
				</Spaced>
			</div>

			<div className={versionStyle}>
				v{version}
				<br />
				<div className={buildTimeStyle}>{build_start || "-"}</div>
			</div>
		</div>
	)
}

export default Home
