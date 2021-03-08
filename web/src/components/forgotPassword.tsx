import * as React from "react"
import { Button } from "baseui/button"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Modal, ModalHeader, ModalBody } from "baseui/modal"
import { useMutation } from "@apollo/react-hooks"
import { graphql } from "../graphql"
import { useForm } from "react-hook-form"
import { ErrorNotification } from "./errorBox"
import { Notification } from "baseui/notification"
import { Spaced } from "./spaced"
import { Spinner } from "baseui/spinner"

export const ForgotPasswordModal = (props: { isOpen: boolean; onClose: () => void }) => {
	const [forgotPass, mutForgotPass] = useMutation<{ forgotPassword: boolean }>(graphql.mutation.FORGOT_PASSWORD)
	const [changeSuccess, setChangeSuccess] = React.useState<boolean>()
	const { register, handleSubmit, errors } = useForm<{ email: string }>()

	const onSubmit = handleSubmit(({ email }) => {
		setChangeSuccess(false)

		forgotPass({ variables: { email } })
	})

	React.useEffect(() => {
		if (!mutForgotPass.data) return
		setChangeSuccess(mutForgotPass.data.forgotPassword)
	}, [mutForgotPass])

	return (
		<Modal onClose={() => props.onClose()} isOpen={props.isOpen}>
			<ModalHeader>Forgot Password</ModalHeader>
			<ModalBody>
				{changeSuccess === true && <div>Email Sent!</div>}

				{mutForgotPass.error && <ErrorNotification message={mutForgotPass.error.message} />}

				{!changeSuccess && (
					<form onSubmit={onSubmit}>
						<FormControl error={errors.email?.message}>
							<Input
								type="email"
								name="email"
								inputRef={register({
									required: "Required",
									pattern: {
										value: /.+\@.+\..+/,
										message: "Invalid email address",
									},
								})}
								disabled={mutForgotPass.loading}
								error={errors.email !== undefined}
							/>
						</FormControl>

						<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
							{mutForgotPass.loading && <Spinner />}
							<Button>Send Password Reset Email</Button>
						</Spaced>
					</form>
				)}
			</ModalBody>
		</Modal>
	)
}
