import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation, useLazyQuery } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { Product, Transaction, Settings } from "../../../types/types"
import { useForm } from "react-hook-form"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { H1 } from "baseui/typography"
import { Button } from "baseui/button"
import { Notification } from "baseui/notification"
import { Spaced } from "../../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { LoadingSimple } from "../../../components/loading"
import { ErrorNotification } from "../../../components/errorBox"
import { Spread } from "../../../components/spread"
import { Value } from "baseui/select"
import { ItemSelectList, SKUSelectList } from "../../../components/itemSelectList"
import { CenteredPage } from "../../../components/common"
import { Datepicker } from "baseui/datepicker"
import { TimePicker } from "baseui/timepicker"
import { Tabs, Tab } from "baseui/tabs"
import { paddingZero, ButtonMarginLeftOverride } from "../../../themeOverrides"
import { QRCode } from "react-qrcode-logo"
import { TransactionsView } from "../../../components/transactionsView"
import { StyledLink } from "baseui/link"
import { invalidateListQueries } from "../../../apollo"
import { Checkbox } from "baseui/checkbox"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"
import { TablePickerSelect } from "../../../components/tablePickerSelect"

type FormData = {
	loyaltyPoints: number
	description: string
}
const ProductEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewProduct = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	const [showDemoQRCode, setShowDemoQRCode] = React.useState(Boolean)

	// Get Product
	const [product, setProduct] = React.useState<Product>()
	const { data, loading, error } = useQuery<{ product: Product }>(graphql.query.PRODUCT, {
		variables: { code },
		fetchPolicy: isNewProduct ? "cache-only" : undefined, // prevent query if new
	})

	// Get Transactions
	const [getTransactions, queryTransactions] = useLazyQuery<{ productByID: Product }>(graphql.query.PRODUCT_TRANSACTIONS_VIEW)

	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)

	// Mutations
	const [updateProduct, mutUpdateProduct] = useMutation(isNewProduct ? graphql.mutation.CREATE_PRODUCT : graphql.mutation.UPDATE_PRODUCT)
	const [archiveProduct, mutArchiveProduct] = useMutation<{ productArchive: Product }>(graphql.mutation.ARCHIVE_PRODUCT)
	const [unarchiveProduct, mutUnarchiveProduct] = useMutation<{ productUnarchive: Product }>(graphql.mutation.UNARCHIVE_PRODUCT)

	// modal states for tableSelectPickers
	const [distributorModalOpen, setDistributorModalOpen] = React.useState(false)
	const [contractModalOpen, setContractModalOpen] = React.useState(false)
	const [orderModalOpen, setOrderModalOpen] = React.useState(false)
	const [cartonModalOpen, setCartonModalOpen] = React.useState(false)

	const toggleArchive = () => {
		if (!product) return

		if (product.archived) {
			unarchiveProduct({
				variables: { id: product.id },
				update: (cache: any) => invalidateListQueries(cache, "products"),
			})
		} else {
			archiveProduct({
				variables: { id: product.id },
				update: (cache: any) => invalidateListQueries(cache, "products"),
			})
		}
	}

	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [carton, setCarton] = React.useState<Value>()
	const [inheritCartonHistory, setInheritCartonHistory] = React.useState(false)
	const [order, setOrder] = React.useState<Value>()
	const [distributor, setDistributor] = React.useState<Value>()
	const [sku, setSKU] = React.useState<Value>()
	const [loyaltyPointsExpire, setLoyaltyPointsExpire] = React.useState(new Date())
	const [contract, setContract] = React.useState<Value>()

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ loyaltyPoints, description }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		const input = {
			code,
			cartonID: carton && carton.length > 0 ? carton[0].id : "-",
			orderID: order && order.length > 0 ? order[0].id : "-",
			distributorID: distributor && distributor.length > 0 ? distributor[0].id : "-",
			skuID: sku && sku.length > 0 ? sku[0].id : "-",
			contractID: contract && contract.length > 0 ? contract[0].id : "-",
			loyaltyPoints: loyaltyPoints || 0,
			loyaltyPointsExpire,
			inheritCartonHistory,
			description,
		}

		if (isNewProduct)
			updateProduct({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "products"),
			})
		else if (product) {
			promiseTimeout(updateProduct({ variables: { id: product.id, input } })).catch((reason) => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	// On load product
	React.useEffect(() => {
		if (!data || !data.product) return
		setProduct(data.product)
	}, [data, loading, error])
	React.useEffect(() => {
		// Load contract on tracking tab visit
		if (!queryTransactions.data && product && activeKey == "#tracking") getTransactions({ variables: { id: product.id } })

		// Set product values when on edit tab
		if (activeKey != "#details") return
		if (!product) return
		setValue("code", product.code)
		setValue("loyaltyPoints", product.loyaltyPoints)
		setLoyaltyPointsExpire(new Date(product.loyaltyPointsExpire))
		setValue("description", product.description)
		if (product.carton) setCarton([{ id: product.carton.id, label: product.carton.code }])
		if (product.order) setOrder([{ id: product.order.id, label: product.order.code }])
		if (product.distributor) setDistributor([{ id: product.distributor.id, label: product.distributor.code }])
		if (product.sku) setSKU([{ id: product.sku.id, label: product.sku.code }])
		if (product.contract) setContract([{ id: product.contract.id, label: `${product.contract.name} (${product.contract.supplierName})` }])
	}, [product, activeKey])

	// On mutation (update/create product)
	React.useEffect(() => {
		if (!mutUpdateProduct.data) return

		if (isNewProduct) {
			if (mutUpdateProduct.data.productCreate) {
				history.push(`/portal/product/${mutUpdateProduct.data.productCreate.code}`)
			}
			return
		}

		if (!mutUpdateProduct.data.productUpdate) return

		setProduct(mutUpdateProduct.data.productUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateProduct])

	React.useEffect(() => {
		if (!mutArchiveProduct.data?.productArchive) return
		setProduct(mutArchiveProduct.data.productArchive)
	}, [mutArchiveProduct])
	React.useEffect(() => {
		if (!mutUnarchiveProduct.data?.productUnarchive) return
		setProduct(mutUnarchiveProduct.data.productUnarchive)
	}, [mutUnarchiveProduct])

	// Styling
	const [css, theme] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "10px",
	})
	const flexStyle = css({
		display: "flex",
	})
	const flexColumnStyle = css({
		display: "flex",
		flexDirection: "column",
	})
	const viewLinkContainerStyle = css({
		display: "flex",
		alignSelf: "flex-end",
	})
	const breadCrumbContainerStyle = css({
		display: "flex",
		paddingTop: "10px",
		alignSelf: "flex-end",
	})
	const breadCrumbStyle = css({
		display: "flex",
		alignItems: "center",
	})
	const breadCrumbChevronStyle = css({
		height: "24px",
		paddingLeft: "10px",
		paddingRight: "10px",
	})
	const breadCrumbChevron = <FontAwesomeIcon icon={["fal", "chevron-right"]} className={breadCrumbChevronStyle} />

	const breakLineStyle = css({
		border: "1px solid rgba(0, 0, 0, 0.1)",
		marginBottom: "10px",
	})
	const breakLine = <div className={breakLineStyle} />

	const qrItemsContainer = css({
		display: "flex",
		justifyContent: "space-evenly",
		paddingTop: "5px",
	})
	const qrItem = css({
		display: "flex",
		textAlign: "center",
		justifyContent: "center",
		flexDirection: "column",
		alignItems: "center",
		padding: "10px",
	})
	const qrLink = css({
		fontSize: "9px",
		lineHeight: "12px",
		color: "grey",
	})

	if (!isNewProduct && !product) {
		return <LoadingSimple />
	}

	if (settingsQuery.loading) {
		return <LoadingSimple />
	}

	if (settingsQuery.error) {
		return <p>Error: {settingsQuery.error.message}</p>
	}
	if (!settingsQuery.data) {
		return <p>No settings returned</p>
	}
	const viewLink = product ? `${settingsQuery.data.settings.consumerHost}/api/steak/view?productID=${product.id}` : ""
	const registerLink = product ? `${settingsQuery.data.settings.consumerHost}/api/steak/detail?steakID=${product.registerID}&productID=${product.id}` : ""
	const registerLinkDemo = product ? `${settingsQuery.data.settings.consumerHost}/api/steak/final?steakID=${product.registerID}&demo=true` : ""

	const editForm = (
		<form onSubmit={onSubmit}>
			{changeSuccess && (
				<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
					The product has been updated.
				</Notification>
			)}

			{mutUpdateProduct.error && <ErrorNotification message={mutUpdateProduct.error.message} />}

			<FormControl label="SKU">
				<SKUSelectList value={sku} setValue={setSKU} />
			</FormControl>

			<FormControl label="Distributor">
				<TablePickerSelect
					queryName="distributors"
					isOpen={distributorModalOpen}
					setIsOpen={setDistributorModalOpen}
					itemName="distributor"
					value={distributor}
					setValue={setDistributor}
					query={graphql.query.DISTRIBUTORS_BASIC}
				/>
			</FormControl>

			<FormControl label="Livestock Specification">
				<TablePickerSelect
					queryName="contracts"
					isOpen={contractModalOpen}
					setIsOpen={setContractModalOpen}
					hash=""
					itemName="contract"
					value={contract}
					setValue={setContract}
					query={graphql.query.CONTRACTS_BASIC}
					identifier="name"
					modalTitle="Livestock Specification"
				/>
			</FormControl>

			<FormControl label="Order">
				<TablePickerSelect
					queryName="orders"
					isOpen={orderModalOpen}
					setIsOpen={setOrderModalOpen}
					itemName="order"
					value={order}
					setValue={setOrder}
					query={graphql.query.ORDERS_BASIC}
				/>
			</FormControl>

			<FormControl label="Carton">
				<>
					<TablePickerSelect
						queryName="cartons"
						isOpen={cartonModalOpen}
						setIsOpen={setCartonModalOpen}
						itemName="carton"
						value={carton}
						setValue={setCarton}
						query={graphql.query.CARTONS_BASIC}
					/>
					<Checkbox checked={inheritCartonHistory} onChange={(e) => setInheritCartonHistory(e.currentTarget.checked)}>
						Inherit Carton Tracking History (ignore if not changing carton)
					</Checkbox>
				</>
			</FormControl>

			{breakLine}

			<FormControl label="Bonus Loyalty Points" error={errors.loyaltyPoints ? errors.loyaltyPoints.message : ""} positive="">
				<Input name="loyaltyPoints" type="number" inputRef={register} />
			</FormControl>

			<FormControl label="Bonus Loyalty Points Expiration" caption="YYYY/MM/DD" error="" positive="">
				<div className={flexStyle}>
					<div
						style={{
							width: "160px",
							marginRight: "10px",
						}}
					>
						<Datepicker value={loyaltyPointsExpire} onChange={({ date }) => setLoyaltyPointsExpire(date as Date)} />
					</div>
					<div style={{ width: "120px" }}>
						<TimePicker value={loyaltyPointsExpire} onChange={setLoyaltyPointsExpire} />
					</div>
					{loyaltyPointsExpire <= new Date() && (
						<Notification
							kind="negative"
							overrides={{
								Body: {
									style: {
										width: "unset",
										marginTop: "0",
										marginBottom: "0",
										marginLeft: "10px",
									},
								},
							}}
						>
							EXPIRED
						</Notification>
					)}
				</div>
			</FormControl>

			{breakLine}
			<FormControl label="Description" error={errors.description ? errors.description.message : ""} positive="">
				<Input name="description" inputRef={register} />
			</FormControl>

			{breakLine}

			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/products")}>
					Cancel
				</Button>
				{product && !isNewProduct ? (
					<Spaced>
						<Button
							type="button"
							kind="secondary"
							isLoading={mutArchiveProduct.loading || mutUnarchiveProduct.loading}
							onClick={toggleArchive}
							startEnhancer={<FontAwesomeIcon icon={["fas", product.archived ? "undo" : "archive"]} size="lg" />}
						>
							{product.archived ? "Unarchive" : "Archive"}
						</Button>
						<Button isLoading={mutUpdateProduct.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
							{timedOut ? "Timed out... Try again" : "Save"}
						</Button>
					</Spaced>
				) : (
					<Button isLoading={mutUpdateProduct.loading}>Create Product</Button>
				)}
			</Spread>
		</form>
	)

	return (
		<CenteredPage>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={["fal", "steak"]} size="3x" />
					<H1>{isNewProduct ? "New Product" : code} </H1>
				</Spaced>

				<div className={flexColumnStyle}>
					<div className={viewLinkContainerStyle}>
						<div>
							{product && product.registeredBy && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/consumer/${product.registeredBy?.wechatID}#loyalty`)}
									overrides={ButtonMarginLeftOverride}
								>
									<FontAwesomeIcon icon={["fal", "smile"]} />
									<div className={marginLeftStyle}>Registered</div>
								</Button>
							)}
						</div>
						<div>
							{product && product.order && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/order/${product.order.code}#products`)}
									overrides={ButtonMarginLeftOverride}
								>
									<FontAwesomeIcon icon={["fas", "shopping-cart"]} />
									<div className={marginLeftStyle}>{product.order.code}</div>
								</Button>
							)}
						</div>
						<div>
							{product && product.contract && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/contract/${product.contract.code}#products`)}
									overrides={ButtonMarginLeftOverride}
								>
									<FontAwesomeIcon icon={["fas", "file-contract"]} />
									<div className={marginLeftStyle}>{product.contract.name}</div>
								</Button>
							)}
						</div>
						<div>
							{product && product.distributor && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/distributor/${product.distributor.code}#products`)}
									overrides={ButtonMarginLeftOverride}
								>
									<FontAwesomeIcon icon={["fas", "shopping-basket"]} />
									<div className={marginLeftStyle}>{product.distributor.code}</div>
								</Button>
							)}
						</div>
						<div>
							{product && product.sku && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/sku/${product.sku.code}#products`)}
									overrides={ButtonMarginLeftOverride}
								>
									<FontAwesomeIcon icon={["fal", "barcode-alt"]} />
									<div className={marginLeftStyle}>{product.sku.code}</div>
								</Button>
							)}
						</div>
					</div>
					{product && product.carton && (
						<div className={breadCrumbContainerStyle}>
							{product.carton.pallet?.container && (
								<>
									<StyledLink href={`/portal/container/${product.carton.pallet.container.code}#pallets`} className={breadCrumbStyle}>
										<FontAwesomeIcon icon={["fal", "container-storage"]} />
										<div className={marginLeftStyle}>{product.carton.pallet.container.code}</div>
									</StyledLink>
									{breadCrumbChevron}
								</>
							)}

							{product.carton.pallet && (
								<>
									<StyledLink href={`/portal/pallet/${product.carton.pallet.code}#cartons`} className={breadCrumbStyle}>
										<FontAwesomeIcon icon={["fal", "pallet-alt"]} />
										<div className={marginLeftStyle}>{product.carton.pallet.code}</div>
									</StyledLink>
									{breadCrumbChevron}
								</>
							)}

							<StyledLink href={`/portal/carton/${product.carton.code}#products`} className={breadCrumbStyle}>
								<FontAwesomeIcon icon={["fas", "box"]} />
								<div className={marginLeftStyle}>{product.carton.code}</div>
							</StyledLink>
							{breadCrumbChevron}

							<StyledLink className={breadCrumbStyle}>
								<FontAwesomeIcon icon={["fas", "steak"]} />
								<div className={marginLeftStyle}>{product.code}</div>
							</StyledLink>
						</div>
					)}
				</div>
			</Spread>

			{isNewProduct ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/product/${code}${activeKey}`)
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

					{!isNewProduct && product && (
						<Tab
							key="#tracking"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fas", "truck-moving"]} />
									<div>Tracking</div>
								</Spaced>
							}
						>
							{queryTransactions.loading ? (
								<LoadingSimple />
							) : (
								<TransactionsView transactions={queryTransactions.data?.productByID.transactions} showManifestDetails />
							)}
						</Tab>
					)}

					{!isNewProduct && product && (
						<Tab
							key="#qrcodes"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "qrcode"]} />
									<div>QR Codes</div>
								</Spaced>
							}
						>
							<div className={qrItemsContainer}>
								<div className={qrItem}>
									<QRCode value={viewLink} size={250} quietZone={2} />
									<div>View Link</div>
									<a href={viewLink} className={qrLink} target="_blank">
										{viewLink}
									</a>
								</div>
								<div className={qrItem}>
									<QRCode value={registerLink} size={250} quietZone={2} />
									<div>Register Link</div>
									<a href={registerLink} className={qrLink} target="_blank">
										{registerLink}
									</a>
								</div>
								{showDemoQRCode || (
									<div className={qrItem}>
										<Button onClick={() => setShowDemoQRCode(true)}>Show Demo Register Link</Button>
									</div>
								)}
								{showDemoQRCode && (
									<div className={qrItem}>
										<QRCode value={registerLinkDemo} size={250} quietZone={2} />
										<div>Register Link (Demo)</div>
										<a href={registerLinkDemo} className={qrLink} target="_blank">
											{registerLinkDemo}
										</a>
									</div>
								)}
							</div>
						</Tab>
					)}
				</Tabs>
			)}
		</CenteredPage>
	)
}

export default ProductEdit
