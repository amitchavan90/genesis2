import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { Pallet } from "../../../types/types"
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
import { Value } from "baseui/select"
import { CenteredPage, LatestTrackActionColumn } from "../../../components/common"
import { ItemList } from "../../../components/itemList"
import { ActionItemSet } from "../../../types/actions"
import { StyledLink } from "baseui/link"
import { invalidateListQueries } from "../../../apollo"
import { QRPreview } from "../../../components/qrCodePreview"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"
import { TablePickerSelect } from "../../../components/tablePickerSelect"

type FormData = {
	quantity: number
	description: string
}

const PalletEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewPallet = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	// Get Pallet
	const [pallet, setPallet] = React.useState<Pallet>()
	const { data, loading, error } = useQuery<{ pallet: Pallet }>(graphql.query.PALLET, {
		variables: { code },
		fetchPolicy: isNewPallet ? "cache-only" : undefined, // prevent query if new
	})

	const [spreadSheetLink, setSpreadSheetLink] = React.useState("")

	// Mutations
	const [updatePallet, mutUpdatePallet] = useMutation(isNewPallet ? graphql.mutation.CREATE_PALLET : graphql.mutation.UPDATE_PALLET)
	const [archivePallet, mutArchivePallet] = useMutation<{ palletArchive: Pallet }>(graphql.mutation.ARCHIVE_PALLET)
	const [unarchivePallet, mutUnarchivePallet] = useMutation<{ palletUnarchive: Pallet }>(graphql.mutation.UNARCHIVE_PALLET)

	// modal
	const [containerModalOpen, setContainerModalOpen] = React.useState(false)

	const toggleArchive = () => {
		if (!pallet) return

		if (pallet.archived) {
			unarchivePallet({
				variables: { id: pallet.id },
				update: (cache: any) => invalidateListQueries(cache, "pallets"),
			})
		} else {
			archivePallet({
				variables: { id: pallet.id },
				update: (cache: any) => invalidateListQueries(cache, "pallets"),
			})
		}
	}

	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [container, setContainer] = React.useState<Value>()

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ quantity, description }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		if (isNewPallet) {
			updatePallet({
				variables: {
					input: {
						quantity,
						description,
						containerID: container && container.length > 0 ? container[0].id : "-",
					},
				},
				update: (cache: any) => invalidateListQueries(cache, "pallets"),
			})
			return
		}

		if (!pallet) return

		const input = {
			code,
			description,
			containerID: container && container.length > 0 ? container[0].id : "-",
		}

		promiseTimeout(updatePallet({ variables: { id: pallet.id, input } })).catch(reason => {
			if (reason !== TIMED_OUT) return
			setTimedOut(true)
		})
	})

	// On load pallet
	React.useEffect(() => {
		if (!data || !data.pallet) return
		setPallet(data.pallet)
	}, [data, loading, error])
	React.useEffect(() => {
		if (activeKey != "#details") return
		if (!pallet) return
		setValue("code", pallet.code)
		setValue("description", pallet.description)
		if (pallet.container) setContainer([{ id: pallet.container.id, label: pallet.container.code }])
	}, [pallet, activeKey])

	// On mutation (update/create pallet)
	React.useEffect(() => {
		if (!mutUpdatePallet.data) return

		if (isNewPallet) {
			if (!mutUpdatePallet.data.palletCreate) return
			setChangeSuccess(true)
			setSpreadSheetLink(mutUpdatePallet.data.palletCreate)
			return
		}

		if (!mutUpdatePallet.data.palletUpdate) return

		setPallet(mutUpdatePallet.data.palletUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdatePallet])

	React.useEffect(() => {
		if (!mutArchivePallet.data?.palletArchive) return
		setPallet(mutArchivePallet.data.palletArchive)
	}, [mutArchivePallet])
	React.useEffect(() => {
		if (!mutUnarchivePallet.data?.palletUnarchive) return
		setPallet(mutUnarchivePallet.data.palletUnarchive)
	}, [mutUnarchivePallet])

	// Styling
	const [css, theme] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "10px",
	})
	const breadCrumbContainerStyle = css({
		display: "flex",
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
	const breakLine = <div className={breakLineStyle} />

	if (!isNewPallet && !pallet) {
		return <LoadingSimple />
	}

	if (isNewPallet && changeSuccess) {
		// Successfully created pallets
		return (
			<CenteredPage>
				<div className={createSuccessStyle}>
					<FontAwesomeIcon icon={["far", "check-circle"]} size="10x" color="#1db954" />
					<H3>Pallets created</H3>

					<Spread>
						<Button type="button" kind="secondary" onClick={() => history.push("/portal/pallets")}>
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
					The pallet has been updated.
				</Notification>
			)}
			{mutUpdatePallet.error && <ErrorNotification message={mutUpdatePallet.error.message} />}
			<form onSubmit={onSubmit}>
				{isNewPallet && (
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

				<FormControl label={isNewPallet ? "Container (optional)" : "Container"}>
					<TablePickerSelect
						isOpen={containerModalOpen}
						setIsOpen={setContainerModalOpen}
						hash="pallets"
						itemName="container"
						value={container}
						setValue={setContainer}
						query={graphql.query.PALLETS_BASIC}
						queryName="containers"
					/>
				</FormControl>

				{breakLine}

				<FormControl label="Description" error={errors.description ? errors.description.message : ""} positive="">
					<Input name="description" inputRef={register} />
				</FormControl>

				{breakLine}

				<Spread>
					<Button type="button" kind="secondary" onClick={() => history.push("/portal/pallets")}>
						Cancel
					</Button>
					{pallet && !isNewPallet ? (
						<Spaced>
							<Button
								type="button"
								kind="secondary"
								isLoading={mutArchivePallet.loading || mutUnarchivePallet.loading}
								onClick={toggleArchive}
								startEnhancer={<FontAwesomeIcon icon={["fas", pallet.archived ? "undo" : "archive"]} size="lg" />}
							>
								{pallet.archived ? "Unarchive" : "Archive"}
							</Button>
							<Button isLoading={mutUpdatePallet.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
								{timedOut ? "Timed out... Try again" : "Save"}
							</Button>
						</Spaced>
					) : (
						<Button isLoading={mutUpdatePallet.loading}>Create Pallets</Button>
					)}
				</Spread>
			</form>
		</>
	)

	return (
		<CenteredPage>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={["fal", "pallet-alt"]} size="3x" />
					<H1>{isNewPallet ? "New Pallets" : code} </H1>
				</Spaced>
				<div className={breadCrumbContainerStyle}>
					{pallet && pallet.container && (
						<>
							<StyledLink href={`/portal/container/${pallet.container.code}#pallets`} className={breadCrumbStyle}>
								<FontAwesomeIcon icon={["fal", "container-storage"]} />
								<div className={marginLeftStyle}>{pallet.container.code}</div>
							</StyledLink>
							{breadCrumbChevron}

							<StyledLink className={breadCrumbStyle}>
								<FontAwesomeIcon icon={["fal", "pallet-alt"]} />
								<div className={marginLeftStyle}>{pallet.code}</div>
							</StyledLink>
						</>
					)}
				</div>
			</Spread>

			{isNewPallet ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/pallet/${code}${activeKey}`)
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

					{!isNewPallet && pallet && (
						<Tab
							key="#cartons"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "box"]} />
									<div>Cartons</div>
								</Spaced>
							}
						>
							<ItemList
								palletID={pallet.id}
								itemName="carton"
								query={graphql.query.CARTONS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_CARTON}
								hash="products"
								columns={[
									{ label: "Weight", value: "weight" },
									{ label: "Products Amount", value: "productCount" },
									{
										label: "Last Track Action",
										value: "latestTrackAction",
										resolver: row => <LatestTrackActionColumn value={row.latestTrackAction} />,
									},
									{ label: "Description", value: "description" },
									{ label: "Meat Type", value: "meatType" },
									{
										label: "Date Created",
										value: "createdAt",
										dateTime: true,
									},
								]}
								actions={ActionItemSet.Cartons}
								showQRCodesToggle
							/>
						</Tab>
					)}

					{/* QR Code */}
					{!isNewPallet && pallet && (
						<Tab
							key="#qrcode"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "qrcode"]} />
									<div>QR Code</div>
								</Spaced>
							}
						>
							<QRPreview item={pallet} />
						</Tab>
					)}
				</Tabs>
			)}
		</CenteredPage>
	)
}

export default PalletEdit
