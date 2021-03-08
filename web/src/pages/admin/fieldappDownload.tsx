import * as React from "react"
import { useStyletron } from "baseui"
import { graphql } from "../../graphql"
import { useQuery } from "@apollo/react-hooks"
import { useHistory } from "react-router-dom"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { QRCode } from "react-qrcode-logo"

const Logo = require("../../assets/images/Genesis_white.png")
const Screenshot = require("../../assets/images/fieldapp_screenshot.png")

const FieldappDownload = () => {
	const settingsQuery = useQuery<{ settings: { fieldappVersion: string } }>(graphql.query.FIELDAPP_VERSION)

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
	const logoStyle = css({
		width: "70%",
		maxWidth: "530px",
		marginBottom: "40px",
		"@media screen and (max-height: 840px)": {
			marginTop: "20px",
		},
	})
	const boxStyle = css({
		...theme.typography.ParagraphMedium,
		padding: "2rem",
		borderWidth: "1px",
		borderStyle: "solid",
		borderColor: theme.colors.primary,
		backgroundColor: "white",
		width: "70%",
		maxWidth: "600px",
	})
	const titleStyle = css({
		...theme.typography.HeadingMedium,
		paddingBottom: "10px",
	})
	const infoStyle = css({
		margin: "0 10px",
		width: "100%",
		textAlign: "right",
	})
	const downloadButtonStyle = css({
		backgroundColor: theme.colors.colorPrimary,
		color: "white",
		padding: "14px 16px",
		cursor: "pointer",
		display: "block",
		textDecoration: "unset",
		textAlign: "center",
		lineHeight: "20px",
		fontWeight: "bold",
		":hover": {
			backgroundColor: "#333333",
		},
	})
	const versionStyle = css({
		marginRight: "5px",
	})
	const screenshotStyle = css({
		maxHeight: "500px",
	})
	const spreadStyle = css({
		display: "flex",
		justifyContent: "space-between",
		alignItems: "center",
		"@media screen and (max-width: 630px)": {
			flexDirection: "column",
		},
	})
	const linkStyle = css({
		...theme.typography.LabelMedium,
		position: "absolute",
		left: "0",
		bottom: "0",
		color: "white",
		cursor: "pointer",
		":hover": {
			color: "#84aaff",
		},
		backgroundColor: theme.colors.colorPrimary,
		width: "100%",
		"@media screen and (max-height: 840px)": {
			position: "unset",
		},
	})
	const linkInnerStyle = css({
		padding: "10px",
		display: "flex",
		alignItems: "center",
	})
	const linkIconStyle = css({
		marginRight: "0.5rem",
	})
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

	const apkURL =
		settingsQuery.data?.settings !== undefined ? `${window.location.origin}/fieldapp/genesis-${settingsQuery.data.settings.fieldappVersion}.apk` : undefined

	return (
		<div className={container}>
			<img src={Logo} alt="Genesis" className={logoStyle} />

			<div className={boxStyle}>
				<div>
					<div className={titleStyle}>Genesis Field Application</div>

					<div className={spreadStyle}>
						<div className={infoStyle}>
							<div className={qrItemsContainer}>
								<div className={qrItem}>
									<QRCode value={apkURL} size={250} quietZone={2} />
								</div>
							</div>
							<a className={downloadButtonStyle} href={apkURL}>
								Download APK
							</a>
							{settingsQuery.data?.settings !== undefined && <div className={versionStyle}>v{settingsQuery.data.settings.fieldappVersion}</div>}
						</div>

						<img src={Screenshot} className={screenshotStyle} />
					</div>
				</div>
			</div>

			<div className={linkStyle} onClick={() => history.push("/")}>
				<div className={linkInnerStyle}>
					<FontAwesomeIcon icon={["fas", "user"]} className={linkIconStyle} />
					<div>Admin Login</div>
				</div>
			</div>
		</div>
	)
}

export default FieldappDownload
