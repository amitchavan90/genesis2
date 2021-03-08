import * as React from "react"
import { ErrorMap } from "../types/types"
import { UserContainer } from "../controllers/user"
import { ErrorNotification } from "./errorBox"
import { KIND, Notification } from "baseui/notification"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Spaced } from "./spaced"
import { Spinner } from "baseui/spinner"
import { ModalButton } from "baseui/modal"

export const ChangePassword = () => {
	const [inputError, setInputError] = React.useState<ErrorMap>({})
	const [oldPassword, setOldPassword] = React.useState<string>("")
	const [password, setPassword] = React.useState<string>("")

	const { apolloError, loading, clearErrors, useChangePassword } = UserContainer.useContainer()
	const { changePassword, changeSuccess, setChangeSuccess } = useChangePassword(oldPassword, password)

	const handleSubmit = () => {
		setChangeSuccess(false)
		clearErrors()

		const errors: ErrorMap = {}
		let foundError = false

		if (oldPassword == "") {
			errors["oldPassword"] = "Please enter your current password"
			foundError = true
		}
		if (password == "") {
			errors["password"] = "Please enter a password"
			foundError = true
		}

		if (foundError) {
			setInputError(errors)
			return
		}

		setInputError({})
		changePassword()
	}

	return (
		<div>
			<div>
				{changeSuccess && <Notification kind={KIND.positive}>Your password has been updated.</Notification>}

				{apolloError && <ErrorNotification message={apolloError.message} />}

				<FormControl label="Old Password" error={inputError["oldPassword"]}>
					<Input
						key={"oldPassword"}
						error={!!inputError["oldPassword"]}
						positive={false}
						value={oldPassword}
						type={"password"}
						onChange={e => setOldPassword(e.currentTarget.value)}
						placeholder={"Enter your current password"}
					/>
				</FormControl>
				<FormControl label="New Password" error={inputError["password"]}>
					<Input
						key={"password"}
						error={!!inputError["password"]}
						positive={false}
						value={password}
						type={"password"}
						onChange={e => setPassword(e.currentTarget.value)}
						placeholder={"Enter a new password"}
					/>
				</FormControl>
			</div>
			<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
				{loading && <Spinner />}
				<ModalButton onClick={handleSubmit}>Save</ModalButton>
			</Spaced>
		</div>
	)
}
