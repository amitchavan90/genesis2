import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { useForm } from "react-hook-form"
import { Textarea } from "baseui/textarea"
import { FormControl,FormControlOverrides} from "baseui/form-control"
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
import { Datepicker } from "baseui/datepicker"
import { Task, SubTask} from "../../../types/types"
import { ItemSelectList } from "../../../components/itemSelectList"
import { Select, Value } from "baseui/select"
import { IconName } from "@fortawesome/fontawesome-svg-core"


type FormData = {
	title: string
	loyaltyPoints: number
    maximumPeople: number
    skuID: string
	subtasks : SubTask[]
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
	const plusCircleMargin = css({
		marginRight: "10px",
	})

	const plusCircle = (
		<div className={plusCircleMargin}>
			<FontAwesomeIcon icon={["fal", "plus-circle"]} size="lg" />
		</div>
	)
	
	
	// Get SKU
	const [task, setTask] = React.useState<Task>()
	// Mutations
	const [updateTask, mutUpdateTask] = useMutation(isNewTask ? graphql.mutation.CREATE_TASK : graphql.mutation.UPDATE_SKU)
	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [description, setDescription] = React.useState("")
	const [isTimeBound, setIsTimeBound] = React.useState<boolean>(false)
	const [isPeopleBound, setIsPeopleBound] = React.useState<boolean>(false)
	const [isProductRelevant, setIsProductRelevant] = React.useState<boolean>(false)
	const [finishDate, setfinishDate] = React.useState(new Date())
	const [sku,setSKU] = React.useState<Value>()
	const [subTasksCount, setSubTasksCount] = React.useState(isNewTask? 0 : 25)
	const breakLine = <div className={breakLineStyle} />

	const { register, setValue, handleSubmit, errors, getValues } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ title, loyaltyPoints, maximumPeople, skuID, subtasks}) => {
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
			//skuID,
			subtasks,
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
		setValue("maximumPeople", task.maximumPeople)
		setfinishDate(new Date(task.finishDate))
		//setValue("L280001", task.skuID)
		setSubTasksCount(task.subtasks.length)
		task.subtasks.forEach((info, index) => {
			setValue(`subtasks[${index}].title`, info.title)
			setValue(`subtasks[${index}].description`, info.description)
		})
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
			{isPeopleBound?
			<FormControl label="Maximum People" error={errors.maximumPeople ? errors.maximumPeople.message : ""} positive="">
				<Input name="maximumPeople" type="number" inputRef={register}/>
			</FormControl>:<div></div>}
			{/* <FormControl label="Finish Date" error={errors.finishDate ? errors.finishDate.message : ""} positive="">
				<Input name="finishDate" type="Date" inputRef={register} />
			</FormControl> */}
			{isTimeBound?
			<FormControl label="Finish Date" caption="YYYY/MM/DD" error="" positive="">
				<div
					style={{
						width: "160px",
						marginRight: "10px",
					}}
				>
					<Datepicker value={finishDate} onChange={({ date }) => setfinishDate(date as Date)} />
				</div>
			</FormControl>:<div></div>}	
			{/* <FormControl label="Role">
				<ItemSelectList itemName="role" value={sku} setValue={setSKU} query={graphql.query.ROLES_LIMITED} identifier="name" disableSearch limit={0} />
			</FormControl> */}
			{breakLine}
			<FormControl label={`Sub Tasks (${subTasksCount}/25)`} error="" positive="">
				<div>
					{Array.from({ length: subTasksCount }).map((_, index) => (
						<InputPair
							key={`sku_info_${index}`}
							prefix="subtasks"
							index={index}
							icon="info-circle"
							titleInputRef={register({ required: "Required" })}
							contentInputRef={register({ required: "Required" })}
							titleError={errors.subtasks && errors.subtasks[index] && errors.subtasks[index].title?.message}
							contentError={errors.subtasks && errors.subtasks[index] && errors.subtasks[index].description?.message}
							onDelete={async () => {
								const subtasks = getValues({ nest: true }).subtasks
								const newInfo = [...subtasks.slice(0, index), ...subtasks.slice(index + 1)]
								subtasks.forEach((_, index) => {
									setValue(`subtasks[${index}].title`, index >= newInfo.length ? "" : newInfo[index].title)
									setValue(`subtasks[${index}].content`, index >= newInfo.length ? "" : newInfo[index].description)
								})
								setSubTasksCount(subTasksCount - 1)
							}}
						/>
					))}
					{subTasksCount < 25 && (
						<Button type="button" kind="secondary" onClick={() => setSubTasksCount(subTasksCount + 1)}>
							{plusCircle} Add SubTasks
						</Button>
					)}
				</div>
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

	interface InputPairProps {
		index: number
		prefix: string
		icon: IconName
		titleInputRef: any
		contentInputRef: any
		titleError?: string
		contentError?: string
		onDelete: () => void
		label1?: string
		label2?: string
	}
	
	export const InputPair = (props: InputPairProps) => {
		const { prefix, index } = props
		console.log("subtasks------------->",prefix,props.titleInputRef)
		const [css, theme] = useStyletron()
		const containerStyle = css({
			display: "flex",
			alignItems: "center",
			padding: "8px 0 8px 10px",
			backgroundColor: index % 2 == 1 ? "rgba(0, 0, 0, 0.015)" : "unset",
		})
		const iconStyle = css({
			marginRight: "10px",
		})
		const half = css({
			display: "flex",
			alignItems: "center",
			width: "50%",
		})
	
		return (
			<div className={containerStyle}>
				<div className={half}>
					<div className={iconStyle}>
						<FontAwesomeIcon icon={["fas", props.icon]} size="lg" />
					</div>
					<FormControl label={props.label1 || "Title"} overrides={InputParFormControlOverrides} error="" positive="">
						<Input name={`${prefix}[${index}].title`} inputRef={props.titleInputRef} error={props.titleError !== undefined} />
					</FormControl>
				</div>
				<div className={half}>
					<FormControl label={props.label2 || "Description"} overrides={InputParFormControlOverrides} error="" positive="">
						<Input name={`${prefix}[${index}].description`} inputRef={props.contentInputRef} error={props.contentError !== undefined} />
					</FormControl>
					<Button
						type="button"
						kind="minimal"
						onClick={props.onDelete}
						overrides={{
							BaseButton: {
								style: {
									color: "grey",
									":hover": {
										color: "#d63916",
									},
								},
							},
						}}
					>
						<FontAwesomeIcon icon={["fas", "trash"]} size="lg" />
					</Button>
				</div>
			</div>
		)
	}

	const InputParFormControlOverrides: FormControlOverrides = {
		Label: {
			style: {
				marginTop: 0,
				marginBottom: 0,
				marginRight: "10px",
				width: "unset",
			},
		},
		ControlContainer: {
			style: {
				marginBottom: 0,
				marginRight: "20px",
				width: "100%",
			},
		},
	}

export default taskEdit