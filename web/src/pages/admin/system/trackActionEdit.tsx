import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { TrackAction } from "../../../types/types"
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
import { CenteredPage } from "../../../components/common"
import { Spread } from "../../../components/spread"
import { Checkbox } from "baseui/checkbox"
import { invalidateListQueries } from "../../../apollo"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"

type FormData = {
	name: string
	nameChinese: string
}

const TrackActionEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewTrackAction = code === "new"

	const history = useHistory()

	// Get TrackAction
	const [trackAction, setTrackAction] = React.useState<TrackAction>()
	const { data, loading, error } = useQuery<{ trackAction: TrackAction }>(graphql.query.TRACK_ACTION, {
		variables: { code },
		fetchPolicy: isNewTrackAction ? "cache-only" : undefined, // prevent query if new
	})

	// Mutations
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [updateTrackAction, mutUpdateTrackAction] = useMutation(isNewTrackAction ? graphql.mutation.CREATE_TRACK_ACTION : graphql.mutation.UPDATE_TRACK_ACTION)
	const [archiveTrackAction, mutArchiveTrackAction] = useMutation<{ trackActionArchive: TrackAction }>(graphql.mutation.ARCHIVE_TRACK_ACTION)
	const [unarchiveTrackAction, mutUnarchiveTrackAction] = useMutation<{ trackActionUnarchive: TrackAction }>(graphql.mutation.UNARCHIVE_TRACK_ACTION)

	const toggleArchive = () => {
		if (!trackAction) return

		if (trackAction.archived) {
			unarchiveTrackAction({
				variables: { id: trackAction.id },
				update: (cache: any) => invalidateListQueries(cache, "trackActions"),
			})
		} else {
			archiveTrackAction({
				variables: { id: trackAction.id },
				update: (cache: any) => invalidateListQueries(cache, "users"),
			})
		}
	}

	// Form submission
	const [requirePhotos, setRequirePhotos] = React.useState<boolean[]>([false, false])
	const [privateAction, setPrivateAction] = React.useState(false)
	const [blockchain, setBlockchain] = React.useState(true)
	const { register, setValue, handleSubmit, errors } = useForm<FormData>()
	const onSubmit = handleSubmit(({ name, nameChinese }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		const input = { name, nameChinese, private: privateAction, blockchain, requirePhotos }

		if (isNewTrackAction)
			updateTrackAction({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "trackActions"),
			})
		else if (trackAction) {
			promiseTimeout(updateTrackAction({ variables: { id: trackAction.id, input } })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	// On load trackAction
	React.useEffect(() => {
		if (!data || !data.trackAction) return
		setTrackAction(data.trackAction)
	}, [data, loading, error])
	React.useEffect(() => {
		if (!trackAction) return
		setValue("name", trackAction.name)
		setValue("nameChinese", trackAction.nameChinese)
		setPrivateAction(trackAction.private)
		setBlockchain(trackAction.blockchain)
		setRequirePhotos(trackAction.requirePhotos)
	}, [trackAction])

	// On mutation (update/create trackAction)
	React.useEffect(() => {
		if (!mutUpdateTrackAction.data) return

		if (isNewTrackAction) {
			if (mutUpdateTrackAction.data.trackActionCreate) {
				history.push(`/portal/trackAction/${mutUpdateTrackAction.data.trackActionCreate.code}`)
			}
			return
		}

		if (!mutUpdateTrackAction.data.trackActionUpdate) return

		setTrackAction(mutUpdateTrackAction.data.trackActionUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateTrackAction.data, mutUpdateTrackAction.loading])

	React.useEffect(() => {
		if (!mutArchiveTrackAction.data?.trackActionArchive) return
		setTrackAction(mutArchiveTrackAction.data.trackActionArchive)
	}, [mutArchiveTrackAction])
	React.useEffect(() => {
		if (!mutUnarchiveTrackAction.data?.trackActionUnarchive) return
		setTrackAction(mutUnarchiveTrackAction.data.trackActionUnarchive)
	}, [mutUnarchiveTrackAction])

	if (!isNewTrackAction && !trackAction) {
		return <LoadingSimple />
	}

	return (
		<form onSubmit={onSubmit}>
			<CenteredPage>
				<Spaced>
					<FontAwesomeIcon icon={["fas", "truck-moving"]} size="3x" />
					<H1>{isNewTrackAction || !trackAction ? "New Track Action" : trackAction.name}</H1>
				</Spaced>

				{changeSuccess && (
					<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
						Track Action has been updated.
					</Notification>
				)}

				{mutUpdateTrackAction.error && <ErrorNotification message={mutUpdateTrackAction.error.message} />}

				<FormControl label="Name" error={errors.name ? errors.name.message : ""} positive="">
					<Input name="name" inputRef={register({ required: "Required" })} />
				</FormControl>

				<FormControl label="Name (Chinese)" error={errors.nameChinese ? errors.nameChinese.message : ""} positive="">
					<Input name="nameChinese" inputRef={register} />
				</FormControl>

				<FormControl caption="Whether or not customers can see this action under Product Tracking Information.">
					<Checkbox checked={privateAction} onChange={e => setPrivateAction(e.currentTarget.checked)}>
						Private
					</Checkbox>
				</FormControl>

				<FormControl caption="Whether or not this action is commited to blockchain.">
					<Checkbox checked={blockchain} onChange={e => setBlockchain(e.currentTarget.checked)}>
						Blockchain
					</Checkbox>
				</FormControl>

				<FormControl caption="Whether or not this action requires a photo of the carton+label to be taken (for cartons only).">
					<Checkbox checked={requirePhotos[0]} onChange={e => setRequirePhotos([e.currentTarget.checked, requirePhotos[1]])}>
						Require Carton+Label photo
					</Checkbox>
				</FormControl>

				<FormControl caption="Whether or not this action requires a photo of the product+qr to be taken (for cartons only).">
					<Checkbox checked={requirePhotos[1]} onChange={e => setRequirePhotos([requirePhotos[0], e.currentTarget.checked])}>
						Require Product+QR photo
					</Checkbox>
				</FormControl>

				<Spread>
					<Button type="button" kind="secondary" onClick={() => history.push("/portal/trackActions")}>
						Cancel
					</Button>
					{trackAction && !isNewTrackAction ? (
						<Spaced>
							<Button
								type="button"
								kind="secondary"
								isLoading={mutArchiveTrackAction.loading || mutUnarchiveTrackAction.loading}
								onClick={toggleArchive}
								startEnhancer={<FontAwesomeIcon icon={["fas", trackAction.archived ? "undo" : "archive"]} size="lg" />}
								disabled={trackAction.system}
							>
								{trackAction.archived ? "Unarchive" : "Archive"}
							</Button>
							<Button isLoading={mutUpdateTrackAction.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
								{timedOut ? "Timed out... Try again" : "Save"}
							</Button>
						</Spaced>
					) : (
						<Button isLoading={mutUpdateTrackAction.loading} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
							Create Track Action
						</Button>
					)}
				</Spread>
			</CenteredPage>
		</form>
	)
}

export default TrackActionEdit
