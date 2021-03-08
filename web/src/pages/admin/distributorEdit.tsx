import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../graphql"
import { Distributor } from "../../types/types"
import { useForm } from "react-hook-form"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { H1 } from "baseui/typography"
import { Button } from "baseui/button"
import { Notification } from "baseui/notification"
import { Spaced } from "../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { LoadingSimple } from "../../components/loading"
import { ErrorNotification } from "../../components/errorBox"
import { CenteredPage } from "../../components/common"
import { Spread } from "../../components/spread"
import { invalidateListQueries } from "../../apollo"
import { Tabs, Tab } from "baseui/tabs"
import { paddingZero } from "../../themeOverrides"
import { ItemList } from "../../components/itemList"
import { FilterOption } from "../../types/enums"
import { ActionItemSet } from "../../types/actions"
import { promiseTimeout, TIMED_OUT } from "../../helpers/timeout"

type FormData = {
	name: string
	code: string
}

const DistributorEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewDistributor = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	// Get Distributor
	const [distributor, setDistributor] = React.useState<Distributor>()
	const { data, loading, error, refetch } = useQuery<{ distributor: Distributor }>(graphql.query.DISTRIBUTOR, {
		variables: { code },
		fetchPolicy: isNewDistributor ? "cache-only" : undefined, // prevent query if new
	})

	// Mutations
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [updateDistributor, mutUpdateDistributor] = useMutation(isNewDistributor ? graphql.mutation.CREATE_DISTRIBUTOR : graphql.mutation.UPDATE_DISTRIBUTOR)
	const [archiveDistributor, mutArchiveDistributor] = useMutation<{ distributorArchive: Distributor }>(graphql.mutation.ARCHIVE_DISTRIBUTOR)
	const [unarchiveDistributor, mutUnarchiveDistributor] = useMutation<{ distributorUnarchive: Distributor }>(graphql.mutation.UNARCHIVE_DISTRIBUTOR)

	const toggleArchive = () => {
		if (!distributor) return

		if (distributor.archived) {
			unarchiveDistributor({
				variables: { id: distributor.id },
				update: (cache: any) => invalidateListQueries(cache, "distributors"),
			})
		} else {
			archiveDistributor({
				variables: { id: distributor.id },
				update: (cache: any) => invalidateListQueries(cache, "distributors"),
			})
		}
	}

	// Form submission
	const { register, setValue, handleSubmit, errors } = useForm<FormData>()
	const onSubmit = handleSubmit(({ name, code }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		const input = { name, code }

		if (isNewDistributor) {
			updateDistributor({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "distributors"),
			})
		} else if (distributor) {
			promiseTimeout(updateDistributor({ variables: { id: distributor.id, input } })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	// On load distributor
	React.useEffect(() => {
		if (!data || !data.distributor) return
		setDistributor(data.distributor)
	}, [data, loading, error])
	React.useEffect(() => {
		if (activeKey != "#details") return
		if (!distributor) return
		setValue("name", distributor.name)
		setValue("code", distributor.code)
	}, [distributor, activeKey])

	// On mutation (update/create distributor)
	React.useEffect(() => {
		if (!mutUpdateDistributor.data) return

		if (isNewDistributor) {
			if (mutUpdateDistributor.data.distributorCreate) {
				history.push(`/portal/distributor/${mutUpdateDistributor.data.distributorCreate.code}`)
			}
			return
		}

		if (!mutUpdateDistributor.data.distributorUpdate) return

		setDistributor(mutUpdateDistributor.data.distributorUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateDistributor.data, mutUpdateDistributor.loading])

	React.useEffect(() => {
		if (!mutArchiveDistributor.data?.distributorArchive) return
		setDistributor(mutArchiveDistributor.data.distributorArchive)
	}, [mutArchiveDistributor])
	React.useEffect(() => {
		if (!mutUnarchiveDistributor.data?.distributorUnarchive) return
		setDistributor(mutUnarchiveDistributor.data.distributorUnarchive)
	}, [mutUnarchiveDistributor])

	if (!isNewDistributor && !distributor) {
		return <LoadingSimple />
	}

	const editForm = (
		<form onSubmit={onSubmit}>
			{changeSuccess && (
				<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
					Distributor has been updated.
				</Notification>
			)}

			{mutUpdateDistributor.error && <ErrorNotification message={mutUpdateDistributor.error.message} />}

			<FormControl label="Name" error={errors.name ? errors.name.message : ""} positive="">
				<Input name="name" inputRef={register({ required: "Required" })} />
			</FormControl>

			<FormControl label="Code" error={errors.code ? errors.code.message : ""} positive="">
				<Input name="code" inputRef={register({ required: "Required" })} />
			</FormControl>

			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/distributors")}>
					Cancel
				</Button>
				{distributor && !isNewDistributor ? (
					<Spaced>
						<Button
							type="button"
							kind="secondary"
							isLoading={mutArchiveDistributor.loading || mutUnarchiveDistributor.loading}
							onClick={toggleArchive}
							startEnhancer={<FontAwesomeIcon icon={["fas", distributor.archived ? "undo" : "archive"]} size="lg" />}
						>
							{distributor.archived ? "Unarchive" : "Archive"}
						</Button>
						<Button isLoading={mutUpdateDistributor.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
							{timedOut ? "Timed out... Try again" : "Save"}
						</Button>
					</Spaced>
				) : (
					<Button isLoading={mutUpdateDistributor.loading} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
						Create Distributor
					</Button>
				)}
			</Spread>
		</form>
	)

	return (
		<CenteredPage>
			<Spaced>
				<FontAwesomeIcon icon={["fal", "shopping-basket"]} size="3x" />
				<H1>{isNewDistributor || !distributor ? "New Distributor" : distributor.name}</H1>
			</Spaced>

			{isNewDistributor ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/distributor/${code}${activeKey}`)
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

					{!isNewDistributor && distributor && (
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
								distributorID={distributor.id}
								itemName="product"
								query={graphql.query.PRODUCTS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_PRODUCT}
								extraFilterOptions={[
									{ label: "Not in Carton", id: FilterOption.ProductWithoutCarton },
									{ label: "Not in Order", id: FilterOption.ProductWithoutOrder },
									{ label: "No SKU", id: FilterOption.ProductWithoutSKU },
								]}
								itemLinks={["order", "sku", "contract", "carton"]}
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

export default DistributorEdit
