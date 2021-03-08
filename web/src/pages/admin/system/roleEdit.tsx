import * as React from "react"
import { useMutation, useQuery } from "@apollo/react-hooks"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useStyletron } from "baseui"
import { Button } from "baseui/button"
import { Checkbox } from "baseui/checkbox"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Notification } from "baseui/notification"
import { Spinner } from "baseui/spinner"
import { Caption1, H1 } from "baseui/typography"
import { useForm } from "react-hook-form"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { $enum } from "ts-enum-util"
import { invalidateListQueries } from "../../../apollo"
import { CenteredPage } from "../../../components/common"
import { ErrorNotification } from "../../../components/errorBox"
import { LoadingSimple } from "../../../components/loading"
import { Spaced } from "../../../components/spaced"
import { Spread } from "../../../components/spread"
import { UserContainer } from "../../../controllers/user"
import { graphql } from "../../../graphql"
import { FilterOption, Perm } from "../../../types/enums"
import { Role, TrackAction } from "../../../types/types"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"

type FormData = {
	name: string
}

type PermItem = {
	checked: boolean
	perm: Perm
}

const PermGroups: string[] = [
	"User",
	"Organisation",
	"Role",
	"SKU",
	"Container",
	"Pallet",
	"Carton",
	"Product",
	"Order",
	"TrackAction",
	"Contract",
	"Distributor",
	"Activity",
	"Other",
]

const RoleEdit = (props: RouteComponentProps<{ name: string }>) => {
	const name = props.match.params.name
	const isNewRole = name === "new"

	// list of permisions that are required for track actions
	const requiredPerms = ["CartonList", "CartonRead", "CartonUpdate", "ProductList", "ProductRead", "ProductUpdate", "TrackActionList", "TrackActionRead"]

	const requiredGroups = ["Carton", "Product", "TrackAction"]

	const history = useHistory()

	const { user: me } = UserContainer.useContainer()

	// Setup Perm groups
	const getPermGroup = (perm: Perm) => {
		const g = PermGroups.find(p => perm.toString().indexOf(p) == 0)
		if (!g) return PermGroups[PermGroups.length - 1] // Other
		return g
	}
	const getPermGroups = () => {
		let result = new Map<string, PermItem[]>()

		$enum(Perm).forEach(p => {
			let group = getPermGroup(p)
			if (!result.has(group)) result.set(group, [])
			result.get(group)?.push({ checked: false, perm: p })
		})

		return result
	}

	// Get Role
	const [role, setRole] = React.useState<Role>()
	const { data, loading, error } = useQuery<{ role: Role }>(graphql.query.ROLE, {
		variables: { name },
		fetchPolicy: isNewRole ? "cache-only" : undefined, // prevent query if new
	})

	// Get Track Actions
	const queryTrackActions = useQuery<{ trackActions: { trackActions: TrackAction[]; total: number } }>(graphql.query.TRACK_ACTIONS, {
		variables: {
			search: { filter: FilterOption.Active },
			limit: 100,
			offset: 0,
		},
	})

	// Form submission
	const [updateRole, mutUpdateRole] = useMutation(isNewRole ? graphql.mutation.CREATE_ROLE : graphql.mutation.UPDATE_ROLE)
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [perms, setPerms] = React.useState<Map<string, PermItem[]>>(getPermGroups())
	const [trackActionIDs, setTrackActionIDs] = React.useState<string[]>([])

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()
	const onSubmit = handleSubmit(({ name }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		let permissions: Perm[] = []
		perms?.forEach(g => g.forEach(p => (p.checked || (trackActionSelected && requiredPerms.includes(p.perm))) && permissions.push(p.perm)))

		const input = { name, permissions, trackActionIDs }

		if (isNewRole) {
			updateRole({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "roles"),
			})
		} else if (role) {
			promiseTimeout(updateRole({ variables: { id: role.id, input } })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	// Archiving
	const [archiveRole, mutArchiveRole] = useMutation<{ roleArchive: Role }>(graphql.mutation.ARCHIVE_ROLE)
	const [unarchiveRole, mutUnarchiveRole] = useMutation<{ roleUnarchive: Role }>(graphql.mutation.UNARCHIVE_ROLE)

	const toggleArchive = () => {
		if (!role) return
		if (role.archived) {
			unarchiveRole({
				variables: { id: role.id },
				update: (cache: any) => invalidateListQueries(cache, "roles"),
			})
		} else {
			archiveRole({
				variables: { id: role.id },
				update: (cache: any) => invalidateListQueries(cache, "roles"),
			})
		}
	}

	React.useEffect(() => {
		if (!mutArchiveRole.data?.roleArchive) return
		setRole(mutArchiveRole.data.roleArchive)
	}, [mutArchiveRole])
	React.useEffect(() => {
		if (!mutUnarchiveRole.data?.roleUnarchive) return
		setRole(mutUnarchiveRole.data.roleUnarchive)
	}, [mutUnarchiveRole])

	// On load role
	React.useEffect(() => {
		if (!data || !data.role) return
		setRole(data.role)
	}, [data, loading, error])
	React.useEffect(() => {
		if (!role) return
		setValue("name", role.name)

		if (role.permissions) {
			const loadedPerms = getPermGroups()
			loadedPerms.forEach(g => {
				g.forEach(p => {
					if (role.permissions.indexOf(p.perm) != -1) p.checked = true
				})
			})
			setPerms(loadedPerms)
		}
		if (role.trackActions) setTrackActionIDs(role.trackActions.map(a => a.id))
	}, [role])

	// On mutation (update/create role)
	React.useEffect(() => {
		if (!mutUpdateRole.data) return

		if (isNewRole) {
			if (mutUpdateRole.data.roleCreate) {
				history.push(`/portal/role/${mutUpdateRole.data.roleCreate.name}`)
			}
			return
		}

		if (!mutUpdateRole.data.roleUpdate) return

		setRole(mutUpdateRole.data.roleUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateRole.data, mutUpdateRole.loading])

	// Styling
	const [css, theme] = useStyletron()
	const permSectionStyle = css({
		display: "grid",
		gridTemplateColumns: "32% 32% 32%",
		columnGap: "2%",
	})
	const permGroupStyle = css({
		marginBottom: "10px",
	})
	const permStyle = css({
		marginLeft: "10px",
	})
	const trackActionStyle = css({
		display: "flex",
		flexDirection: "column",
		lineHeight: "16px",
	})
	const trackActionAltStyle = css({
		color: "grey",
		fontSize: "12px",
	})

	if (!isNewRole && !role) {
		return <LoadingSimple />
	}

	const canEdit = isNewRole || (me && role && me.role.tier < role.tier)

	const trackActionSelected = trackActionIDs.length > 0

	return (
		<form onSubmit={onSubmit}>
			<CenteredPage>
				<Spaced>
					<FontAwesomeIcon icon={["fas", "user-tag"]} size="3x" />
					<H1>{isNewRole || !role ? "New Role" : role.name}</H1>
				</Spaced>

				{changeSuccess && (
					<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
						The role has been updated.
					</Notification>
				)}

				{mutUpdateRole.error && <ErrorNotification message={mutUpdateRole.error.message} />}

				<FormControl label="Name" error={errors.name ? errors.name.message : ""} positive="">
					<Input name="name" inputRef={register({ required: "Required" })} disabled={!canEdit} />
				</FormControl>

				{perms && (
					<FormControl label="Permissions" positive="" error="">
						<div className={permSectionStyle}>
							{Array.from(perms.keys()).map(group => {
								// Perm Group
								const permGroup = perms.get(group)
								if (!permGroup) return <></>

								const allChecked = permGroup.every(p => p.checked)

								const isIndeterminate = !allChecked && permGroup.some(p => p.checked)
								return (
									<div key={`permGroup-${group}`} className={permGroupStyle}>
										<Checkbox
											key={`RoleEditPermGroup${group}`}
											checked={allChecked}
											isIndeterminate={isIndeterminate || (!allChecked && trackActionSelected && requiredGroups.includes(group))}
											onChange={_ => {
												setPerms(
													new Map(
														perms.set(
															group,
															permGroup.map(p => {
																return { checked: !allChecked, perm: p.perm }
															}),
														),
													),
												)
											}}
											disabled={!canEdit}
										>
											{group}
										</Checkbox>

										{/* Perm */}
										<div className={permStyle}>
											{permGroup.map(p => {
												return (
													<Checkbox
														key={`RoleEditPerm${p.perm}`}
														checked={p.checked || (trackActionSelected && requiredPerms.includes(p.perm))}
														onChange={e => {
															const index = permGroup.indexOf(p)
															setPerms(
																new Map(
																	perms.set(group, [
																		...permGroup.slice(0, index),
																		{ checked: e.currentTarget.checked || (trackActionSelected && requiredPerms.includes(p.perm)), perm: p.perm },
																		...permGroup.slice(index + 1),
																	]),
																),
															)
														}}
														disabled={!canEdit}
													>
														{p.perm
															.toString()
															.replace(group, "")
															.replace(/([A-Z])/g, " $1")}
													</Checkbox>
												)
											})}
										</div>
									</div>
								)
							})}
						</div>
					</FormControl>
				)}

				<FormControl
					label="Track Actions (If a track action is selected, some permisions will be automatically ticked)"
					overrides={{ Label: { style: { marginBottom: 0 } } }}
					positive=""
					error=""
				>
					<div>
						<Caption1
							overrides={{
								Block: { style: { marginTop: 0 } },
							}}
						>
							The actions this role can do when updating product tracking.
						</Caption1>
						{queryTrackActions.data?.trackActions && (
							<div className={permSectionStyle}>
								{queryTrackActions.data.trackActions.trackActions.map(action => (
									<Checkbox
										key={`RoleEditTrackAction${action.id}`}
										overrides={{
											Root: {
												style: {
													alignItems: "center",
													marginTop: "5px",
												},
											},
										}}
										checked={trackActionIDs.indexOf(action.id) != -1}
										onChange={e => {
											// force
											const index = trackActionIDs.indexOf(action.id)
											if (index != -1) setTrackActionIDs([...trackActionIDs.slice(0, index), ...trackActionIDs.slice(index + 1)])
											else setTrackActionIDs([...trackActionIDs, action.id])
										}}
									>
										<div className={trackActionStyle}>
											<div>{action.name}</div>
											<div className={trackActionAltStyle}>{action.nameChinese}</div>
										</div>
									</Checkbox>
								))}
							</div>
						)}
						{queryTrackActions.loading && <Spinner />}
					</div>
				</FormControl>

				<Spread>
					<Button type="button" kind="secondary" onClick={() => history.push("/portal/roles")}>
						Cancel
					</Button>
					{role && !isNewRole ? (
						<Spaced>
							<Button
								type="button"
								kind="secondary"
								isLoading={mutArchiveRole.loading || mutUnarchiveRole.loading}
								onClick={toggleArchive}
								startEnhancer={<FontAwesomeIcon icon={["fas", role.archived ? "undo" : "archive"]} size="lg" />}
							>
								{role.archived ? "Unarchive" : "Archive"}
							</Button>
							<Button isLoading={mutUpdateRole.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
								{timedOut ? "Timed out... Try again" : "Save"}
							</Button>
						</Spaced>
					) : (
						<Button isLoading={mutUpdateRole.loading}>Create Role</Button>
					)}
				</Spread>
			</CenteredPage>
		</form>
	)
}

export default RoleEdit
