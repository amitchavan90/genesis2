import * as React from "react"
import { useStyletron } from "baseui"
import { useHistory } from "react-router-dom"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { Button } from "baseui/button"
import { GetItemIcon } from "../themeOverrides"

interface SmallItemLinkProps {
	title?: string
	code: string
	link: string
	/** Will not append /${code} to end of link */
	fullLink?: boolean
	hash?: string
	icon?: IconName
	iconLight?: boolean
	itemName?: string
	disabled?: boolean
}
export const SmallItemLink = (props: SmallItemLinkProps) => {
	const { title, code, link, fullLink, hash, itemName, disabled } = props
	const [css, theme] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "5px",
		whiteSpace: "nowrap",
	})

	const icon = props.icon ? { icon: props.icon, light: props.iconLight } : GetItemIcon(itemName || "", true)

	const history = useHistory()

	return (
		<div
			onClick={
				disabled
					? undefined
					: e => {
							e.stopPropagation()
							if (fullLink) window.open(link, "_blank")
							else history.push(`${link}/${code}${hash ? `#${hash}` : ""}`)
					  }
			}
		>
			<Button
				kind="minimal"
				type="button"
				overrides={{
					BaseButton: {
						style: {
							...theme.typography.font200,
							paddingLeft: "2px",
							paddingRight: "2px",
							paddingTop: "2px",
							paddingBottom: "2px",
							color: "inherit",
							":hover": disabled
								? {}
								: {
										color: theme.colors.colorPrimary,
								  },
							marginRight: "0.5rem",
						},
					},
				}}
			>
				<FontAwesomeIcon icon={[icon.light ? "fal" : "fas", icon.icon]} size="sm" />
				<div className={marginLeftStyle}>{title || code}</div>
			</Button>
		</div>
	)
}
