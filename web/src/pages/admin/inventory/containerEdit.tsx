import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { Container } from "../../../types/types"
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
import { ActionItemSet } from "../../../types/actions"
import { invalidateListQueries } from "../../../apollo"
import { QRPreview } from "../../../components/qrCodePreview"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"

type FormData = {
	quantity: number
	description: string
}

const ContainerEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewContainer = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	// Get Container
	const [container, setContainer] = React.useState<Container>()
	const { data, loading, error } = useQuery<{ container: Container }>(graphql.query.CONTAINER, {
		variables: { code },
		fetchPolicy: isNewContainer ? "cache-only" : undefined, // prevent query if new
	})

	const [spreadSheetLink, setSpreadSheetLink] = React.useState("")

	// Mutations
	const [updateContainer, mutUpdateContainer] = useMutation(
		isNewContainer ? graphql.mutation.CREATE_CONTAINER : graphql.mutation.UPDATE_CONTAINER,
	)
	const [archiveContainer, mutArchiveContainer] = useMutation<{ containerArchive: Container }>(
		graphql.mutation.ARCHIVE_CONTAINER,
	)
	const [unarchiveContainer, mutUnarchiveContainer] = useMutation<{ containerUnarchive: Container }>(
		graphql.mutation.UNARCHIVE_CONTAINER,
	)

	const toggleArchive = () => {
		if (!container) return

		if (container.archived) {
			unarchiveContainer({
				variables: { id: container.id },
				update: (cache: any) => invalidateListQueries(cache, "containers"),
			})
		} else {
			archiveContainer({
				variables: { id: container.id },
				update: (cache: any) => invalidateListQueries(cache, "containers"),
			})
		}
	}

	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ quantity, description }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		if (isNewContainer) {
			updateContainer({
				variables: { input: { quantity, description } },
				update: (cache: any) => invalidateListQueries(cache, "containers"),
			})
			return
		}

		if (!container) return

		promiseTimeout(updateContainer({ variables: { id: container.id, input: { description } } })).catch((reason) => {
			if (reason !== TIMED_OUT) return
			setTimedOut(true)
		})
	})

	// On load container
	React.useEffect(() => {
		if (!data || !data.container) return
		setContainer(data.container)
	}, [data, loading, error])
	React.useEffect(() => {
		if (activeKey != "#details") return
		if (!container) return
		setValue("code", container.code)
		setValue("description", container.description)
	}, [container, activeKey])

	// On mutation (update/create container)
	React.useEffect(() => {
		if (!mutUpdateContainer.data) return

		if (isNewContainer) {
			if (!mutUpdateContainer.data.containerCreate) return
			setChangeSuccess(true)
			setSpreadSheetLink(mutUpdateContainer.data.containerCreate)
			return
		}

		if (!mutUpdateContainer.data.containerUpdate) return

		setContainer(mutUpdateContainer.data.containerUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateContainer])

	React.useEffect(() => {
		if (!mutArchiveContainer.data?.containerArchive) return
		setContainer(mutArchiveContainer.data.containerArchive)
	}, [mutArchiveContainer])
	React.useEffect(() => {
		if (!mutUnarchiveContainer.data?.containerUnarchive) return
		setContainer(mutUnarchiveContainer.data.containerUnarchive)
	}, [mutUnarchiveContainer])

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
	const breakLine = <div className={breakLineStyle} />

	if (!isNewContainer && !container) {
		return <LoadingSimple />
	}

	if (isNewContainer && changeSuccess) {
		// Successfully created containers
		return (
			<CenteredPage>
				<div className={createSuccessStyle}>
					<FontAwesomeIcon icon={["far", "check-circle"]} size="10x" color="#1db954" />
					<H3>Containers created</H3>

					<Spread>
						<Button type="button" kind="secondary" onClick={() => history.push("/portal/containers")}>
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
					The container has been updated.
				</Notification>
			)}

			{mutUpdateContainer.error && <ErrorNotification message={mutUpdateContainer.error.message} />}

			<form onSubmit={onSubmit}>
				{isNewContainer && (
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
				)}

				{breakLine}
				<FormControl label="Description" error={errors.description ? errors.description.message : ""} positive="">
					<Input name="description" inputRef={register} />
				</FormControl>
				{breakLine}
				<Spread>
					<Button type="button" kind="secondary" onClick={() => history.push("/portal/containers")}>
						Cancel
					</Button>
					{container && !isNewContainer ? (
						<Spaced>
							<Button
								type="button"
								kind="secondary"
								isLoading={mutArchiveContainer.loading || mutUnarchiveContainer.loading}
								onClick={toggleArchive}
								startEnhancer={<FontAwesomeIcon icon={["fas", container.archived ? "undo" : "archive"]} size="lg" />}>
								{container.archived ? "Unarchive" : "Archive"}
							</Button>
							<Button
								isLoading={mutUpdateContainer.loading && !timedOut}
								startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
								{timedOut ? "Timed out... Try again" : "Save"}
							</Button>
						</Spaced>
					) : (
						<Button isLoading={mutUpdateContainer.loading}>Create Containers</Button>
					)}
				</Spread>
			</form>
		</>
	)

	return (
		<CenteredPage>
			<Spaced>
				<FontAwesomeIcon icon={["fal", "container-storage"]} size="3x" />
				<H1>{isNewContainer ? "New Containers" : code} </H1>
			</Spaced>

			{isNewContainer ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/container/${code}${activeKey}`)
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

					{!isNewContainer && container && (
						<Tab
							key="#pallets"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "pallet-alt"]} />
									<div>Pallets</div>
								</Spaced>
							}>
							<ItemList
								containerID={container.id}
								itemName="pallet"
								query={graphql.query.PALLETS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_PALLET}
								hash="cartons"
								columns={[{ label: "Cartons", value: "cartonCount" }]}
								actions={ActionItemSet.Pallets}
							/>
						</Tab>
					)}

					{/* QR Code */}
					{!isNewContainer && container && (
						<Tab
							key="#qrcode"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "qrcode"]} />
									<div>QR Code</div>
								</Spaced>
							}>
							<QRPreview item={container} />
						</Tab>
					)}
				</Tabs>
			)}
		</CenteredPage>
	)
}

export default ContainerEdit
