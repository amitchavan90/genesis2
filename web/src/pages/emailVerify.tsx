import * as React from "react"
import { Card } from "baseui/card"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Spaced } from "../components/spaced"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { ErrorBox } from "../components/errorBox"
import { useStyletron } from "baseui"
import { AuthContainer } from "../controllers/auth"
import { Redirect, RouteComponentProps } from "react-router-dom"

const EmailVerify = (props: RouteComponentProps<{ code?: string }>) => {
	const code = props.match.params.code

	const [css, theme] = useStyletron()
	const cardContainer: string = css({
		height: "100vh",
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
		backgroundColor: theme.colors.colorPrimary,
	})

	const { userErrors, loading, loggedIn, verify } = AuthContainer.useContainer()

	const [verifyCode, setVerifyCode] = React.useState<string>(code ? code : "")
	const [inputError, setInputError] = React.useState<string>("")

	React.useEffect(() => {
		if (code) verify(code)
	}, [code])

	if (loggedIn) {
		return <Redirect to={"/portal"} />
	}

	return (
		<div className={cardContainer}>
			<Card overrides={{ Root: { style: { flexGrow: 0.3 } } }}>
				<div>
					<ErrorBox userErrors={userErrors} />
					{code === undefined && (
						<FormControl label="Verification Code" error={inputError}>
							<Input
								error={!!inputError}
								positive={false}
								value={verifyCode}
								onChange={e => setVerifyCode(e.currentTarget.value)}
								placeholder={"Enter your verification code"}
							/>
						</FormControl>
					)}
					<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
						{loading && <Spinner />}
						<Button onClick={() => verify(verifyCode)}>Verify</Button>
					</Spaced>
				</div>
			</Card>
		</div>
	)
}

export default EmailVerify
