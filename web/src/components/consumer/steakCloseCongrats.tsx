import * as React from "react"
import { useStyletron } from "baseui"

import SteakCooked from "../../assets/images/steak_cooked.png"
import Smoke from "../../assets/images/elshmoko.png"
import Grill from "../../assets/images/grill.png"
import GrillCircle from "../../assets/images/grill_circle_white.svg"
import Congrats from "../../assets/images/congrats.svg"
import L28Logo from "../../assets/images/l28_logo.svg"
import ProducedInAustralia from "../../assets/images/produced_in_australia.svg"
import TrueAussieBeef from "../../assets/images/true_aussie_beef.svg"

const SteakCloseCongrats = (props: { pointsGiven: string | null }) => {
	const { pointsGiven } = props

	// Styling
	const [css, theme] = useStyletron()
	const containerStyle = css({
		backgroundColor: "#85D3E1",
		textAlign: "center",
		justifyContent: "center",
		color: "white",
		display: "flex",
		flexDirection: "column",
	})
	const innerStyle = css({
		maxWidth: "400px",
		margin: "0 auto",
		padding: "0 26px",
		display: "flex",
	})
	const titleStyle = css({
		fontFamily: "Elephant Black",
		fontSize: "30px",
		margin: "23px auto",
		maxWidth: "320px",
		textTransform: "uppercase",
	})
	const grillContainerStyle = css({
		display: "flex",
		justifyContent: "center",
		width: "116vw",
		marginLeft: "-8vw",
		background: "linear-gradient(0deg, #FFFFFF 35%, #85D3E1 55%)",
	})
	const grillCircleStyle = css({
		pointerEvents: "none",
		backgroundImage: `url(${GrillCircle})`,
		backgroundSize: "contain",
		backgroundRepeat: "no-repeat",
		textAlign: "center",
		display: "flex",
		justifyContent: "center",
		userSelect: "none",

		width: "116vw",
		height: "116vw",
		maxWidth: "400px",
		maxHeight: "400px",
	})
	const grillStyle = css({
		position: "absolute",
		width: "70vw",
		maxWidth: "260px",
		marginLeft: "-6px",
		marginTop: "30px",
		userSelect: "none",
	})
	const steakStyle = css({
		position: "absolute",
		width: "40vw",
		maxWidth: "150px",
		marginLeft: "-6px",
		marginTop: "30px",
		userSelect: "none",
	})
	const smokeStyle = css({
		position: "absolute",
		width: "40vw",
		maxWidth: "150px",
		marginLeft: "-6px",
		marginTop: "30px",
		pointerEvents: "none",
		mixBlendMode: "exclusion",
		userSelect: "none",
	})
	const iconContainerStyle = css({
		display: "flex",
		justifyContent: "space-evenly",
		padding: "20px 0",
	})
	const iconContainerStyle2 = css({
		display: "flex",
		justifyContent: "center",
		padding: "20px 0",
	})
	const iconStyle = css({
		userSelect: "none",
	})
	const iconRightStyle = css({
		userSelect: "none",
		marginLeft: "15px",
	})
	const linkContainerStyle = css({
		marginTop: "15px",
		marginBottom: "25px",
	})
	const linkStyle = css({
		fontFamily: "Elephant Black",
		fontSize: "22px",
		textDecoration: "none",
		color: "white",
		":visited": {
			color: "white",
		},
	})

	return (
		<div className={containerStyle}>
			<div className={innerStyle}>
				<div>
					<h1 className={titleStyle}>THANK YOU</h1>

					<div>You have closed this Blockchain</div>
					{pointsGiven !== null && <div>{`and earned ${pointsGiven} point${pointsGiven === "1" ? "" : "s"}`}</div>}
					<div className={iconContainerStyle}>
						<img src={L28Logo} className={iconStyle} draggable={false} />
						<img src={Congrats} draggable={false} />
					</div>

					<p>
						Latitude28 is proud to bring quality Australian beef to China, using leading blockchain technology to track the journey of our products from paddock
						to plate. We’re proud to help connect people through food, bringing energy, fun and flavour to any table. Visit our website to find out more.
					</p>

					<div className={linkContainerStyle}>
						<a href="https://www.l28produce.cn" className={linkStyle}>
							WWW.L28PRODUCE.CN
						</a>
					</div>

					<div className={iconContainerStyle2}>
						<img src={TrueAussieBeef} className={iconStyle} draggable={false} />
						<img src={ProducedInAustralia} className={iconRightStyle} draggable={false} />
					</div>
				</div>
			</div>

			<div className={grillContainerStyle}>
				<div className={grillCircleStyle} />
				<img src={Grill} className={grillStyle} />
				<img src={Smoke} className={smokeStyle} />
				<img src={SteakCooked} className={steakStyle} />
			</div>
		</div>
	)
}

export default SteakCloseCongrats
