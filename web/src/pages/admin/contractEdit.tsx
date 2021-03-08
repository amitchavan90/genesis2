import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../graphql"
import { Contract } from "../../types/types"
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
import { Textarea } from "baseui/textarea"
import { Datepicker } from "baseui/datepicker"
import { invalidateListQueries } from "../../apollo"
import { Tabs, Tab } from "baseui/tabs"
import { paddingZero } from "../../themeOverrides"
import { ItemList } from "../../components/itemList"
import { ActionItemSet } from "../../types/actions"
import { FilterOption } from "../../types/enums"
import { promiseTimeout, TIMED_OUT } from "../../helpers/timeout"

type FormData = {
	name: string
	supplierName: string
}

const ContractEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewContract = code === "new"

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")

	// Get Contract
	const [contract, setContract] = React.useState<Contract>()
	const { data, loading, error, refetch } = useQuery<{ contract: Contract }>(graphql.query.CONTRACT, {
		variables: { code },
		fetchPolicy: isNewContract ? "cache-only" : undefined, // prevent query if new
	})

	// Mutations
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [updateContract, mutUpdateContract] = useMutation(isNewContract ? graphql.mutation.CREATE_CONTRACT : graphql.mutation.UPDATE_CONTRACT)
	const [archiveContract, mutArchiveContract] = useMutation<{ contractArchive: Contract }>(graphql.mutation.ARCHIVE_CONTRACT)
	const [unarchiveContract, mutUnarchiveContract] = useMutation<{ contractUnarchive: Contract }>(graphql.mutation.UNARCHIVE_CONTRACT)

	const toggleArchive = () => {
		if (!contract) return

		if (contract.archived) {
			unarchiveContract({
				variables: { id: contract.id },
				update: (cache: any) => invalidateListQueries(cache, "contracts"),
			})
		} else {
			archiveContract({
				variables: { id: contract.id },
				update: (cache: any) => invalidateListQueries(cache, "contracts"),
			})
		}
	}

	// Form submission
	const [description, setDescription] = React.useState("")
	const [dateSigned, setDateSigned] = React.useState(new Date())

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()
	const onSubmit = handleSubmit(({ name, supplierName }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		const input = { name, description, supplierName, dateSigned }

		if (isNewContract) {
			updateContract({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "contracts"),
			})
		} else if (contract) {
			promiseTimeout(updateContract({ variables: { id: contract.id, input } })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	// On load contract
	React.useEffect(() => {
		if (!data || !data.contract) return
		setContract(data.contract)
	}, [data, loading, error])
	React.useEffect(() => {
		if (activeKey != "#details") return
		if (!contract) return
		setValue("name", contract.name)
		setValue("supplierName", contract.supplierName)
		setDescription(contract.description)
		setDateSigned(new Date(contract.dateSigned))
	}, [contract, activeKey])

	// On mutation (update/create contract)
	React.useEffect(() => {
		if (!mutUpdateContract.data) return

		if (isNewContract) {
			if (mutUpdateContract.data.contractCreate) {
				history.push(`/portal/contract/${mutUpdateContract.data.contractCreate.code}`)
			}
			return
		}

		if (!mutUpdateContract.data.contractUpdate) return

		setContract(mutUpdateContract.data.contractUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateContract.data, mutUpdateContract.loading])

	React.useEffect(() => {
		if (!mutArchiveContract.data?.contractArchive) return
		setContract(mutArchiveContract.data.contractArchive)
	}, [mutArchiveContract])
	React.useEffect(() => {
		if (!mutUnarchiveContract.data?.contractUnarchive) return
		setContract(mutUnarchiveContract.data.contractUnarchive)
	}, [mutUnarchiveContract])

	if (!isNewContract && !contract) {
		return <LoadingSimple />
	}

	const editForm = (
		<form onSubmit={onSubmit}>
			{changeSuccess && (
				<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
					Contract has been updated.
				</Notification>
			)}

			{mutUpdateContract.error && <ErrorNotification message={mutUpdateContract.error.message} />}

			<FormControl label="Name" error={errors.name ? errors.name.message : ""} positive="">
				<Input name="name" inputRef={register({ required: "Required" })} />
			</FormControl>

			<FormControl label="Description" error="" positive="">
				<Textarea
					name="description"
					value={description}
					onChange={e => setDescription(e.currentTarget.value)}
					overrides={{
						Input: {
							style: {
								resize: "vertical",
								height: "170px",
							},
						},
					}}
				/>
			</FormControl>

			<FormControl label="Supplier Name" error={errors.supplierName ? errors.supplierName.message : ""} positive="">
				<Input name="supplierName" inputRef={register({ required: "Required" })} />
			</FormControl>

			<FormControl label="Date Signed" caption="YYYY/MM/DD" error="" positive="">
				<div
					style={{
						width: "160px",
						marginRight: "10px",
					}}
				>
					<Datepicker value={dateSigned} onChange={({ date }) => setDateSigned(date as Date)} />
				</div>
			</FormControl>

			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/contracts")}>
					Cancel
				</Button>
				{contract && !isNewContract ? (
					<Spaced>
						<Button
							type="button"
							kind="secondary"
							isLoading={mutArchiveContract.loading || mutUnarchiveContract.loading}
							onClick={toggleArchive}
							startEnhancer={<FontAwesomeIcon icon={["fas", contract.archived ? "undo" : "archive"]} size="lg" />}
						>
							{contract.archived ? "Unarchive" : "Archive"}
						</Button>
						<Button isLoading={mutUpdateContract.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
							{timedOut ? "Timed out... Try again" : "Save"}
						</Button>
					</Spaced>
				) : (
					<Button isLoading={mutUpdateContract.loading} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
						Create Livestock Specification
					</Button>
				)}
			</Spread>
		</form>
	)

	return (
		<CenteredPage>
			<Spaced>
				<FontAwesomeIcon icon={["fal", "file-contract"]} size="3x" />
				<H1>{isNewContract || !contract ? "New Livestock Specification" : contract.name}</H1>
			</Spaced>

			{isNewContract ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/contract/${code}${activeKey}`)
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

					{!isNewContract && contract && (
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
								contractID={contract.id}
								itemName="product"
								query={graphql.query.PRODUCTS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_PRODUCT}
								extraFilterOptions={[
									{ label: "Not in Carton", id: FilterOption.ProductWithoutCarton },
									{ label: "Not in Order", id: FilterOption.ProductWithoutOrder },
									{ label: "No SKU", id: FilterOption.ProductWithoutSKU },
								]}
								itemLinks={["order", "sku", "distributor", "carton"]}
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

export default ContractEdit
