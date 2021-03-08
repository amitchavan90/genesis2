import * as React from "react"
import { Button } from "baseui/button"
import { StatefulPopover } from "baseui/popover"
import { StatefulMenu } from "baseui/menu"
import { MenuOptionLabel, MenuListOverride, MenuOptionStyleOverride } from "./menuList"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { useMutation } from "@apollo/react-hooks"
import { ItemSelectList } from "./itemSelectList"
import { Modal, ModalHeader, ModalBody, ModalFooter } from "baseui/modal"
import { ErrorNotification } from "./errorBox"
import { Spinner } from "baseui/spinner"
import { useStyletron } from "baseui"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { Value } from "baseui/select"
import { graphql } from "../graphql"
import { GetItemIcon } from "../themeOverrides"
import { Action, ActionItem } from "../types/actions"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Datepicker } from "baseui/datepicker"
import { TimePicker } from "baseui/timepicker"
import { Checkbox } from "baseui/checkbox"

interface ActionsButtonProps {
	options: ActionItem[]
	items: any[]
	batchActionMutation?: any
	refetch?: () => void
	disabled?: boolean
}

interface BatchActionInput {
	str?: string
	no?: number
	dateTime?: Date
	bool?: boolean
}

export const ActionsButton = (props: ActionsButtonProps) => {
	const { options, items, batchActionMutation, refetch, disabled } = props

	// Batch action mutation
	const [batchAction, mutBatchAction] = useMutation(batchActionMutation)
	const [showModal, setShowModal] = React.useState<boolean>()
	const [modalAction, setModalAction] = React.useState<number>()
	const [selectValue, setSelectValue] = React.useState<Value>()
	const [numberValue, setNumberValue] = React.useState<number>()
	const [dateTimeValue, setDateTimeValue] = React.useState<Date>()
	const [boolValue, seBoolValue] = React.useState<boolean>(false)

	const onAction = (action: Action, value?: BatchActionInput) => {
		if (!batchActionMutation) return

		let ids: string[] = []
		items.forEach(i => i.selected && ids.push(i.id))

		batchAction({ variables: { ids: ids, action, value } })
	}

	React.useEffect(() => {
		if (mutBatchAction.loading || !mutBatchAction.data) return
		if (refetch) refetch()
	}, [mutBatchAction])

	const [css] = useStyletron()
	const flexStyle = css({
		display: "flex",
	})
	const iconStyle = css({
		marginRight: "10px",
		width: "20px",
		height: "20px",
		textAlign: "center",
	})

	const getListQuery = (action: Action) => {
		switch (action) {
			case Action.SetSKU:
				return graphql.query.SKUS
			case Action.SetOrder:
				return graphql.query.ORDERS_BASIC
			case Action.SetDistributor:
				return graphql.query.DISTRIBUTORS_BASIC
			case Action.SetCarton:
				return graphql.query.CARTONS_BASIC
			case Action.SetPallet:
				return graphql.query.PALLETS_BASIC
			case Action.SetContainer:
				return graphql.query.CONTAINERS_BASIC
			case Action.SetContract:
				return graphql.query.CONTRACTS_BASIC
		}
		return graphql.query.CARTONS_BASIC
	}

	return (
		<>
			<StatefulPopover
				showArrow
				placement="bottomRight"
				content={({ close }) => (
					<StatefulMenu
						items={options.map((option, index) => {
							const icon = option.itemName ? GetItemIcon(option.itemName) : { icon: option.icon, light: option.iconLight }
							return {
								label: (
									<MenuOptionLabel
										label={option.label}
										icon={icon.icon}
										iconLight={icon.light}
										borderTop={index != 0 && (option.itemName != options[index - 1].itemName || options[index - 1].action == Action.Unarchive)}
									/>
								),
								id: option.action,
								index: index,
							}
						})}
						onItemSelect={({ item }) => {
							if (item?.id) {
								if ((item.id as Action).toString().startsWith("Set") || (item.id as Action) == Action.InheritCartonHistory) {
									// Open Set Modal
									setModalAction(item.index)
									setShowModal(true)
								} else {
									// Do Action
									onAction(item.id)
								}
							}
							close()
						}}
						overrides={{
							List: MenuListOverride,
							Option: { style: MenuOptionStyleOverride },
						}}
					/>
				)}
				accessibilityType="menu"
				overrides={{
					Body: { style: { boxShadow: "rgba(0, 0, 0, 0.13) 0px 4px 4px" } },
				}}
			>
				<Button disabled={disabled} isLoading={mutBatchAction.loading}>
					Actions
				</Button>
			</StatefulPopover>

			<Modal
				isOpen={showModal}
				onClose={() => setShowModal(false)}
				overrides={{
					Dialog: {
						style: {
							width: modalAction ? "70%" : "40%",
						},
					},
				}}
			>
				{modalAction && (
					<>
						<ModalHeader>
							<div className={flexStyle}>
								{options[modalAction].icon && (
									<div className={iconStyle}>
										<FontAwesomeIcon icon={[options[modalAction].iconLight ? "fal" : "fas", options[modalAction].icon as IconName]} />
									</div>
								)}
								<div>{options[modalAction].label}</div>
							</div>
						</ModalHeader>
						<ModalBody>
							{mutBatchAction.loading && <Spinner />}

							{mutBatchAction.error && <ErrorNotification message={mutBatchAction.error.message} />}

							{!mutBatchAction.loading && !mutBatchAction.error && (
								<>
									{options[modalAction].action == Action.SetBonusLoyaltyPoints ? (
										<SetBonusLoyaltyPointsComponent
											numberValue={numberValue || 0}
											setNumberValue={setNumberValue}
											dateTimeValue={dateTimeValue || new Date()}
											setDateTimeValue={setDateTimeValue}
										/>
									) : (
										<ItemSelectList
											useTable={
												options[modalAction].action === Action.SetCarton ||
												options[modalAction].action === Action.SetContainer ||
												options[modalAction].action === Action.SetPallet ||
												options[modalAction].action === Action.InheritCartonHistory ||
												options[modalAction].action === Action.SetDistributor ||
												options[modalAction].action === Action.SetSKU ||
												options[modalAction].action === Action.SetOrder ||
												options[modalAction].action === Action.SetContract
											}
											itemName={options[modalAction].itemName || ""}
											value={selectValue}
											setValue={setSelectValue}
											identifier={options[modalAction].identifier}
											query={getListQuery(options[modalAction].action)}
										/>
									)}

									{options[modalAction].action == Action.SetCarton && (
										<Checkbox checked={boolValue} onChange={e => seBoolValue(e.currentTarget.checked)}>
											Inherit Carton Tracking History
										</Checkbox>
									)}
								</>
							)}
						</ModalBody>
						<ModalFooter>
							<Button
								onClick={() => {
									onAction(options[modalAction].action, {
										str: selectValue && selectValue.length > 0 ? (selectValue[0].id as string) : "-",
										no: numberValue,
										dateTime: dateTimeValue,
										bool: boolValue,
									})
									setSelectValue(undefined)
									setShowModal(false)
								}}
							>
								Confirm
							</Button>
						</ModalFooter>
					</>
				)}
			</Modal>
		</>
	)
}

interface SetBonusLoyaltyPointsComponentProps {
	numberValue: number
	setNumberValue: (value: number) => void
	dateTimeValue: Date
	setDateTimeValue: (value: Date) => void
}
const SetBonusLoyaltyPointsComponent = (props: SetBonusLoyaltyPointsComponentProps) => {
	const { numberValue, setNumberValue, dateTimeValue, setDateTimeValue } = props

	return (
		<>
			<FormControl label="Bonus Loyalty Points" error="" positive="">
				<Input type="number" value={numberValue.toString()} onChange={e => setNumberValue(+e.currentTarget.value)} />
			</FormControl>

			<FormControl
				label="Bonus Loyalty Points Expiration"
				caption="After the expiration date, no bonus points will be earned from these products."
				error=""
				positive=""
			>
				<div style={{ display: "flex" }}>
					<div
						style={{
							width: "160px",
							marginRight: "10px",
						}}
					>
						<Datepicker value={dateTimeValue} onChange={({ date }) => setDateTimeValue(date as Date)} />
					</div>
					<div style={{ width: "120px" }}>
						<TimePicker value={dateTimeValue} onChange={setDateTimeValue} />
					</div>
				</div>
			</FormControl>
		</>
	)
}
