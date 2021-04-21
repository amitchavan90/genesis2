import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { Order } from "../../../types/types"
import { useForm } from "react-hook-form"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { H1, H3 } from "baseui/typography"
import { Button } from "baseui/button"
import { Notification } from "baseui/notification"
import { Spaced } from "../../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { LoadingSimple } from "../../../components/loading"
import { ErrorNotification } from "../../../components/errorBox"
import { Spread } from "../../../components/spread"
import { Tabs, Tab } from "baseui/tabs"
import { paddingZero } from "../../../themeOverrides"
import { CenteredPage } from "../../../components/common"
import { ItemList } from "../../../components/itemList"
import { FilterOption } from "../../../types/enums"
import { ActionItemSet } from "../../../types/actions"
import { Value } from "baseui/select"
import { SKUSelectList, ItemSelectList } from "../../../components/itemSelectList"
import { SKUItemPreview } from "../../../components/itemPreview"
import { invalidateListQueries } from "../../../apollo"
import { TablePickerSelect } from "../../../components/tablePickerSelect"
import { Checkbox } from "baseui/checkbox"

type FormData = {
	quantity: number
}

const OrderEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewOrder = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	// Get Order
	const [order, setOrder] = React.useState<Order>()
	const { data, loading, error, refetch } = useQuery<{ order: Order }>(graphql.query.ORDER, {
		variables: { code },
		fetchPolicy: isNewOrder ? "cache-only" : undefined, // prevent query if new
	})

	const [createdOrder, setCreatedOrder] = React.useState<Order>()

	// Mutations
	const [updateOrder, mutUpdateOrder] = useMutation(isNewOrder ? graphql.mutation.CREATE_ORDER : graphql.mutation.UPDATE_ORDER)
	const [archiveOrder, mutArchiveOrder] = useMutation<{ orderArchive: Order }>(graphql.mutation.ARCHIVE_ORDER)
	const [unarchiveOrder, mutUnarchiveOrder] = useMutation<{ orderUnarchive: Order }>(graphql.mutation.UNARCHIVE_ORDER)

	// modal
	const [contractModalOpen, setContractModalOpen] = React.useState(false)

	const toggleArchive = () => {
		if (!order) return

		if (order.archived) {
			unarchiveOrder({
				variables: { id: order.id },
				update: (cache: any) => invalidateListQueries(cache, "orders"),
			})
		} else {
			archiveOrder({
				variables: { id: order.id },
				update: (cache: any) => invalidateListQueries(cache, "orders"),
			})
		}
	}

	// Form submission
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [contract, setContract] = React.useState<Value>()
	const [sku, setSKU] = React.useState<Value>()
	const [contractError, setContractError] = React.useState<string>()
	const [isAppBound, setIsAppBound] = React.useState<boolean>(true)

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ quantity }) => {
		setChangeSuccess(false)

		if (isNewOrder) {
			// New order
			updateOrder({
				variables: {
					input: {
						contractID: contract && contract.length > 0 ? contract[0].id : undefined,
						skuID: sku && sku.length > 0 ? sku[0].id : undefined,
						quantity,
						isAppBound
					},
				},
				update: (cache: any) => {
					invalidateListQueries(cache, "orders")
					invalidateListQueries(cache, "products")
				},
			})
			return
		}

		// Update order
		if (!order) return
		updateOrder({
			variables: {
				id: order.id,
				input: { code },
			},
		})
	})

	// On load order
	React.useEffect(() => {
		if (!data || !data.order) return
		setOrder(data.order)
	}, [data, loading, error])
	React.useEffect(() => {
		if (activeKey != "#details") return
		if (!order) return
		setValue("code", order.code)
		setIsAppBound(order.isAppBound)
	}, [order, activeKey])

	// On mutation (update/create order)
	React.useEffect(() => {
		if (!mutUpdateOrder.data) return

		if (isNewOrder) {
			if (!mutUpdateOrder.data.orderCreate) return

			setChangeSuccess(true)
			setCreatedOrder(mutUpdateOrder.data.orderCreate)
			return
		}

		if (!mutUpdateOrder.data.orderUpdate) return

		setOrder(mutUpdateOrder.data.orderUpdate)
		setChangeSuccess(true)
	}, [mutUpdateOrder])

	React.useEffect(() => {
		if (!mutArchiveOrder.data?.orderArchive) return
		setOrder(mutArchiveOrder.data.orderArchive)
	}, [mutArchiveOrder])
	React.useEffect(() => {
		if (!mutUnarchiveOrder.data?.orderUnarchive) return
		setOrder(mutUnarchiveOrder.data.orderUnarchive)
	}, [mutUnarchiveOrder])

	// Styling
	const [css, theme] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "10px",
	})
	const createSuccessStyle = css({
		textAlign: "center",
		marginTop: "30px",
	})
	const breakLineStyle = css({
		border: "1px solid rgba(0, 0, 0, 0.1)",
		marginBottom: "10px",
	})
	const skuPreview = css({
		border: "2px solid rgba(0, 0, 0, 0.25)",
		margin: "10px",
		padding: "10px",
		":hover": {
			backgroundColor: theme.colors.colorPrimary,
			color: "white",
		},
		cursor: "pointer",
	})
	const breakLine = <div className={breakLineStyle} />

	if (!isNewOrder && !order) {
		return <LoadingSimple />
	}

	if (isNewOrder && changeSuccess && createdOrder) {
		// Successfully created order
		return (
			<CenteredPage>
				<div className={createSuccessStyle}>
					<FontAwesomeIcon icon={["far", "check-circle"]} size="10x" color="#1db954" />
					<H3>Order created</H3>

					<Spread>
						<Button
							type="button"
							kind="secondary"
							onClick={() => {
								setChangeSuccess(false)
								history.push(`/portal/order/${createdOrder.code}#products`)
							}}
						>
							<FontAwesomeIcon icon={["fas", "shopping-cart"]} />
							<div className={marginLeftStyle}>{createdOrder.code}</div>
						</Button>
						<Button kind="secondary" type="button" onClick={() => window.open(`/api/files/order?id=${createdOrder.id}`)}>
							<FontAwesomeIcon icon={["fas", "file-excel"]} color="#1db954" />
							<div className={marginLeftStyle}>Download Spreadsheet</div>
						</Button>
					</Spread>
				</div>
			</CenteredPage>
		)
	}

	const editForm = (
		<>
			{changeSuccess && (
				<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
					The order has been updated.
				</Notification>
			)}
			{mutUpdateOrder.error && <ErrorNotification message={mutUpdateOrder.error.message} />}
			<form onSubmit={onSubmit}>
				{!isNewOrder && order && (
					<>
						{order.sku && (
							<FormControl label="SKU" disabled error="" positive="">
								<div className={skuPreview} onClick={() => history.push(`/portal/sku/${order.sku.code}`)}>
									<SKUItemPreview sku={order.sku} />
								</div>
							</FormControl>
						)}
					</>
				)}

				{isNewOrder && (
					<>
						<FormControl label="SKU" error="">
							<SKUSelectList
								value={sku}
								setValue={value => {
									setSKU(value)
								}}
							/>
						</FormControl>

						<FormControl caption="">
							<Checkbox checked={isAppBound} onChange={e => setIsAppBound(e.currentTarget.checked)}>
								Is App
							</Checkbox>
						</FormControl>

						<FormControl label="Livestock Specification" error={contractError || ""}>
							<TablePickerSelect
								queryName="contracts"
								isOpen={contractModalOpen}
								setIsOpen={setContractModalOpen}
								hash=""
								itemName="contract"
								value={contract}
								setValue={setContract}
								query={graphql.query.CONTRACTS_BASIC}
								modalTitle="Livestock Specification"
							/>
						</FormControl>

						<FormControl label="Quantity" error={errors.quantity ? errors.quantity.message : ""} positive="">
							<Input
								name="quantity"
								type="number"
								inputRef={register({
									required: "Required",
									min: { value: 1, message: "Quantity must be greater than 0" },
									max: { value: 10000, message: "Quantity cannot be greater than 10000" },
								})}
								error={errors.quantity !== undefined}
							/>
						</FormControl>
					</>
				)}

				{breakLine}

				<Spread>
					<Button type="button" kind="secondary" onClick={() => history.push("/portal/orders")}>
						Cancel
					</Button>
					{order && !isNewOrder ? (
						<Button
							type="button"
							kind="secondary"
							isLoading={mutArchiveOrder.loading || mutUnarchiveOrder.loading}
							onClick={toggleArchive}
							startEnhancer={<FontAwesomeIcon icon={["fas", order.archived ? "undo" : "archive"]} size="lg" />}
						>
							{order.archived ? "Unarchive" : "Archive"}
						</Button>
					) : (
						<Button isLoading={mutUpdateOrder.loading}>Create Order</Button>
					)}
				</Spread>
			</form>
		</>
	)

	return (
		<CenteredPage>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={["fas", "shopping-cart"]} size="3x" />
					<H1>{isNewOrder ? "New Order" : code} </H1>
				</Spaced>
				<div>
					{order && (
						<Spaced>
							<Button kind="secondary" type="button" onClick={() => window.open(`/api/files/order?id=${order.id}`)}>
								<FontAwesomeIcon icon={["fas", "file-excel"]} color="#1db954" />
								<div className={marginLeftStyle}>Download</div>
							</Button>

							{order.sku && (
								<Button kind="secondary" type="button" onClick={() => history.push(`/portal/sku/${order.sku.code}`)}>
									<FontAwesomeIcon icon={["fal", "barcode-alt"]} />
									<div className={marginLeftStyle}>{order.sku.code}</div>
								</Button>
							)}
						</Spaced>
					)}
				</div>
			</Spread>

			{isNewOrder ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/order/${code}${activeKey}`)
					}}
					activeKey={activeKey}
					overrides={{
						TabContent: {
							style: { ...paddingZero },
						},
						TabBar: {
							style: { ...paddingZero },
						},
					}}
				>
					<Tab
						key="#details"
						title={
							<Spaced>
								<FontAwesomeIcon icon={["fal", "pencil"]} />
								<div>Details</div>
							</Spaced>
						}
					>
						{editForm}
					</Tab>

					{!isNewOrder && order && (
						<Tab
							key="#products"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "steak"]} />
									<div>Products</div>
								</Spaced>
							}
						>
							<ItemList
								orderID={order.id}
								itemName="product"
								query={graphql.query.PRODUCTS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_PRODUCT}
								extraFilterOptions={[
									{ label: "Not in Carton", id: FilterOption.ProductWithoutCarton },
									{ label: "No SKU", id: FilterOption.ProductWithoutSKU },
								]}
								itemLinks={["sku", "contract", "distributor", "carton"]}
								actions={ActionItemSet.Products}
								onListUpdate={refetch}
								showQRCodesToggle
							/>
						</Tab>
					)}
				</Tabs>
			)}
		</CenteredPage>
	)
}

export default OrderEdit
