import gql from "graphql-tag"
import { createContainer } from "unstated-next"
import * as React from "react"
import { UserError } from "../types/types"
import { ServerPlan, ServerInstance } from "../types/server"
import { useMutation } from "@apollo/react-hooks"
import { ApolloError } from "apollo-client"

const QUERY_VPS_PLAN_LIST = gql`
	mutation onboardStart($email: String!) {
		onboardStart(email: $email) {
			userErrors {
				message
				field
			}
			onboard {
				prospect {
					id
					email
					firstName
					lastName
				}
			}
		}
	}
`

const QUERY_VPS_LIST = gql`
	{
		vultrServerList {
			InstanceID
			servers {
				AllowedBandwidth
				AppID
			}
		}
	}
`

const MUTATION_VPS_CREATE = gql`
	mutation onboardUpdate($id: String!, $input: OnboardingInput!) {
		onboardUpdate(id: $id, input: $input) {
			userErrors {
				message
				field
			}
			onboard {
				prospect {
					id
					email
					firstName
					lastName
				}
			}
		}
	}
`
const MUTATION_VPS_DELETE = gql``
const MUTATION_VPS_START = gql``
const MUTATION_VPS_STOP = gql``
const MUTATION_VPS_RESTART = gql``

const useOnboarding = () => {
	// const [submitError, setSubmitError] = React.useState<UserError[]>([])
	// const [apolloError, setApolloError] = React.useState<ApolloError | undefined>(undefined)
	// const [prospect, setProspect] = React.useState<Prospect>({
	// 	id: "",
	// 	email: "",
	// 	firstName: "",
	// 	lastName: "",
	// 	onboardingComplete: false,
	// })
	// const [loading, setLoading] = React.useState<boolean>(false)
	// const [current, setCurrent] = React.useState<number>(0)
	// const useStartOnboarding = () => {
	// 	const [startOnboarding, { data, loading, error }] = useMutation<{ onboardStart: OnboardingOutput }, { email: string }>(MUTATION_ONBOARD_START, {
	// 		variables: { email: prospect.email },
	// 	})
	// 	React.useEffect(() => {
	// 		setLoading(loading)
	// 		setApolloError(error)
	// 		if (data && data.onboardStart.userErrors.length > 0) {
	// 			setSubmitError(data.onboardStart.userErrors)
	// 			return
	// 		}
	// 		if (data) {
	// 			setProspect({ ...prospect, email: data.onboardStart.onboard.prospect.email, id: data.onboardStart.onboard.prospect.id })
	// 			stepForward()
	// 		}
	// 		return
	// 	}, [data, loading, error])
	// 	return { startOnboarding }
	// }
	// const useUpdateOnboarding = () => {
	// 	const input: OnboardingInput = {
	// 		email: prospect.email,
	// 		lastName: prospect.lastName,
	// 		firstName: prospect.firstName,
	// 	}
	// 	const [updateOnboarding, { data, loading, error }] = useMutation<{ onboardUpdate: OnboardingOutput }, { id: string; input: OnboardingInput }>(
	// 		MUTATION_ONBOARD_UPDATE,
	// 		{
	// 			variables: { id: prospect.id, input: input },
	// 		},
	// 	)
	// 	React.useEffect(() => {
	// 		setLoading(loading)
	// 		setApolloError(error)
	// 		if (data && data.onboardUpdate.userErrors.length > 0) {
	// 			setSubmitError(data.onboardUpdate.userErrors)
	// 			return
	// 		}
	// 		if (data) {
	// 			setProspect(data.onboardUpdate.onboard.prospect)
	// 			stepForward()
	// 		}
	// 	}, [data, loading, error])
	// 	return { updateOnboarding }
	// }
	// const useFinishOnboarding = () => {
	// 	const [finishOnboarding, { data, loading, error }] = useMutation<{ onboardFinish: OnboardingOutput }, { id: string }>(MUTATION_ONBOARD_FINISH, {
	// 		variables: { id: prospect.id },
	// 	})
	// 	React.useEffect(() => {
	// 		setLoading(loading)
	// 		setApolloError(error)
	// 		if (data && data.onboardFinish.userErrors.length > 0) {
	// 			setSubmitError(data.onboardFinish.userErrors)
	// 			return
	// 		}
	// 		if (data) {
	// 			stepForward()
	// 		}
	// 	}, [data, loading, error])
	// 	return { finishOnboarding }
	// }
	// const stepForward = () => {
	// 	setCurrent(current + 1)
	// }
	// const stepBack = () => {
	// 	setCurrent(current - 1)
	// }
	// return {
	// 	useStartOnboarding,
	// 	useUpdateOnboarding,
	// 	useFinishOnboarding,
	// 	stepBack,
	// 	current,
	// 	loading,
	// 	submitError,
	// 	apolloError,
	// 	prospect,
	// 	setProspect,
	// }
}

export const Onboarding = createContainer(useOnboarding)
