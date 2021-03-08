import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../../graphql"
import { User, UserLoyaltyActivity, Carton } from "../../../types/types"
import { AffiliateOrgs } from "../../../types/enums"
import { useForm } from "react-hook-form"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { H1 } from "baseui/typography"
import { Button } from "baseui/button"
import { Notification } from "baseui/notification"
import { Select, Value } from "baseui/select"
import { Spaced } from "../../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { LoadingSimple } from "../../../components/loading"
import { UserContainer } from "../../../controllers/user"
import { Block } from "baseui/block"
import { ErrorNotification } from "../../../components/errorBox"
import { Spread } from "../../../components/spread"
import { CenteredPage } from "../../../components/common"
import { ItemSelectList } from "../../../components/itemSelectList"
import { Tabs, Tab } from "baseui/tabs"
import { paddingZero } from "../../../themeOverrides"
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic"
import { SmallItemLink } from "../../../components/smallItemLink"
import { invalidateListQueries } from "../../../apollo"
import { HashButton } from "../../../components/hashButton"
import { promiseTimeout, TIMED_OUT } from "../../../helpers/timeout"

type FormData = {
	email: string
	firstName: string
	lastName: string
	password?: string
	mobilePhone: string
}

export const UserEdit = (props: RouteComponentProps<{ email: string }>) => <UserEditComponent id={props.match.params.email} hash={props.location.hash} />
export const ConsumerEdit = (props: RouteComponentProps<{ wechatID: string }>) => (
	<UserEditComponent id={props.match.params.wechatID} hash={props.location.hash} consumer />
)

const UserEditComponent = (props: { id?: string; hash: string; consumer?: boolean }) => {
	const { id, hash, consumer } = props
	const isNewUser = id === "new"

	const history = useHistory()

	const { user: me } = UserContainer.useContainer()

	// Get User
	const [user, setUser] = React.useState<User>()
	const { data, loading, error } = useQuery<{ user: User }>(consumer ? graphql.query.CONSUMER : graphql.query.USER, {
		variables: consumer ? { wechatID: id } : { email: id },
		fetchPolicy: isNewUser ? "cache-only" : undefined, // prevent query if new
	})

	// Form submission
	const [updateUser, mutUpdateUser] = useMutation(isNewUser ? graphql.mutation.CREATE_USER : graphql.mutation.UPDATE_USER)
	const [timedOut, setTimedOut] = React.useState(false)
	const [changeSuccess, setChangeSuccess] = React.useState(false)
	const [role, setRole] = React.useState<Value>()
	const [affiliateOrg, setAffiliateOrg] = React.useState<Value>()
	const [changePassword, setChangePassword] = React.useState(false)
	const [archiveUser, mutArchiveUser] = useMutation<{ userArchive: User }>(graphql.mutation.ARCHIVE_USER)
	const [unarchiveUser, mutUnarchiveUser] = useMutation<{ userUnarchive: User }>(graphql.mutation.UNARCHIVE_USER)

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()
	const onSubmit = handleSubmit(({ email, firstName, lastName, password, mobilePhone }) => {
		setChangeSuccess(false)
		setTimedOut(false)

		const input = {
			email,
			firstName,
			lastName,
			affiliateOrg: affiliateOrg && affiliateOrg.length > 0 ? affiliateOrg[0].id : undefined,
			roleID: role && role.length > 0 ? role[0].id : undefined,
			password,
			mobilePhone,
		}

		if (isNewUser) {
			updateUser({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "users"),
			})
		} else if (user) {
			promiseTimeout(updateUser({ variables: { id: user.id, input } })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	const toggleArchive = () => {
		if (!user) return

		if (user.archived) {
			unarchiveUser({
				variables: { id: user.id },
				update: (cache: any) => invalidateListQueries(cache, "users"),
			})
		} else {
			archiveUser({
				variables: { id: user.id },
				update: (cache: any) => invalidateListQueries(cache, "users"),
			})
		}
	}

	// On load user
	React.useEffect(() => {
		if (!data || !data.user) return
		setUser(data.user)
	}, [data, loading, error])
	React.useEffect(() => {
		if (!user) return
		setValue("email", user.email)
		setValue("firstName", user.firstName)
		setValue("lastName", user.lastName)
		setValue("mobilePhone", user.mobilePhone)

		if (user.affiliateOrg) {
			const v = { id: user.affiliateOrg }
			if (AffiliateOrgs.findIndex(o => o.id == user.affiliateOrg) == -1) AffiliateOrgs.push(v)
			setAffiliateOrg([v])
		}

		if (user.role) setRole([{ id: user.role.id, label: user.role.name }])
	}, [user])

	// On mutation (update/create user)
	React.useEffect(() => {
		if (!mutUpdateUser.data) return

		if (isNewUser) {
			if (mutUpdateUser.data.userCreate) {
				history.push(`/portal/user/${mutUpdateUser.data.userCreate.email}`)
			}
			return
		}

		if (!mutUpdateUser.data.userUpdate) return

		setUser(mutUpdateUser.data.userUpdate)
		setChangeSuccess(true)
		setTimedOut(false)
	}, [mutUpdateUser.data, mutUpdateUser.loading])

	// On hook-form update (register drop-down components)
	React.useEffect(() => {
		register({ name: "affiliateOrg" })
		register({ name: "role" })
	}, [register])

	React.useEffect(() => {
		if (!mutArchiveUser.data?.userArchive) return
		setUser(mutArchiveUser.data.userArchive)
	}, [mutArchiveUser])

	React.useEffect(() => {
		if (!mutUnarchiveUser.data?.userUnarchive) return
		setUser(mutUnarchiveUser.data.userUnarchive)
	}, [mutUnarchiveUser])

	// Styling
	const [css, theme] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "10px",
	})

	if (!isNewUser && !user) {
		return <LoadingSimple />
	}

	const title = (() => {
		if (isNewUser || !user) return "New User"
		if (consumer) return user.wechatID
		return `${user.firstName} ${user.lastName}`
	})()

	const canEdit = isNewUser || (me && user && me.role.tier < user.role.tier)

	const editForm = (
		<form onSubmit={onSubmit}>
			{changeSuccess && (
				<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
					The user has been updated.
				</Notification>
			)}

			{mutUpdateUser.error && <ErrorNotification message={mutUpdateUser.error.message} />}

			<FormControl label="Email" error={errors.email ? errors.email.message : ""} positive="">
				<Input
					name="email"
					type="email"
					inputRef={register({
						required: "Required",
						pattern: {
							value: /.+\@.+\..+/,
							message: "Invalid email address",
						},
					})}
					disabled={!canEdit}
				/>
			</FormControl>
			<FormControl label="First Name" error={errors.firstName ? errors.firstName.message : ""} positive="">
				<Input name="firstName" inputRef={register({ required: "Required" })} disabled={!canEdit} />
			</FormControl>
			<FormControl label="Last Name" error={errors.lastName ? errors.lastName.message : ""} positive="">
				<Input name="lastName" inputRef={register({ required: "Required" })} disabled={!canEdit} />
			</FormControl>
			<FormControl label="Mobile" error={errors.mobilePhone?.message}>
				<Input name="mobilePhone" inputRef={register({ required: "Required" })} />
			</FormControl>

			<FormControl label="Role">
				<ItemSelectList itemName="role" value={role} setValue={setRole} query={graphql.query.ROLES_LIMITED} identifier="name" disableSearch limit={0} />
			</FormControl>

			<FormControl label="Affiliate Organization" error="" positive="" caption="(Type a name to create a new organization)">
				<Select
					creatable
					options={AffiliateOrgs}
					labelKey="id"
					valueKey="id"
					value={affiliateOrg}
					onChange={({ value }) => setAffiliateOrg(value)}
					disabled={!canEdit}
				/>
			</FormControl>

			<FormControl label="Password" error={errors.password ? errors.password.message : ""} positive="">
				{changePassword ? (
					<Block display="flex">
						<Input name="password" type="password" inputRef={register({ required: { value: changePassword, message: "Required" } })} disabled={!canEdit} />
						<Button type="button" onClick={() => setChangePassword(false)}>
							Cancel
						</Button>
					</Block>
				) : (
					<Button type="button" disabled={!canEdit} onClick={() => setChangePassword(!changePassword)}>
						{isNewUser ? "Set Password" : "Change Password"}
					</Button>
				)}
			</FormControl>

			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/users")}>
					Cancel
				</Button>
				{user && !isNewUser ? (
					<Spaced>
						<Button
							type="button"
							kind="secondary"
							isLoading={mutArchiveUser.loading || mutUnarchiveUser.loading}
							onClick={toggleArchive}
							startEnhancer={<FontAwesomeIcon icon={["fas", user.archived ? "undo" : "archive"]} size="lg" />}
						>
							{user.archived ? "Unarchive" : "Archive"}
						</Button>

						<Button isLoading={mutUpdateUser.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
							{isNewUser ? "Create User" : timedOut ? "Timed out... Try again" : "Save"}
						</Button>
					</Spaced>
				) : (
					<Button isLoading={mutUpdateUser.loading && !timedOut}>Create User</Button>
				)}
			</Spread>
		</form>
	)

	return (
		<CenteredPage>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={consumer ? ["fal", "smile"] : ["fas", "user"]} size="3x" />
					<H1>{title}</H1>
				</Spaced>

				{consumer && user ? (
					<Button kind="secondary" type="button">
						<FontAwesomeIcon icon={["fas", "star"]} />
						<div className={marginLeftStyle}>{`${user.loyaltyPoints} Loyalty Points`}</div>
					</Button>
				) : (
					<div />
				)}
			</Spread>

			{consumer && user ? <ConsumerPage editForm={editForm} hash={hash} user={user} /> : editForm}
		</CenteredPage>
	)
}

const ConsumerPage = (props: { editForm: JSX.Element; hash: string; user: User }) => {
	const { editForm, hash, user } = props

	const history = useHistory()
	const [activeKey, setActiveKey] = React.useState(hash || "#details")

	const [loyaltyActivity, setLoyaltyActivity] = React.useState<UserLoyaltyActivity[]>([])
	const { data, loading, error } = useQuery<{ getLoyaltyActivity: UserLoyaltyActivity[] }>(graphql.query.GET_LOYALTY_ACTIVITY, {
		variables: { userID: user.id },
	})
	React.useEffect(() => {
		if (!data || !data.getLoyaltyActivity) return
		setLoyaltyActivity(data.getLoyaltyActivity)
	}, [data, loading, error])

	return (
		<Tabs
			onChange={({ activeKey }) => {
				setActiveKey(activeKey.toString())
				history.push(`/portal/consumer/${user.wechatID}${activeKey}`)
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

			<Tab
				key="#loyalty"
				title={
					<Spaced>
						<FontAwesomeIcon icon={["fal", "star"]} />
						<div>Loyalty Activity</div>
					</Spaced>
				}
			>
				<TableBuilder
					data={loyaltyActivity}
					overrides={{
						TableBodyRow: {
							style: ({ $theme }) => ({
								":hover": {
									backgroundColor: $theme.colors.colorPrimary,
									color: "white",
								},
							}),
						},
					}}
				>
					<TableBuilderColumn
						key={`column-points`}
						header={"Points"}
						overrides={{
							TableBodyCell: {
								style: {
									color: "inherit",
								},
							},
						}}
					>
						{item => (
							<div>
								{item.amount}
								{item.bonus > 0 && ` (${item.bonus} bonus)`}
							</div>
						)}
					</TableBuilderColumn>

					<TableBuilderColumn
						key={`column-message`}
						header={"Message"}
						overrides={{
							TableBodyCell: {
								style: {
									color: "inherit",
								},
							},
						}}
					>
						{item => <div>{item.message}</div>}
					</TableBuilderColumn>

					<TableBuilderColumn
						key={`column-product`}
						header={"Product"}
						overrides={{
							TableBodyCell: {
								style: {
									color: "inherit",
								},
							},
						}}
					>
						{item => (item.product ? <SmallItemLink code={item.product.code} link="/portal/product" itemName="product" /> : <div />)}
					</TableBuilderColumn>

					<TableBuilderColumn
						key="column-date"
						header="Date/Time"
						overrides={{
							TableBodyCell: {
								style: {
									color: "inherit",
								},
							},
						}}
					>
						{item => <div>{new Date(item.createdAt).toLocaleString()}</div>}
					</TableBuilderColumn>

					<TableBuilderColumn
						key="column-hash"
						header="Hash"
						overrides={{
							TableBodyCell: {
								style: {
									color: "inherit",
								},
							},
						}}
					>
						{item => (item.transactionHash ? <HashButton hash={item.transactionHash} /> : <></>)}
					</TableBuilderColumn>
				</TableBuilder>
			</Tab>
		</Tabs>
	)
}
