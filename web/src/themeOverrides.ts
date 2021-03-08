import { createTheme, lightThemePrimitives, darkThemePrimitives } from "baseui"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import "./styles.scss"
import { ButtonOverrides } from "baseui/button"

// See https://github.com/uber-web/baseui/blob/master/src/themes/creator.js for full list of theme properties

// export const primaryOrange = "#FE7B6B"
// export const secondaryPurple = "#6C44A0"
// export const darkPurple = "#3E046A"

export const boxShadowCommon = "rgba(0, 0, 0, 0.13) 0px 4px 4px"
export const paddingZero = { paddingLeft: "0", paddingRight: "0", paddingTop: "0", paddingBottom: "0" }

export const LightBlue = "#59B7E8"

export const LightTheme = createTheme(
	{
		...lightThemePrimitives,
		primaryFontFamily: "Konnect, -apple-system, BlinkMacSystemFont, San Francisco, Roboto, Segoe UI, Helvetica Neue, sans-serif",
	},
	{
		// add all the theme overrides here - under the hood it uses deep merge
		// animation: {
		// 	timing100: '0.50s',
		// },
		colors: {
			primary: "#1a1a1a",
			colorPrimary: "#1a1a1a",
			// progressStepsCompletedFill: primaryOrange,
			// buttonPrimaryFill: primaryOrange,
			// buttonPrimaryHover: "#FF9585",
			// buttonSecondaryFill: secondaryPurple,
			// buttonSecondaryHover: "#865EBA",
			// buttonSecondaryText: "#FFFFFF",
		},
	},
)

export const ButtonMarginLeftOverride: ButtonOverrides = {
	BaseButton: {
		style: {
			marginLeft: "0.5rem",
		},
	},
}

export const GetItemIcon = (name: string, small?: boolean) => {
	let icon = "steak"
	let light = !small
	switch (name) {
		case "order":
			icon = "shopping-cart"
			light = false
			break
		case "distributor":
			icon = "shopping-basket"
			break
		case "carton":
			icon = "box"
			break
		case "pallet":
			icon = "pallet-alt"
			break
		case "container":
			icon = "container-storage"
			light = true
			break
		case "sku":
			icon = "barcode-alt"
			light = true
			break
		case "user":
			icon = "user"
			light = false
			break
		case "consumer":
			icon = "smile"
			break
		case "role":
			icon = "user-tag"
			light = false
			break
		case "registered":
			icon = "smile"
			light = true
			break
		case "trackAction":
			icon = "truck-moving"
			light = false
			break
		case "contract":
			icon = "file-contract"
			break
		case "transaction":
			icon = "cubes"
			light = false
			break
		case "userActivity":
			icon = "user-chart"
			light = false
			break
		case "blob":
			icon = "file"
			break
	}
	return { icon: icon as IconName, light: light }
}
