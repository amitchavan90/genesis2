import * as React from "react"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useStyletron } from "baseui"
import { StyleObject } from "styletron-react"
import { paddingZero } from "../themeOverrides"
import { IconName } from "@fortawesome/fontawesome-svg-core"

interface MenuOptionLabelProps {
	label: string
	checked?: boolean
	borderTop?: boolean
	icon?: IconName
	iconLight?: boolean
}

export const MenuOptionLabel = (props: MenuOptionLabelProps) => {
	const { label, checked, borderTop, icon, iconLight } = props

	const [css, theme] = useStyletron()
	const menuItem = css({
		display: "flex",
		padding: "8px 16px",
	})
	const borderTopStyle = css({
		borderTop: "solid rgba(0, 0, 0, 0.2) 1px",
	})
	const checkIconStyle = css({
		marginRight: "10px",
		width: "16px",
	})
	const iconStyle = css({
		marginRight: "10px",
		width: "20px",
		height: "20px",
		textAlign: "center",
	})

	return (
		<div className={menuItem + (borderTop ? ` ${borderTopStyle}` : "")}>
			{!icon && <div className={checkIconStyle}>{checked && <FontAwesomeIcon icon={["fas", "check"]} />}</div>}
			{icon && (
				<div className={iconStyle}>
					<FontAwesomeIcon icon={[iconLight ? "fal" : "fas", icon]} />
				</div>
			)}
			<div>{label}</div>
		</div>
	)
}

export const MenuListOverride = {
	style: {
		outline: "unset",
		boxShadow: "unset",
	},
}
export const MenuOptionStyleOverride: StyleObject = {
	color: "black",
	":hover": {
		backgroundColor: "#1A1A1A",
		color: "white",
	},
	transitionTimingFunction: "unset",
	transitionDuration: "unset",
	transitionProperty: "unset",
	fontSize: "16px",
	lineHeight: "24px",
	...paddingZero,
}
