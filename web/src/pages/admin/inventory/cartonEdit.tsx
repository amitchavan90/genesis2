import * as React from "react"
import { useLazyQuery, useMutation, useQuery } from "@apollo/react-hooks"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useStyletron } from "baseui"
import { Button } from "baseui/button"
import { Datepicker } from "baseui/datepicker"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { StyledLink } from "baseui/link"
import { Modal } from "baseui/modal"
import { Notification } from "baseui/notification"
import { Value } from "baseui/select"
import { Tab, Tabs } from "baseui/tabs"
import { H1, H3 } from "baseui/typography"
import { useForm } from "react-hook-form"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { invalidateListQueries } from "../../../apollo"
import { CenteredPage } from "../../../components/common"
import { ErrorNotification } from "../../../components/errorBox"
import { ItemList } from "../../../components/itemList"
import { LoadingSimple } from "../../../components/loading"
import { QRPreview } from "../../../components/qrCodePreview"
import { Spaced } from "../../../components/spaced"
import { Spread } from "../../../components/spread"
import { TablePickerSelect } from "../../../components/tablePickerSelect"
import { TransactionsView } from "../../../components/transactionsView"
import { graphql } from "../../../graphql"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"
import { ButtonMarginLeftOverride, paddingZero } from "../../../themeOverrides"
import { ActionItemSet } from "../../../types/actions"
import { FilterOption } from "../../../types/enums"
import { Carton, Settings, Transaction } from "../../../types/types"

const fallbackImage = require("../../../assets/images/default-image.png")

type FormData = {
	quantity: number
	weight: string
	description: string
	meatType: string
}

const CartonEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewCarton = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	// Get Carton
	const [carton, setCarton] = React.useState<Carton>()
	const { data, loading, error, refetch } = useQuery<{ carton: Carton }>(graphql.query.CARTON, {
		variables: { code },
		fetchPolicy: isNewCarton ? "cache-only" : undefined, // prevent query if new
	})

	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)

	// Get Transactions
	const [transactions, setTransactions] = React.useState<Transaction[]>()
	const [getTransactions, queryTransactions] = useLazyQuery<{ carton: Carton }>(graphql.query.CARTON_TRANSACTIONS, {
		fetchPolicy: "network-only",
		pollInterval: 60000, // refetches every 60 secs
	})

	const [spreadSheetLink, setSpreadSheetLink] = React.useState("")
	// On load transactions
	React.useEffect(() => {
		if (!queryTransactions.data || !queryTransactions.data.carton) return
		setTransactions(queryTransactions.data.carton.transactions)
	}, [queryTransactions])

	// Mutations
	const [updateCarton, mutUpdateCarton] = useMutation(
		isNewCarton ? graphql.mutation.CREATE_CARTON : graphql.mutation.UPDATE_CARTON,
	)
	const [archiveCarton, mutArchiveCarton] = useMutation<{ cartonArchive: Carton }>(graphql.mutation.ARCHIVE_CARTON)
	const [unarchiveCarton, mutUnarchiveCarton] = useMutation<{ cartonUnarchive: Carton }>(
		graphql.mutation.UNARCHIVE_CARTON,
	)

	const toggleArchive = () => {
		if (!carton) return

		if (carton.archived) {
			unarchiveCarton({
				variables: { id: carton.id },
				update: (cache: any) => invalidateListQueries(cache, "cartons"),
			})
		} else {
			archiveCarton({
				variables: { id: carton.id },
				update: (cache: any) => invalidateListQueries(cache, "cartons"),
			})
		}
	}

	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [pallet, setPallet] = React.useState<Value>()
	const [processedAt, setProcessedAt] = React.useState<Date>()

	// modal
	const [palletModalOpen, setPalletModalOpen] = React.useState(false)

	// form
	const { register, setValue, handleSubmit, errors } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ quantity, weight, description, meatType }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		if (isNewCarton) {
			updateCarton({
				variables: {
					input: {
						description,
						meatType,
						quantity,
						palletID: pallet && pallet.length > 0 ? pallet[0].id : "-",
					},
				},
				update: (cache: any) => invalidateListQueries(cache, "cartons"),
			})
			return
		}

		if (!carton) return

		const input = {
			code,
			weight,
			processedAt,
			description,
			meatType,

			palletID: pallet && pallet.length > 0 ? pallet[0].id : "-",
		}

		promiseTimeout(updateCarton({ variables: { id: carton.id, input } })).catch((reason) => {
			if (reason !== TIMED_OUT) return
			setTimedOut(true)
		})
	})

	// On load carton
	React.useEffect(() => {
		if (!data || !data.carton) return
		setCarton(data.carton)
	}, [data, loading, error])
	React.useEffect(() => {
		// Load contract on tracking tab visit
		if (!transactions && carton && activeKey == "#tracking") getTransactions({ variables: { code } })

		if (activeKey != "#details") return
		if (!carton) return
		setValue("code", carton.code)
		setValue("weight", carton.weight)
		setValue("description", carton.description)
		setValue("meatType", carton.meatType)

		if (carton.processedAt) setProcessedAt(new Date(carton.processedAt))

		if (carton.pallet) setPallet([{ id: carton.pallet.id, label: carton.pallet.code }])
	}, [carton, activeKey])

	// On mutation (update/create carton)
	React.useEffect(() => {
		if (!mutUpdateCarton.data) return

		if (isNewCarton) {
			if (!mutUpdateCarton.data.cartonCreate) return
			setChangeSuccess(true)
			setSpreadSheetLink(mutUpdateCarton.data.cartonCreate)
			return
		}

		if (!mutUpdateCarton.data.cartonUpdate) return

		setCarton(mutUpdateCarton.data.cartonUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateCarton])

	React.useEffect(() => {
		if (!mutArchiveCarton.data?.cartonArchive) return
		setCarton(mutArchiveCarton.data.cartonArchive)
	}, [mutArchiveCarton])
	React.useEffect(() => {
		if (!mutUnarchiveCarton.data?.cartonUnarchive) return
		setCarton(mutUnarchiveCarton.data.cartonUnarchive)
	}, [mutUnarchiveCarton])

	// photo view on edit page
	const [showPreviewModal, setShowPreviewModal] = React.useState<boolean>()
	const [previewModalImage, setPreviewModalImage] = React.useState<string>()

	let transactionWithPhotos: Transaction | undefined
	if (!isNewCarton && carton && transactions) {
		transactionWithPhotos = transactions.find((tx) => tx.photos)
	}

	// Styling
	const [css] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "10px",
	})
	const flexColumnStyle = css({
		display: "flex",
		flexDirection: "column",
	})
	const viewLinkContainerStyle = css({
		display: "flex",
		alignSelf: "flex-end",
	})

	const imagePreviewStyle = css({
		width: "100%",
		cursor: "pointer",
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

	const createSuccessStyle = css({
		textAlign: "center",
		marginTop: "30px",
	})
	const breakLineStyle = css({
		border: "1px solid rgba(0, 0, 0, 0.1)",
		marginBottom: "10px",
	})
	const imageBox = css({
		width: "100%",
		display: "block",
		objectFit: "contain",
		cursor: "pointer",
		marginRight: "15px",
	})
	const photos = css({
		display: "flex",
		maxWidth: "500px",
		marginTop: "8px",
	})
	const photoDescriptionStyle = css({
		fontSize: "0.8rem",
		lineHeight: "0.8rem",
		color: "#0000008f",
		marginTop: "5px",
		marginBottom: "5px",
	})

	const breakLine = <div className={breakLineStyle} />

	if (!isNewCarton && !carton) {
		return <LoadingSimple />
	}

	if (!settingsQuery.data || !settingsQuery.data.settings || !settingsQuery.data.settings.adminHost) {
		return <LoadingSimple />
	}

	if (isNewCarton && changeSuccess) {
		// Successfully created cartons
		return (
			<CenteredPage>
				<div className={createSuccessStyle}>
					<FontAwesomeIcon icon={["far", "check-circle"]} size="10x" color="#1db954" />
					<H3>Cartons created</H3>

					<Spread>
						<Button type="button" kind="secondary" onClick={() => history.push("/portal/cartons")}>
							Back
						</Button>
						<Button kind="secondary" type="button" onClick={() => window.open(spreadSheetLink)}>
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
					The carton has been updated.
				</Notification>
			)}
			{mutUpdateCarton.error && <ErrorNotification message={mutUpdateCarton.error.message} />}
			<form onSubmit={onSubmit}>
				{isNewCarton && (
					<>
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
						{breakLine}
					</>
				)}

				{!isNewCarton && carton && (
					<>
						<FormControl label="Weight" error={errors.weight ? errors.weight.message : ""} positive="">
							<Input name="weight" inputRef={register} />
						</FormControl>

						<FormControl label="Meat Type" error={errors.meatType ? errors.meatType.message : ""} positive="">
							<Input name="meatType" inputRef={register} />
						</FormControl>

						{breakLine}

						<FormControl label="Processed Date" caption="YYYY/MM/DD" error="" positive="">
							<div className={flexColumnStyle}>
								<div
									style={{
										width: "160px",
										marginRight: "10px",
									}}>
									<Datepicker
										value={processedAt}
										maxDate={new Date()}
										onChange={({ date }) => setProcessedAt(date as Date)}
									/>
								</div>
							</div>
						</FormControl>

						{breakLine}
					</>
				)}

				<FormControl label={isNewCarton ? "Pallet (optional)" : "Pallet"}>
					<TablePickerSelect
						isOpen={palletModalOpen}
						setIsOpen={setPalletModalOpen}
						hash="cartons"
						itemName="pallet"
						value={pallet}
						setValue={setPallet}
						query={graphql.query.PALLETS_BASIC}
						queryName="pallets"
					/>
				</FormControl>

				{breakLine}

				<FormControl label="Description" error={errors.description ? errors.description.message : ""} positive="">
					<Input name="description" inputRef={register} />
				</FormControl>

				{transactionWithPhotos && transactionWithPhotos.photos && (
					<>
						{breakLine}

						<FormControl label="Photos" error={errors.weight ? errors.weight.message : ""} positive="">
							{/* Photos of carton and product */}
							<div className={photos}>
								{transactionWithPhotos.photos.cartonPhoto && (
									<div style={{ marginRight: "10px" }}>
										<div className={photoDescriptionStyle}>photo of carton:</div>
										<img
											className={imageBox}
											src={transactionWithPhotos.photos.cartonPhoto.file_url || fallbackImage}
											alt="carton image"
											onClick={() => {
												if (!transactionWithPhotos || !transactionWithPhotos.photos.cartonPhoto?.file_url) return
												setPreviewModalImage(transactionWithPhotos.photos.cartonPhoto?.file_url)
												setShowPreviewModal(true)
											}}
										/>
									</div>
								)}
								{transactionWithPhotos.photos.productPhoto && (
									<div>
										<div className={photoDescriptionStyle}>photo of product:</div>
										<img
											className={imageBox}
											src={transactionWithPhotos.photos.productPhoto.file_url || fallbackImage}
											alt="product image"
											onClick={() => {
												if (!transactionWithPhotos || !transactionWithPhotos.photos.productPhoto?.file_url) return
												setPreviewModalImage(transactionWithPhotos.photos.productPhoto?.file_url)
												setShowPreviewModal(true)
											}}
										/>
									</div>
								)}
							</div>
						</FormControl>

						<Modal
							isOpen={showPreviewModal}
							onClose={() => setShowPreviewModal(false)}
							overrides={{
								Dialog: {
									style: {
										width: "unset",
										backgroundColor: "unset",
									},
								},
							}}>
							<img className={imagePreviewStyle} src={previewModalImage} />
						</Modal>
					</>
				)}

				{breakLine}

				<Spread>
					<Button type="button" kind="secondary" onClick={() => history.push("/portal/cartons")}>
						Cancel
					</Button>
					{carton && !isNewCarton ? (
						<Spaced>
							<Button
								type="button"
								kind="secondary"
								isLoading={mutArchiveCarton.loading || mutUnarchiveCarton.loading}
								onClick={toggleArchive}
								startEnhancer={<FontAwesomeIcon icon={["fas", carton.archived ? "undo" : "archive"]} size="lg" />}>
								{carton.archived ? "Unarchive" : "Archive"}
							</Button>
							<Button
								isLoading={mutUpdateCarton.loading && !timedOut}
								startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
								{timedOut ? "Timed out... Try again" : "Save"}
							</Button>
						</Spaced>
					) : (
						<Button isLoading={mutUpdateCarton.loading}>Create Cartons</Button>
					)}
				</Spread>
			</form>
		</>
	)

	return (
		<CenteredPage>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={["fal", "box"]} size="3x" />
					<H1>{isNewCarton ? "New Cartons" : code}</H1>
				</Spaced>

				<div className={flexColumnStyle}>
					<div className={viewLinkContainerStyle}>
						{!isNewCarton && (
							<Button
								disabled={!carton?.spreadsheetLink}
								kind="secondary"
								type="button"
								onClick={() => {
									window.open(carton?.spreadsheetLink)
								}}>
								<FontAwesomeIcon icon={["fas", "file-excel"]} color="#1db954" />
								<div className={marginLeftStyle}>Download</div>
							</Button>
						)}
						<div>
							{carton && carton.order && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/order/${carton.order.code}#products`)}
									overrides={ButtonMarginLeftOverride}>
									<FontAwesomeIcon icon={["fas", "shopping-cart"]} />
									<div className={marginLeftStyle}>{carton.order.code}</div>
								</Button>
							)}
						</div>
						<div>
							{carton && carton.distributor && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/distributor/${carton.distributor.code}#products`)}
									overrides={ButtonMarginLeftOverride}>
									<FontAwesomeIcon icon={["fas", "shopping-basket"]} />
									<div className={marginLeftStyle}>{carton.distributor.code}</div>
								</Button>
							)}
						</div>
						<div>
							{carton && carton.sku && (
								<Button
									kind="secondary"
									type="button"
									onClick={() => history.push(`/portal/sku/${carton.sku.code}#products`)}
									overrides={ButtonMarginLeftOverride}>
									<FontAwesomeIcon icon={["fal", "barcode-alt"]} />
									<div className={marginLeftStyle}>{carton.sku.code}</div>
								</Button>
							)}
						</div>
					</div>
					{carton && carton.pallet && (
						<div className={breadCrumbContainerStyle}>
							{carton.pallet && carton.pallet.container && (
								<>
									<StyledLink
										href={`/portal/container/${carton.pallet.container.code}#pallets`}
										className={breadCrumbStyle}>
										<FontAwesomeIcon icon={["fal", "container-storage"]} />
										<div className={marginLeftStyle}>{carton.pallet.container.code}</div>
									</StyledLink>
									{breadCrumbChevron}
								</>
							)}

							<StyledLink href={`/portal/pallet/${carton.pallet.code}#cartons`} className={breadCrumbStyle}>
								<FontAwesomeIcon icon={["fal", "pallet-alt"]} />
								<div className={marginLeftStyle}>{carton.pallet.code}</div>
							</StyledLink>
							{breadCrumbChevron}

							<StyledLink className={breadCrumbStyle}>
								<FontAwesomeIcon icon={["fas", "box"]} />
								<div className={marginLeftStyle}>{carton.code}</div>
							</StyledLink>
						</div>
					)}
				</div>
			</Spread>

			{isNewCarton ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/carton/${code}${activeKey}`)
					}}
					activeKey={activeKey}
					overrides={{
						TabContent: {
							style: { ...paddingZero },
						},
						TabBar: {
							style: { ...paddingZero },
						},
					}}>
					{/* Details/Edit */}
					<Tab
						key="#details"
						title={
							<Spaced>
								<FontAwesomeIcon icon={["fal", "pencil"]} />
								<div>Details</div>
							</Spaced>
						}>
						{editForm}
					</Tab>

					{/* Tracking */}
					{!isNewCarton && carton && (
						<Tab
							key="#tracking"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fas", "truck-moving"]} />
									<div>Tracking</div>
								</Spaced>
							}>
							{queryTransactions.loading ? (
								<LoadingSimple />
							) : (
								<TransactionsView transactions={transactions} showManifestDetails />
							)}
						</Tab>
					)}

					{/* Products */}
					{!isNewCarton && carton && (
						<Tab
							key="#products"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "steak"]} />
									<div>Products</div>
								</Spaced>
							}>
							{/* Warning(s) */}
							{carton.productCount > 0 && (
								<>
									{!carton.order && (
										<Notification kind="warning" overrides={{ Body: { style: { width: "unset" } } }}>
											Not all of the products in this Carton belong to the same Order
										</Notification>
									)}
									{!carton.sku && (
										<Notification kind="warning" overrides={{ Body: { style: { width: "unset" } } }}>
											Not all of the products in this Carton have the same SKU
										</Notification>
									)}
								</>
							)}

							{/* List */}
							<ItemList
								cartonID={carton.id}
								itemName="product"
								query={graphql.query.PRODUCTS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_PRODUCT}
								extraFilterOptions={[
									{ label: "Not in Order", id: FilterOption.ProductWithoutOrder },
									{ label: "No SKU", id: FilterOption.ProductWithoutSKU },
								]}
								itemLinks={["order", "sku", "contract", "distributor"]}
								actions={ActionItemSet.Products}
								onListUpdate={refetch}
								showQRCodesToggle
							/>
						</Tab>
					)}

					{/* QR Code */}
					{!isNewCarton && carton && (
						<Tab
							key="#qrcode"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "qrcode"]} />
									<div>QR Code</div>
								</Spaced>
							}>
							<QRPreview item={carton} />
						</Tab>
					)}
				</Tabs>
			)}
		</CenteredPage>
	)
}

export default CartonEdit
