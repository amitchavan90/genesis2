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
import { Spread } from "../../../components/spread"
import { CenteredPage } from "../../../components/common"
import { UserTask} from "../../../types/types"
import { Tabs, Tab } from "baseui/tabs"
import { Modal, ModalButton, ModalFooter, ModalHeader } from "baseui/modal"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"
import { paddingZero} from "../../../themeOverrides"
import { Notification } from "baseui/notification"
import {
	Label1,
	Label2,
	Paragraph2,
  } from 'baseui/typography';





const taskApprove = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewTask = "new" ===code 
	
	const history = useHistory()
	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")
	const [showPreviewModal, setShowPreviewModal] = React.useState<boolean>()

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
	
	
	// Get TASK
	const [usertask, setUserTask] = React.useState<UserTask>()
	console.log("after set user task ----------------------->",usertask);
	const { data, loading, error, refetch } = useQuery<{ userTask: UserTask}>(graphql.query.USER_TASK, {
		variables: { code }
	})
	// Mutations
	const [approveTask,mutApproveTask] = useMutation(graphql.mutation.USER_TASK_APPROVE)
	
	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	// const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [showSuccessModal, setShowSuccessModal] = React.useState(false)
	const [showConfirmationModal, setShowConfirmationModal] = React.useState(false)
	
	const breakLine = <div className={breakLineStyle} />

	const { register, setValue, handleSubmit, errors, getValues } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({}) => {
		console.log("on submit");
		setTimedOut(false)
		setShowConfirmationModal(true)
		
	})
	 
	const updateUserTask = () => {
		 if (usertask) {
			promiseTimeout(approveTask({ variables: { id: usertask.id} })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
			setShowConfirmationModal(false)
			setShowSuccessModal(true)
			setTimedOut(false)
		}
		
	}

	React.useEffect(() => {
		if (!data || !data.userTask) 
		return
		setUserTask(data.userTask)
	}, [data, loading, error])

	React.useEffect(() => {
		if (!usertask) return
	}, [usertask])

	
	
	const editForm = (
		<form onSubmit={onSubmit}>
			{usertask?.isComplete ?
				<Notification closeable kind="positive">
					UserTask  alreday Approved
				</Notification>:<div></div>
			}
			<Label2>Code :</Label2>
			<Paragraph2>{usertask?.code}</Paragraph2>
			<Label2>Status :</Label2>
			<Paragraph2>{usertask?.status}</Paragraph2>
			<Label1>User Details:</Label1>
			{breakLine}
			<Label2>Name :</Label2>
			<Paragraph2>{usertask?.user.firstName} {usertask?.user.lastName}</Paragraph2>
			<Label2>Email :</Label2>
			<Paragraph2>{usertask?.user.email}</Paragraph2>
			<Label1>Task Details:</Label1>
			{breakLine}
			<Label2>Title :</Label2>
			<Paragraph2>{usertask?.task.title}</Paragraph2>
			<Label2>Description:</Label2>
			<Paragraph2>{usertask?.task.description}</Paragraph2>
			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/userTasks")}>
					Cancel
				</Button>
				{usertask?.isComplete?<div></div>:
				<Button isLoading={mutApproveTask.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
						{isNewTask ? "Create Task" : timedOut ? "Timed out... Try again" : "Approve"}
				</Button>}
			</Spread>
		</form>
	)
	
		// Successfully created order
		return (
			<CenteredPage>
		
				<Spaced>
					<FontAwesomeIcon icon={["fal", "file-contract"]} size="3x" />
					<H1>{isNewTask ? "New Task" : "Approve Task"}</H1>
				</Spaced>
				{isNewTask ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/sku/${code}${activeKey}`)
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
				</Tabs>
			)}
			
			<Modal isOpen={showPreviewModal} onClose={() => setShowPreviewModal(false)}>
				
			</Modal>
			{/* Confirmation Modal */}
			<Modal onClose={() => setShowConfirmationModal(false)} isOpen={showConfirmationModal}>
				<ModalHeader>
					<span>
						<FontAwesomeIcon icon={["fas", "pencil"]} />
						<span style={{ marginLeft: "10px" }}>Are you sure you want to Approve UsetTask ?</span>
					</span>
				</ModalHeader>
				<ModalFooter>
					<ModalButton onClick={() => history.push("/portal/userTasks")}>Cancel</ModalButton>
					<ModalButton onClick={() => updateUserTask()}>OK</ModalButton>
				</ModalFooter>
			</Modal>
			{/* Success Modal */}
			<Modal onClose={() => setShowSuccessModal(false)} isOpen={showSuccessModal}>
				<ModalHeader>
					<span>
						<FontAwesomeIcon icon={["fas", "check"]} />
						<span style={{ marginLeft: "10px" }}>Task Approve Successfully</span>
					</span>
				</ModalHeader>
				<ModalFooter>
					<ModalButton onClick={() => setShowSuccessModal(false)}>OK</ModalButton>
				</ModalFooter>
			</Modal>
			</CenteredPage>
		)
	}

export default taskApprove
