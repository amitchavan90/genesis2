import * as React from "react"
import { useStyletron } from "baseui"
import Confetti from "react-dom-confetti"
import { PressHereToScan } from "../../components/consumer/pressHereToScanButton"
import { LightBlue } from "../../themeOverrides"
import * as Scroll from "react-scroll"
import DragSteak, { FloatingSteak } from "./steak"

import LandingBackBoard from "../../assets/images/landing_back_board.svg"
import StickyLandingBackBoard from "../../assets/images/backboard-sticky.svg"

import Kangaroo from "../../assets/images/kangaroo.png"
import Plane from "../../assets/images/plane.png"
import Wine from "../../assets/images/wine.png"
import WineGlass from "../../assets/images/wine_glass.png"
import WineSingle from "../../assets/images/wine_single.png"
import WineLiquid from "../../assets/images/wine_liquid.svg"
import WineGlassMask from "../../assets/images/wine_glass_mask.png"

const scroll = Scroll.animateScroll

const confettiConfig = {
	angle: 90,
	spread: 55,
	startVelocity: 60,
	elementCount: 120,
	dragFriction: 0.12,
	duration: 3000,
	stagger: 3,
	width: "10px",
	height: "10px",
	colors: ["#a864fd", "#29cdff", "#78ff44", "#ff718d", "#fdff6a"],
}

const SteakAndGrill = () => {
	// QR Scanning

	const [steakMoved, setSteakMoved] = React.useState(false)

	const [qrError, setQrError] = React.useState<string>()
	const startScan = () => {
		setQrError(undefined)
		const w: any = window
		w.wx.error((res: any) => setQrError(res.errMsg))
		w.wx.miniProgram.redirectTo({ url: "../scan/index?page=" + encodeURIComponent(location.href) }) //Redirect to scan page in MP
	}

	React.useEffect(() => {
		const script = document.createElement("script")
		script.src = `https://res.wx.qq.com/open/js/jweixin-1.3.2.js`
		document.body.appendChild(script)

		return () => {
			document.body.removeChild(script)
		}
	}, [])

	// Animation

	// Scaling and scrolling
	const [scrollTop, setScrollTop] = React.useState(window.pageYOffset)
	const [scrolledPassedSteak, setScrolledPassed] = React.useState(window.pageYOffset > 600)
	const [pageScale, setPageScale] = React.useState(window.innerWidth > 400 ? 1 : window.innerWidth / 400)
	const onResize = () => setPageScale(window.innerWidth > 400 ? 1 : window.innerWidth / 400)
	React.useEffect(() => {
		const onScroll = (e: any) => {
			setScrollTop(window.pageYOffset)
		}
		window.addEventListener("scroll", onScroll)
		window.addEventListener("resize", onResize)
		return () => {
			window.removeEventListener("resize", onResize)
			window.removeEventListener("scroll", onScroll)
		}
	}, [])

	React.useEffect(() => {
		if ((scrollTop > 600 && !scrolledPassedSteak) || (scrolledPassedSteak && scrollTop > 400)) {
			setScrolledPassed(true)
		} else if (scrollTop <= 600 && scrolledPassedSteak === true) {
			setScrolledPassed(false)
		}
	}, [scrollTop])

	// Styling
	const [css, theme] = useStyletron()
	const scallableStyle = css({
		alignSelf: "center",
		userSelect: "none",
	})
	const headingStyle = css({
		padding: "0 26px",
		display: "flex",
	})
	const titleStyle = css({
		fontFamily: "Elephant Black",
		fontSize: "27px",
		color: LightBlue,
		marginTop: "23px",
		marginBottom: "23px",
		maxWidth: "320px",
		textTransform: "uppercase",
	})
	const earnPointsToWinStyle = css({
		color: LightBlue,
		margin: "23px 0",
		marginBottom: "0",
	})

	const backboardStyle = css({
		display: "flex",
		width: "380px",
		margin: "auto",
		pointerEvents: "none",
		userSelect: "none",
		perspective: "unset",
		zIndex: -1,
		[`@supports (-webkit-touch-callout: none) `]: {
			perspective: "unset !important",
		},
	})
	const confettiContainerStyle = css({
		display: "flex",
		justifyContent: "center",
	})
	const floatingTopBar = css({
		alignSelf: "center",
		userSelect: "none",

		zIndex: 20,
		position: "fixed",
		top: scrolledPassedSteak ? "-340px" : "-460px",
		left: scrolledPassedSteak ? "50%" : "unset",
		transform: scrolledPassedSteak && !steakMoved ? "translateX(-50%)" : "unset",
		transition: "top 0.25s",
	})

	return (
		<>
			{!steakMoved && (
				<div
					className={floatingTopBar}
					onClick={() => {
						if (scrolledPassedSteak && !steakMoved) scroll.scrollTo(100)
					}}
				>
					<img src={StickyLandingBackBoard} className={backboardStyle + " " + "un-perspectiver"} />
					<FloatingSteak />
				</div>
			)}

			<div className={scallableStyle}>
				<div className={headingStyle}>
					<div>
						<h1 className={titleStyle}>100% Authenticated Beef Blockchain</h1>

						<div style={{ maxWidth: "225px" }}>Drag the steak into the pan to close this Blockchain</div>
						<div className={earnPointsToWinStyle}>Earn Points and Win!</div>
					</div>
				</div>

				{steakMoved && <CelebrationAnimations />}

				<img src={LandingBackBoard} className={backboardStyle + " " + "un-perspectiver"} />

				<PressHereToScan active={steakMoved} onClick={startScan} />
				<div className={confettiContainerStyle}>
					<Confetti active={steakMoved} config={confettiConfig} />
				</div>

				<DragSteak pageScale={pageScale} steakMoved={steakMoved} setSteakMoved={setSteakMoved} />
			</div>
		</>
	)
}

const CelebrationAnimations = () => {
	const [css, theme] = useStyletron()

	const planeStyle = css({
		position: "absolute",
		zIndex: 4,
		width: "250px",
		marginTop: "45px",
		marginLeft: "60px",
		userSelect: "none",

		animationDuration: "2s",
		animationName: {
			"0%": {
				transform: "translate(-75vw, -100px) rotate(25deg) scale(0)",
			},
			"100%": {
				transform: "translate(0vw, 0px) rotate(0deg) scale(1)",
			},
		} as any,
	})

	const kangarooStyle = css({
		position: "absolute",
		zIndex: 4,
		width: "145px",
		marginLeft: "-10px",
		transform: "translate(0vw, 325px)",
		userSelect: "none",
		animationDuration: "2s",
		animationName: {
			"0%": {
				transform: "translate(-45vw, 325px) scale(0) rotate(0deg)",
				opacity: 0,
			},
			"15%": {
				transform: "translate(-45vw, 325px) scale(0.3) rotate(0deg)",
				opacity: 1,
			},

			"20%": {
				transform: "translate(-42vw, 365px) scale(0.3) rotate(-10deg)",
				opacity: 1,
			},
			"25%, 28%": {
				transform: "translate(-40vw, 305px) scale(0.4) rotate(10deg)",
				opacity: 1,
			},
			"40%": {
				transform: "translate(-31vw, 365px) scale(0.5) rotate(-10deg)",
				opacity: 1,
			},
			"45%, 48%": {
				transform: "translate(-28vw, 305px) scale(0.6) rotate(10deg)",
				opacity: 1,
			},
			"60%": {
				transform: "translate(-19vw, 365px) scale(0.7) rotate(-10deg)",
				opacity: 1,
			},
			"65%, 68%": {
				transform: "translate(-16vw, 305px) scale(0.8) rotate(10deg)",
				opacity: 1,
			},
			"80%": {
				transform: "translate(-8vw, 365px) scale(0.9) rotate(-10deg)",
				opacity: 1,
			},
			"85%, 88%": {
				transform: "translate(-5vw, 305px) scale(0.95) rotate(10deg)",
				opacity: 1,
			},

			"100%": {
				transform: "translate(0vw, 325px) scale(1) rotate(0deg)",
				opacity: 1,
			},
		} as any,
	})

	const wineStyle = css({
		position: "absolute",
		zIndex: 4,
		userSelect: "none",
		width: "100px",
		transform: "translate(342px, 250px)",

		animationDuration: "3s",
		animationName: {
			"0%, 20%": {
				opacity: 0,
			},
			"40%, 100%": {
				opacity: 1,
			},
		} as any,
	})
	const wineGlassStyle = css({
		position: "absolute",
		zIndex: 7,
		userSelect: "none",
		width: "50px",
		transform: "translate(306px, 344px)",

		animationDuration: "3s",
		animationName: {
			"0%, 20%": {
				opacity: 0,
			},
			"40%, 100%": {
				opacity: 1,
			},
		} as any,
	})
	const wineGlassMaskStyle = css({
		position: "absolute",
		zIndex: 6,
		userSelect: "none",
		width: "50px",
		transform: "translate(308px, 365px)",

		animationDuration: "3s",
		animationName: {
			"0%, 20%": {
				opacity: 0,
				clipPath: "inset(100% 0 0 0)",
			},
			"40%, 87%": {
				opacity: 1,
				clipPath: "inset(100% 0 0 0)",
			},
			"100%": {
				opacity: 1,
				clipPath: "inset(0% 0 0 0)",
			},
		} as any,
	})
	const wineLiquidStyle = css({
		position: "absolute",
		zIndex: 6,
		userSelect: "none",
		transform: "translate(325px, 247px) scaleY(1.4)",
		opacity: 0,
		animationDuration: "3s",
		animationName: {
			"0%, 20%": {
				opacity: 0,
				clipPath: "inset(0 0 100% 0)",
			},
			"60%, 80%": {
				opacity: 1,
				clipPath: "inset(0 0 100% 0)",
			},
			"80%, 90%": {
				opacity: 1,
				clipPath: "inset(0 0 0% 0)",
			},
			"100%": {
				opacity: 0,
			},
		} as any,
	})
	const wineSingleStyle = css({
		position: "absolute",
		zIndex: 4,
		userSelect: "none",
		width: "170px",
		transform: "translate(342px, 155px) rotate(-40deg)",
		opacity: 0,
		animationDuration: "3s",
		animationName: {
			"0%, 20%": {
				transform: "translate(342px, 155px) rotate(0deg)",
				opacity: 0,
			},
			"40%": {
				transform: "translate(342px, 155px) rotate(0deg)",
				opacity: 1,
			},

			"80%, 90%": {
				transform: "translate(342px, 155px) rotate(-40deg)",
				opacity: 1,
			},
			"100%": {
				opacity: 0,
			},
		} as any,
	})
	const animationsStyle = css({
		pointerEvents: "none",
		position: "relative",
	})

	return (
		<div className={animationsStyle + " un-perspectiver"}>
			<img src={Plane} className={planeStyle} draggable={false} />
			<img src={Kangaroo} className={kangarooStyle} draggable={false} />

			<img src={Wine} className={wineStyle} draggable={false} />
			<img src={WineLiquid} className={wineLiquidStyle} draggable={false} />
			<img src={WineGlassMask} className={wineGlassMaskStyle} draggable={false} />
			<img src={WineGlass} className={wineGlassStyle} draggable={false} />
			<img src={WineSingle} className={wineSingleStyle} draggable={false} />
			<img src={Plane} className={planeStyle} draggable={false} />
			<img src={Kangaroo} className={kangarooStyle} draggable={false} />

			<img src={Wine} className={wineStyle} draggable={false} />
			<img src={WineLiquid} className={wineLiquidStyle} draggable={false} />
			<img src={WineGlassMask} className={wineGlassMaskStyle} draggable={false} />
			<img src={WineGlass} className={wineGlassStyle} draggable={false} />
			<img src={WineSingle} className={wineSingleStyle} draggable={false} />
		</div>
	)
}

export default SteakAndGrill
