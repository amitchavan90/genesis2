import { Value, Select } from "baseui/select"
import React from "react"
import { Modal, ModalHeader, ModalBody, ModalFooter } from "baseui/modal"
import { ItemSelectList } from "./itemSelectList"
import { Button } from "baseui/button"

interface TablePickerSelectProps {
	/* modal */
	isOpen: boolean
	setIsOpen: (b: boolean) => void
	/* select value */
	setValue: (v: Value) => void
	value?: Value
	/* list query eg. graphql.query.PALLETS_BASIC */
	query: any
	hash?: string
	itemName: string
	identifier?: string
	/** The column to use for the label */
	labelResolver?: (value: any) => void
	modalTitle?: string
	queryName: string
}
export const TablePickerSelect = (props: TablePickerSelectProps) => {
	/* modal title with first letter uppercased */
	const title = (props.modalTitle || props.itemName).charAt(0).toUpperCase() + (props.modalTitle || props.itemName).slice(1)

	return (
		<div>
			<Select
				labelKey="label"
				valueKey="id"
				autoFocus={false}
				openOnClick
				value={props.value}
				onChange={({ value }) => {
					props.setValue(value)
				}}
				onOpen={() => props.setIsOpen(true)}
				clearable
			/>
			<Modal
				overrides={{
					Dialog: {
						style: {
							width: "70%",
						},
					},
				}}
				onClose={() => props.setIsOpen(false)}
				isOpen={props.isOpen}
			>
				<ModalHeader>{`Select ${title}`} </ModalHeader>
				<ModalBody>
					<ItemSelectList
						queryName={props.queryName}
						identifier={props.identifier}
						useTable
						itemName={props.itemName}
						value={props.value}
						setValue={props.setValue}
						labelResolver={props.labelResolver}
						query={props.query}
						hash={props.hash}
					/>
				</ModalBody>
				<ModalFooter>
					<Button
						onClick={() => {
							props.setIsOpen(false)
						}}
						overrides={{
							BaseButton: {
								style: {
									display: "inline-block",
								},
							},
						}}
					>
						OK
					</Button>
				</ModalFooter>
			</Modal>
		</div>
	)
}
