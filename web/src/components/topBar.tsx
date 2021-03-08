import * as React from "react"
import { useStyletron, styled } from "baseui"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { StatefulTooltip } from "baseui/tooltip"
import { PLACEMENT, StatefulPopover } from "baseui/popover"
import { StatefulMenu } from "baseui/menu"
import { UserContainer } from "../controllers/user"
import { AuthContainer } from "../controllers/auth"
import { useHistory } from "react-router-dom"
import { Button } from "baseui/button"
import { Block } from "baseui/block"
import { paddingZero, boxShadowCommon } from "../themeOverrides"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { VerifyBanner } from "./verifyBanner"

export const TopBar = () => {
	const [css] = useStyletron()

	const { user } = UserContainer.useContainer()

	const containerStyle: string = css({
		display: "flex",
		justifyContent: "flex-end",
		alignItems: "center",
		width: "100%",
		height: "70px",
		borderBottom: `1px solid #D9D9D9`,
		fontSize: "1.5rem",
		boxShadow: boxShadowCommon,
	})

	const topBarContentStyle: string = css({
		marginRight: "25px",
		display: "flex",
		alignItems: "center",
	})

	return (
		<>
			<div className={containerStyle}>
				<div className={topBarContentStyle}>
					<Account />
				</div>
			</div>
			{user && !user.verified && <VerifyBanner />}
		</>
	)
}

const Account = () => {
	const [css, theme] = useStyletron()
	const [isOpen, setIsOpen] = React.useState<boolean>(false)
	const { logout } = AuthContainer.useContainer()
	const { user } = UserContainer.useContainer()

	const back = css({
		backgroundColor: "white",
	})
	const userStyle: string = css({
		cursor: "pointer",
		display: "flex",
		alignItems: "center",
	})

	const userCircleStyle: string = css({
		marginRight: "0.5rem",
	})

	if (!user) {
		return null
	}

	const menuItems = (close: () => void) => {
		const items = [
			{
				label: (
					<MenuLink to="/portal/settings" onClick={close} icon="user-cog">
						Settings
					</MenuLink>
				),
			},
			{
				label: (
					<MenuLink onClick={() => logout()} icon="sign-out">
						Log Out
					</MenuLink>
				),
			},
		]

		return items
	}

	return (
		<StatefulPopover
			onOpen={() => setIsOpen(true)}
			onClose={() => setIsOpen(false)}
			placement={PLACEMENT.bottomRight}
			content={({ close }) => (
				<div className={back}>
					<MenuItemText>
						Signed in as <strong>{`${user.firstName} ${user.lastName}`}</strong>
					</MenuItemText>
					<StatefulMenu
						items={menuItems(close)}
						overrides={{
							List: {
								style: {
									outline: "unset",
									boxShadow: "unset",
								},
							},
							Option: {
								style: ({ $theme }) => {
									return {
										...paddingZero,
										color: "black",
										":hover": {
											backgroundColor: $theme.colors.primary,
											color: "white",
										},
										transitionTimingFunction: "unset",
										transitionDuration: "unset",
										transitionProperty: "unset",
									}
								},
							},
						}}
					/>
				</div>
			)}
			accessibilityType={"tooltip"}
			overrides={{
				Body: { style: { boxShadow: boxShadowCommon } },
			}}
		>
			<div className={userStyle}>
				<div className={userCircleStyle}>
					<FontAwesomeIcon icon={["fal", "user-circle"]} />
				</div>
				<FontAwesomeIcon icon={["fal", isOpen ? "chevron-up" : "chevron-down"]} size={"xs"} />
			</div>
		</StatefulPopover>
	)
}

const MenuItemText = styled("div", ({ $theme }) => ({
	...$theme.typography.font200,
	fontSize: "16px",
	lineHeight: "24px",
	color: "black",
	padding: "8px 16px",
	borderBottom: "solid rgba(27, 31, 35, 0.15) 1px",
}))

const MenuLink: React.FunctionComponent<{ to?: string; onClick?: () => void; icon?: IconName }> = props => {
	const history = useHistory()

	return (
		<>
			<Button
				onClick={() => {
					if (props.onClick) props.onClick()
					if (props.to) history.push(props.to)
				}}
				overrides={{
					BaseButton: {
						style: {
							paddingTop: "8px",
							paddingBottom: "8px",
							paddingLeft: "16px",
							paddingRight: "16px",
							width: "100%",
							justifyContent: "left",
							backgroundColor: "unset",
							fontSize: "16px",
							lineHeight: "24px",
							color: "unset",
							":hover": {
								backgroundColor: "unset",
							},
							":active": {
								backgroundColor: "unset",
							},
							":focus": {
								backgroundColor: "unset",
							},
						},
					},
				}}
			>
				<Block display="inline-flex">
					{props.icon && <FontAwesomeIcon icon={["fal", props.icon]} size="lg" />}
					<Block marginLeft="8px">{props.children}</Block>
				</Block>
			</Button>
		</>
	)
}
