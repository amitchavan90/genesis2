import * as React from "react"
import { useStyletron } from "baseui"
import QRScanIcon from "../../assets/images/qr-scan.svg"

export const PressHereToScan = (props: { active: boolean; onClick: () => void }) => {
	const { active, onClick } = props

	const [css] = useStyletron()
	const container = css({
		opacity: active ? 1 : 0,
		visibility: active ? "visible" : "hidden",
		cursor: active ? "pointer" : "default",
		transition: "opacity 0.4s",
		userSelect: "none",

		display: "flex",
		justifyContent: "center",

		width: "100%",
		position: "absolute",
		top: "345px",
		left: "0px",
	})
	const box = css({
		fontFamily: "Elephant Black",
		borderRadius: "4px",
		backgroundColor: "white",
		padding: "15px",
		boxShadow: "#00000029 0px 3px 6px",

		width: "300px",
	})
	const textStyle = css({
		fontSize: "32px",
		textAlign: "center",
	})
	const flexStyle = css({
		display: "flex",
		justifyContent: "center",
		marginTop: "5px",
	})
	const subTextStyle = css({
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		fontSize: "14px",
		marginTop: "12px",
		whiteSpace: "nowrap",
	})
	const twoCircle = css({
		backgroundColor: "white",
		borderRadius: "15px",
		width: "30px",
		height: "30px",
		lineHeight: "34px",
		fontSize: "20px",
		marginLeft: "10px",
		color: "rgb(24 133 189)",
		textAlign: "center",
		alignSelf: "center",
	})
	const twoCircleSmall = css({
		backgroundColor: "#59B7E8",
		borderRadius: "18px",
		width: "18px",
		height: "18px",
		lineHeight: "20px",
		fontSize: "12px",
		marginLeft: "5px",
		marginRight: "5px",
		color: "white",
		textAlign: "center",
		alignSelf: "center",
	})
	const qrScanIconStyle = css({
		marginLeft: "5px",
		height: "24px",
	})
	const subButton = css({
		backgroundColor: "rgb(24 133 189)",
		color: "white",
		padding: "15px",
		borderRadius: "3px",
		boxShadow: "rgba(255, 255, 255, 0.5) 0.2em 0.2em 0.2em 0px inset, rgba(0, 0, 0, 0.5) -0.2em -0.2em 0.2em 0px inset, rgba(0, 0, 0, 0.3) 0px 3px 5px 2px",
	})

	return (
		<div className={container}>
			<div className={box} onClick={active ? onClick : undefined}>
				<div className={subButton + " shimmer"}>
					<div className={textStyle}>PRESS HERE</div>
					<div className={flexStyle}>
						<div className={textStyle}>TO SCAN QR</div>
						<div className={twoCircle}>2</div>
					</div>
				</div>
				<div className={subTextStyle}>
					<span>SCAN QR</span>
					<div className={twoCircleSmall}>2</div>
					<span>HIDDEN UNDER YOUR STEAK</span>
					<img src={QRScanIcon} className={qrScanIconStyle} />
				</div>
			</div>
		</div>
	)
}
