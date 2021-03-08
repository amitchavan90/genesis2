import * as React from "react"
import { useStyletron } from "baseui"
import { graphql } from "../graphql"
import { useMutation } from "@apollo/react-hooks"
import { Spinner } from "baseui/spinner"
import { UserContainer } from "../controllers/user"

export const VerifyBanner = () => {
	const { user } = UserContainer.useContainer()

	const [resendVerify, mutResendVerify] = useMutation<{ resendEmailVerification: boolean }>(graphql.mutation.RESEND_EMAIL_VERIFICATION, {
		variables: { email: user?.email },
	})
	const [resendSuccess, setResendSuccess] = React.useState<boolean>()
	React.useEffect(() => {
		if (resendSuccess || !mutResendVerify.data) return
		setResendSuccess(mutResendVerify.data.resendEmailVerification)
	}, [mutResendVerify])

	const [css] = useStyletron()
	const containerStyle: string = css({
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		height: "30px",
		backgroundColor: "#4CAF50",
		color: "white",
		textAlign: "center",
		cursor: resendSuccess ? "default" : "pointer",
		":hover": {
			backgroundColor: resendSuccess ? "#4CAF50" : "#81C784",
		},
	})
	const paddingLeft: string = css({
		paddingLeft: "10px",
		paddingTop: "5px",
	})

	return (
		<div className={containerStyle} onClick={resendSuccess ? undefined : () => resendVerify()}>
			{resendSuccess ? (
				<>Verification Email Sent</>
			) : (
				<>
					<div>Please click here to verify your email</div>
					{mutResendVerify.loading && (
						<div className={paddingLeft}>
							<Spinner size="20px" color="white" />
						</div>
					)}
				</>
			)}
		</div>
	)
}
