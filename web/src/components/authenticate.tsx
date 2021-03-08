import * as React from "react"
import { Button } from "baseui/button"
import { useForm } from "react-hook-form"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Notification, KIND } from "baseui/notification"
import { useStyletron } from "baseui"
import { Block } from "baseui/block"
import { useHistory } from "react-router"

interface Props {
	setShowLogin?: (showLogin: boolean) => void
}

export const Authenticate = (props: Props) => {
	const [css, theme] = useStyletron()
	const history = useHistory()

	const { register, handleSubmit, watch, errors } = useForm<{ email: string; password: string }>()
	const onSubmit = (data: Record<string, any>) => {
		console.log(data)
	}
	return (
		<Block maxWidth={["420px"]} paddingLeft={["scale800", "scale1200"]} paddingRight={["scale800", "scale1200"]}>
			<form onSubmit={handleSubmit(onSubmit)}>
				<FormControl label="Email">
					<Input name="email" inputRef={register({ required: { value: true, message: "Email is required" } })} />
				</FormControl>
				{errors.email && (
					<Notification
						overrides={{
							Body: {
								style: ({ $theme }) => {
									return {
										width: "auto",
									}
								},
							},
						}}
						kind={KIND.warning}
					>
						{errors.email.message}
					</Notification>
				)}
				<FormControl label="Password">
					<Input name="password" type="password" inputRef={register({ required: { value: true, message: "Password is required" } })} />
				</FormControl>
				{errors.password && (
					<Notification
						kind={KIND.warning}
						overrides={{
							Body: {
								style: ({ $theme }) => {
									return {
										width: "auto",
									}
								},
							},
						}}
					>
						{errors.password.message}
					</Notification>
				)}
				<Button
					overrides={{
						BaseButton: {
							style: ({ $theme }) => {
								return {
									width: "100%",
								}
							},
						},
					}}
					type={"submit"}
				>
					Sign in
				</Button>
			</form>

			<hr />
			<Button
				onClick={() => {
					history.push("/create_account")
				}}
				overrides={{
					BaseButton: {
						style: ({ $theme }) => {
							return {
								width: "100%",
							}
						},
					},
				}}
			>
				Create account
			</Button>
			<p
				onClick={() => {
					history.push("/forgot_password")
				}}
			>
				Forgot your password?
			</p>
		</Block>
	)
}
