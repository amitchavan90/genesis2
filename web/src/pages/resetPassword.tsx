import * as React from "react"
import { Card } from "baseui/card"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Spaced } from "../components/spaced"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { ErrorNotification } from "../components/errorBox"
import { useStyletron } from "baseui"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useMutation, useQuery } from "@apollo/react-hooks"
import { graphql } from "../graphql"
import { useForm } from "react-hook-form"
import { LoadingSimple } from "../components/loading"

const ResetPassword = (props: RouteComponentProps<{ token: string }>) => {
	const token = props.match.params.token

	const history = useHistory()

	const [css, theme] = useStyletron()
	const cardContainer: string = css({
		height: "100vh",
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
		backgroundColor: theme.colors.colorPrimary,
	})

	const { data, loading, error } = useQuery<{ verifyResetToken: boolean }>(graphql.query.VERIFY_RESET_TOKEN, { variables: { token } })
	const [validToken, setValidToken] = React.useState<boolean>()

	const [resetPassword, mutResetPassword] = useMutation<{ resetPassword: boolean }>(graphql.mutation.RESET_PASSWORD)
	const [changeSuccess, setChangeSuccess] = React.useState<boolean>()
	const { register, handleSubmit, errors } = useForm<{ password: string }>()

	const onSubmit = handleSubmit(({ password }) => {
		setChangeSuccess(false)
		resetPassword({ variables: { token, password } })
	})

	React.useEffect(() => {
		if (loading) return
		setValidToken(data ? data.verifyResetToken : false)
	}, [data, loading, error])

	React.useEffect(() => {
		if (!mutResetPassword.data) return
		setChangeSuccess(mutResetPassword.data.resetPassword)
	}, [mutResetPassword])

	if (validToken === undefined) {
		return <LoadingSimple />
	}

	return (
		<div className={cardContainer}>
			<Card overrides={{ Root: { style: { flexGrow: 0.3 } } }}>
				{changeSuccess && (
					<div>
						<div>Your password has been updated.</div>
						<Button onClick={() => history.push("/")}>Login</Button>
					</div>
				)}

				{mutResetPassword.error && <ErrorNotification message={mutResetPassword.error?.message} />}

				{!validToken && <ErrorNotification closeable={false} message={error ? error.message : "Invalid Token"} />}

				{validToken && !changeSuccess && (
					<form onSubmit={onSubmit}>
						<FormControl label="Password" error={errors.password?.message}>
							<Input
								name="password"
								type="password"
								inputRef={register({ required: "Required" })}
								disabled={mutResetPassword.loading}
								error={errors.password !== undefined}
							/>
						</FormControl>

						<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
							{mutResetPassword.loading && <Spinner />}
							<Button>Reset Password</Button>
						</Spaced>
					</form>
				)}
			</Card>
		</div>
	)
}
export default ResetPassword
