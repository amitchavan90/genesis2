import { ErrorNotification } from "./errorBox"
import { Notification } from "baseui/notification"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import * as React from "react"
import { UserContainer } from "../controllers/user"
import { Spinner } from "baseui/spinner"
import { Spaced } from "./spaced"
import { useForm } from "react-hook-form"
import { Select, Value } from "baseui/select"
import { Button } from "baseui/button"
import { AffiliateOrgs } from "../types/enums"

type FormData = {
	email: string
	firstName: string
	lastName: string
	mobilePhone: string
}

export const ChangeSettings = () => {
	const { apolloError, loading, useChangeDetails, user } = UserContainer.useContainer()
	const { changeDetails, changeSuccess, setChangeSuccess } = useChangeDetails()

	// Form Subbmision
	const [affiliateOrg, setAffiliateOrg] = React.useState<Value>()

	const { register, setValue, handleSubmit, errors } = useForm<FormData>()
	const onSubmit = handleSubmit(({ email, firstName, lastName, mobilePhone }) => {
		if (!user) return

		const input = {
			email,
			firstName,
			lastName,
			affiliateOrg: affiliateOrg && affiliateOrg.length > 0 ? affiliateOrg[0].id : undefined,
			mobilePhone,
		}

		changeDetails({ variables: { input } })
	})

	// On load user
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
	}, [user])

	return (
		<form onSubmit={onSubmit}>
			{apolloError && <ErrorNotification message={apolloError.message} />}

			{changeSuccess && (
				<Notification closeable kind="positive" onClose={() => setChangeSuccess(false)}>
					Your details have been updated.
				</Notification>
			)}

			<FormControl label="Email" error={errors.email?.message}>
				<Input name="email" type="email" inputRef={register({ required: "Required" })} />
			</FormControl>
			<FormControl label="First Name" error={errors.firstName?.message}>
				<Input name="firstName" inputRef={register({ required: "Required" })} />
			</FormControl>
			<FormControl label="Last Name" error={errors.lastName?.message}>
				<Input name="lastName" inputRef={register({ required: "Required" })} />
			</FormControl>
			<FormControl label="Mobile" error={errors.mobilePhone?.message}>
				<Input name="mobilePhone" inputRef={register} />
			</FormControl>

			<FormControl label="Affiliate Organization">
				<Select creatable options={AffiliateOrgs} labelKey="id" valueKey="id" value={affiliateOrg} onChange={({ value }) => setAffiliateOrg(value)} />
			</FormControl>

			<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
				{loading && <Spinner />}
				<Button>Save</Button>
			</Spaced>
		</form>
	)
}
