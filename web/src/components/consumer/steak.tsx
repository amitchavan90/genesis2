import * as React from "react"
import { styled } from "baseui"
import Draggable, { DraggableData, DraggableEvent } from "react-draggable"

import Steak from "../../assets/images/steak.png"
import Hand from "../../assets/images/hand.svg"

import Smoke from "../../assets/images/elshmoko.png"
import Grill from "../../assets/images/grill.png"
import GrillCircle from "../../assets/images/grill_circle.svg"

const steakTop = 45
const steakBottom = 420

const SteakContainer = styled("div", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	position: "absolute",
	marginLeft: "-6px",
	zIndex: 4,
	top: "225px",
	left: "calc(50% - 120px)",
})
const StuckContainer = styled("div", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	position: "absolute",
	top: "140px",
	left: "50%",
	transform: "translateX(-50%)",
})
const SteakImage = styled("img", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	maxWidth: "240px",
	width: "100%",
	userSelect: "none",
})
const SteakAniImg = styled("img", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	maxWidth: "240px",
	userSelect: "none",
	width: "100%",
	animationDuration: "3.5s",
	animationIterationCount: "infinite",
	WebkitBackfaceVisibility: "hidden",
	animationName: {
		"0%, 10%": {
			opacity: 0,
		},
		"65%, 80%": {
			opacity: 0.5,
		},
		"80%, 90%": {
			opacity: 0.5,
		},
		"100%": {
			opacity: 0,
		},
	} as any,
})

const AnimSteakDiv = styled("div", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	position: "relative",
	marginLeft: "-6px",
	zIndex: 3,
	top: "-401px",
	right: "-1px",
	animationDuration: "3.5s",
	animationIterationCount: "infinite",
	userSelect: "none",
	transform: "translate(0px, 0px)",
	WebkitBackfaceVisibility: "hidden",
	animationName: {
		"0%, 10%": {
			transform: "translate(0px, 0px) scale(1)",
			opacity: 1,
		},

		"65%, 80%": {
			transform: "translate(0px, 380px) scale(0.7)",
			opacity: 1,
		},
		"80%, 90%": {
			transform: "translate(0px, 380px) scale(0.7)",
			opacity: 1,
		},
		"100%": {
			transform: "translate(0px, 380px) scale(0)",
			opacity: 0,
		},
	} as any,
})
const HandImg = styled("img", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	width: "100px",
	position: "absolute",
	left: "55%",
	top: "50%",
	userSelect: "none",
	transform: "translate(-50%, -50%)",
	animationDuration: "3.5s",
	animationIterationCount: "infinite",
	WebkitBackfaceVisibility: "hidden",
	animationName: {
		"0%": {
			width: "100px",
		},
		"30%": {
			width: "60px",
		},
		"100%": {
			width: "60px",
		},
	} as any,
})

const Circle = styled("div", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	width: "100%",
	pointerEvents: "none",
	backgroundImage: `url(${GrillCircle})`,
	backgroundSize: "contain",
	backgroundRepeat: "no-repeat",
	textAlign: "center",
	display: "flex",
	justifyContent: "center",
	userSelect: "none",
})
const ImageGrill = styled("img", {
	MozUserSelect: "none",
	WebkitUserSelect: "none",
	msUserSelect: "none",
	maxWidth: "76%",
	width: "100%",
	marginTop: "25px",
	marginLeft: "-12px",
	pointerEvents: "none",
	userSelect: "none",
	height: "100%",
})
const ImageSmoke = styled("img", ({ steakY }: { steakY: number }) => {
	return {
		MozUserSelect: "none",
		WebkitUserSelect: "none",
		msUserSelect: "none",
		display: "flex",
		position: "absolute",
		width: "50%",
		maxWidth: "192px",
		pointerEvents: "none",
		mixBlendMode: "exclusion",
		opacity: Math.max((steakY - 85) / 375.0, 0),
		userSelect: "none",
	}
})

interface Props {
	pageScale: number
	setSteakMoved(moved: boolean): void
	steakMoved: boolean
}

export default ({ pageScale, steakMoved, setSteakMoved }: Props) => {
	const [steakY, setSteakY] = React.useState(steakTop)
	const [slideSteak, setSlideSteak] = React.useState(false)

	const onDrag = (e: DraggableEvent, data: DraggableData) => {
		if (data.y === steakY) return
		setSteakY(data.y)
		setIsHideSteak(true)
	}
	const onSteakClick = () => {
		if (steakMoved) return
		if (slideSteak) setSlideSteak(false)
	}
	const onDragStop = (e: DraggableEvent, data: DraggableData) => {
		const moved = data.y >= steakTop + (steakBottom - steakTop) / 2
		if (moved) setSteakMoved(true)
		setSlideSteak(true)
		setSteakY(moved ? steakBottom : steakTop)
	}

	// start of sticky banner on top
	const [scrollTop, setScrollTop] = React.useState(-1)
	const [stickSteak, setStickSteak] = React.useState(true)

	React.useEffect(() => {
		const onScroll = (e: any) => {
			if (scrollTop != -1) {
				setStickSteak(scrollTop - window.pageYOffset >= 0 || window.pageYOffset < 0)
			}
			setScrollTop(window.pageYOffset)
		}
		window.addEventListener("scroll", onScroll)

		return () => window.removeEventListener("scroll", onScroll)
	}, [scrollTop])

	// Hide animated steak
	var [isHideSteak, setIsHideSteak] = React.useState(false)
	const onTouchStart = (evt: TouchEvent) => setIsHideSteak(true)
	const onTouchEnd = (evt: TouchEvent) => setIsHideSteak(false)
	React.useEffect(() => {
		window.addEventListener("touchstart", onTouchStart)
		window.addEventListener("touchend", onTouchEnd)
	}, [])

	return (
		<>
			<Draggable
				axis="y"
				defaultPosition={{ x: 0, y: steakTop }}
				bounds={{ top: steakTop, bottom: steakBottom }}
				onDrag={onDrag}
				onStop={onDragStop}
				onMouseDown={onSteakClick}
				disabled={steakMoved}
				position={slideSteak ? { x: 0, y: steakY } : undefined}
				scale={pageScale}
			>
				<SteakContainer>
					<SteakImage
						src={Steak}
						draggable={false}
						style={{ transform: `scale(${1 - 0.3 * ((steakY - steakTop) / 375.0)})`, transition: slideSteak ? "0.5s" : "unset" }}
					/>
					{!steakMoved && (
						<AnimSteakDiv style={{ display: isHideSteak ? "none" : "block" }}>
							<SteakAniImg src={Steak} style={{ transform: `scale(${1 - 0.3 * ((steakY - steakTop) / 375.0)})` }} draggable={false} />
							<HandImg src={Hand} draggable={false} />
						</AnimSteakDiv>
					)}
				</SteakContainer>
			</Draggable>
			<Circle>
				<ImageSmoke src={Smoke} steakY={steakY} />
				<ImageGrill src={Grill} />
			</Circle>
		</>
	)
}

export const FloatingSteak = () => {
	return (
		<StuckContainer>
			<SteakImage src={Steak} draggable={false} />
		</StuckContainer>
	)
}
