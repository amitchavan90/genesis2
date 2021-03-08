import * as React from "react"
import { Accordion, Panel } from "baseui/accordion"
import { H2 } from "baseui/typography"
import { UserContainer } from "../../controllers/user"
import { ChangePassword } from "../../components/changePassword"
import { ChangeSettings } from "../../components/changeDetails"
import { CenteredPage } from "../../components/common"

const Settings = () => {
	const { clearErrors } = UserContainer.useContainer()

	return (
		<CenteredPage>
			<H2>User Settings</H2>
			<Accordion
				onChange={({ expanded }) => {
					clearErrors()
				}}
			>
				<Panel title="Your Details">
					<ChangeSettings />
				</Panel>
				<Panel title="Change Password">
					<ChangePassword />
				</Panel>
			</Accordion>
		</CenteredPage>
	)
}

export default Settings
