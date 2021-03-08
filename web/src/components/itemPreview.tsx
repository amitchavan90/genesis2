import * as React from "react"
import { useStyletron } from "baseui"
import { Blob, SKU } from "../types/types"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { IconName } from "@fortawesome/pro-light-svg-icons"
import { Button } from "baseui/button"

interface ItemPreviewProps {
	name: string
	code: string
	description: string
	thumbnail?: Blob
	icon: IconName
	changable?: Boolean
	changeOnClick?: (evt: React.MouseEvent<HTMLButtonElement, MouseEvent>) => void
}
export const ItemPreview = (props: ItemPreviewProps) => {
	const [css, theme] = useStyletron()
	const container = css({
		...theme.typography.font200,
		display: "flex",
		paddingLeft: "2px",
		paddingRight: "2px",
		paddingTop: "2px",
		paddingBottom: "2px",
		color: "inherit",
		marginRight: "0.5rem",
	})
	const infoContainer = css({
		alignSelf: "flex-start",
		marginTop: "10px",
		marginBottom: "10px",
		marginLeft: "15px",
		textAlign: "left",
		whiteSpace: "normal",
	})
	const nameStyle = css({
		fontSize: "1.15rem",
		lineHeight: "1rem",
	})
	const nameSubStyle = css({
		fontSize: "0.8rem",
		lineHeight: "1rem",
		color: "grey",
		marginTop: "3px",
	})
	const descStyle = css({
		fontSize: "0.8rem",
		color: "grey",
		marginTop: "6px",
	})
	const imageStyle = css({
		width: "120px",
		minWidth: "120px",
		objectFit: "contain",
		textAlign: "center",
		margin: "auto 0",
	})

	return (
		<div>
			<div className={container}>
				{props.thumbnail && <img className={imageStyle} src={props.thumbnail.file_url} />}
				{!props.thumbnail && (
					<div className={imageStyle}>
						<FontAwesomeIcon icon={["fal", props.icon]} size="4x" />
					</div>
				)}
				<div className={infoContainer}>
					<div className={nameStyle}>{props.name}</div>
					<div className={nameSubStyle}>{props.code}</div>
					<div className={descStyle}>{props.description}</div>
				</div>
			</div>
			{props.changable && (
				<Button
					onClick={props.changeOnClick}
					overrides={{
						BaseButton: {
							style: {
								marginTop: "10px",
								width: "120px",
							},
						},
					}}
				>
					Change
				</Button>
			)}
		</div>
	)
}

export const SKUItemPreview = (props: { sku: SKU }) => (
	<ItemPreview name={props.sku.name} code={props.sku.code} description={props.sku.description} thumbnail={props.sku.masterPlan} icon="barcode-alt" />
)
