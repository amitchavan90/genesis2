import * as React from "react"
import { Button } from "baseui/button"
import { paddingZero } from "../themeOverrides"

interface StyledTextButtonProps {
	onClick?: () => void
	disabled?: boolean
	underline?: boolean
	color?: string
}
export const StyledTextButton: React.FunctionComponent<StyledTextButtonProps> = props => {
	return (
		<Button
			type="button"
			onClick={props.onClick}
			disabled={props.disabled}
			overrides={{
				BaseButton: {
					style: ({ $theme }) => ({
						color: props.color || $theme.colors.colorPrimary,
						fontWeight: 500,
						textDecoration: props.underline ? "underline" : "unset",
						whiteSpace: "nowrap",
						...paddingZero,

						backgroundColor: "unset",
						":hover": {
							backgroundColor: "unset",
							color: $theme.colors.accent,
						},
						":active": {
							backgroundColor: "unset",
						},
						":focus": {
							backgroundColor: "unset",
						},
					}),
				},
			}}
		>
			{props.children}
		</Button>
	)
}
