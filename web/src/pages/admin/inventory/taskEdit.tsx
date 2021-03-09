import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { useForm } from "react-hook-form"
import { Textarea } from "baseui/textarea"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { H1, H3 } from "baseui/typography"
import { Button } from "baseui/button"
import { Spaced } from "../../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { ErrorNotification } from "../../../components/errorBox"
import { Spread } from "../../../components/spread"
import { CenteredPage } from "../../../components/common"
import { invalidateListQueries } from "../../../apollo"
import { Checkbox } from "baseui/checkbox"
import { Task} from "../../../types/types"

type FormData = {
	title: string
	loyaltyPoints: number
    finishDate: Date
    maximumPeople: number
    skuID: string
}

const taskEdit = (props: RouteComponentProps<{ code: string }>) => {
	console.log("i am in task update")

	const isNewTask = "new"
	const history = useHistory()

	//style
	const [css, theme] = useStyletron()
	const breakLineStyle = css({
		border: "1px solid rgba(0, 0, 0, 0.1)",
		marginBottom: "10px",
	})
	
	// Get SKU
	const [task, setTask] = React.useState<Task>()
	// Mutations
	const [updateTask, mutUpdateTask] = useMutation(isNewTask ? graphql.mutation.CREATE_TASK : graphql.mutation.UPDATE_SKU)
	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [description, setDescription] = React.useState("")
	const [isTimeBound, setIsTimeBound] = React.useState<boolean>()
	const [isPeopleBound, setIsPeopleBound] = React.useState<boolean>()
	const [isProductRelevant, setIsProductRelevant] = React.useState<boolean>(false)

	const breakLine = <div className={breakLineStyle} />

	const { register, setValue, handleSubmit, errors, getValues } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ title, loyaltyPoints, finishDate, maximumPeople, skuID}) => {
		setTimedOut(false)
		const input = {
			title,
			description,
			isTimeBound,
			isPeopleBound,
			isProductRelevant,
			loyaltyPoints,
			finishDate,
			maximumPeople,
			skuID
		}

		if (isNewTask) {
			console.log("input---------------->",input);
			updateTask({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "tasks"),
			})
		} else if (1) {
			// promiseTimeout(updateTask({ variables: { id: sku.id, input } })).catch(reason => {
			// 	if (reason !== TIMED_OUT) return
			// 	setTimedOut(true)
			// })
		}
	})

	React.useEffect(() => {
		if (!task) return
		setValue("title",task.title)
		setValue("loyaltyPoints", task.loyaltyPoints)
		setValue("finishDate", task.finishDate)
		setValue("maximumPeople", task.maximumPeople)
		setValue("L280001", task.skuID)
	}, [task])
	
	const editForm = (
		<form onSubmit={onSubmit}>
			{/* {isNewTask && cloneFrom && <div>Cloned from {cloneFrom}</div>} */}

			{mutUpdateTask.error && <ErrorNotification message={mutUpdateTask.error.message} />}

			{/* {!isNewTask && (
				<FormControl label="Code" error={errors.code ? errors.code.message : ""} positive="" disabled>
					<Input name="code" inputRef={register} disabled />
				</FormControl>
			)} */}

			<FormControl label="Name" error={errors.title ? errors.title.message : ""} positive="">
				<Input name="title" inputRef={register({ required: "Required" })} />
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

			<FormControl label="LoyaltyPoints" error={errors.title ? errors.title.message : ""} positive="">
				<Input name="loyaltyPoints" type="number" inputRef={register}/>
			</FormControl>

			<FormControl caption="">
				<Checkbox checked={isTimeBound} onChange={e => setIsTimeBound(e.currentTarget.checked)}>
					Is Time Bound
				</Checkbox>
			</FormControl>

			<FormControl caption="">
				<Checkbox checked={isPeopleBound} onChange={e => setIsPeopleBound(e.currentTarget.checked)}>
					Is People Bound
				</Checkbox>
			</FormControl>
			<FormControl caption="">
				<Checkbox checked={isProductRelevant} onChange={e => setIsProductRelevant(e.currentTarget.checked)}>
					Is Product Relevant
				</Checkbox>
			</FormControl>
			{breakLine}
			<FormControl label="Maximum People" error={errors.maximumPeople ? errors.maximumPeople.message : ""} positive="">
				<Input name="maximumPeople" type="number" inputRef={register}/>
			</FormControl>
			<FormControl label="Finish Date" error={errors.finishDate ? errors.finishDate.message : ""} positive="">
				<Input name="finishDate" type="Date" inputRef={register} />
			</FormControl>
			{breakLine}
			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/skus")}>
					Cancel
				</Button>

				<Button isLoading={mutUpdateTask.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
						{isNewTask ? "Create Task" : timedOut ? "Timed out... Try again" : "Save"}
				</Button>
			</Spread>
		</form>
	)
	
		// Successfully created order
		return (
			<CenteredPage>
				<Spaced>
					<FontAwesomeIcon icon={["fal", "file-contract"]} size="3x" />
					<H1>{isNewTask ? "New Task" : "Edit Task"}</H1>
				</Spaced>
				{isNewTask ? (
				editForm
			) : (<div>view task</div>)}
			</CenteredPage>
		)
	}

export default taskEdit
