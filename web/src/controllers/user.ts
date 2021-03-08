import { createContainer } from "unstated-next"
import * as React from "react"
import { UserError, User } from "../types/types"
import { useMutation, useQuery } from "@apollo/react-hooks"
import { ApolloError } from "apollo-client"
import { graphql } from "../graphql"
import { Perm } from "../types/enums"

const useUser = () => {
	const [userErrors, setUserErrors] = React.useState<UserError[]>([])
	const [apolloError, setApolloError] = React.useState<ApolloError | undefined>(undefined)
	const [loading, setLoading] = React.useState<boolean>(false)
	const [user, setUser] = React.useState<User | undefined>(undefined)

	const clearErrors = () => {
		setUserErrors([])
		setApolloError(undefined)
	}

	const useChangePassword = (oldPassword: string, password: string) => {
		const [changePassword, { data, loading, error }] = useMutation(graphql.mutation.CHANGE_PASSWORD, {
			variables: { oldPassword, password },
		})
		const [changeSuccess, setChangeSuccess] = React.useState<boolean>(false)

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)

			if (!data) return
			setChangeSuccess(data.changePassword)
		}, [data, loading, error])

		return { changePassword, changeSuccess, setChangeSuccess }
	}

	const useFetchUser = () => {
		const { data, loading, error } = useQuery(graphql.query.ME)
		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data && data.me) {
				setUser(data.me)
			}
			return
		}, [data, loading, error])
	}

	const useChangeDetails = () => {
		const [changeDetails, { data, loading, error }] = useMutation<{ changeDetails: User }>(graphql.mutation.CHANGE_DETAILS)
		const [changeSuccess, setChangeSuccess] = React.useState<boolean>(false)

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)

			if (!data) return
			setChangeSuccess(true)
			setUser(data.changeDetails)
		}, [data, loading, error])

		return { changeDetails, changeSuccess, setChangeSuccess }
	}

	const hasPermission = (perm: Perm) => {
		if (!user || !user.role || !user.role.permissions) return false
		return user.role.permissions.includes(perm)
	}

	return {
		loading,
		userErrors,
		apolloError,
		useChangePassword,
		user,
		useFetchUser,
		clearErrors,
		useChangeDetails,
		hasPermission,
	}
}

export const UserContainer = createContainer(useUser)
