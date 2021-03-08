import { createContainer } from "unstated-next"
import * as React from "react"
import { UserError } from "../types/types"
import { ApolloError } from "apollo-client"
import { resetClient } from "../app"

interface IAuthError {
	error: string
	message: string
}

const getUserError = async (response: Response): Promise<string> => {
	const defaultMessage: string = "Something went wrong. Please try again."
	try {
		const jsonResp: IAuthError = await response.json()
		if (jsonResp.message) {
			return jsonResp.message
		}
		return defaultMessage
	} catch {
		return defaultMessage
	}
}

const useAuth = () => {
	const [userErrors, setUserErrors] = React.useState<UserError[]>([])
	const [apolloError, setApolloError] = React.useState<ApolloError | undefined>(undefined)
	const [loggedIn, setLoggedIn] = React.useState<boolean>(false)
	const [logoutRedirect, setLogoutRedirect] = React.useState<boolean>(false)
	const [loading, setLoading] = React.useState<boolean>(false)
	const [showVerifyComplete, setShowVerifyComplete] = React.useState<boolean>(false)

	const login = async (email: string, password: string) => {
		await setLoading(true)
		const response = await fetch("/api/auth/login", {
			method: "POST",
			body: JSON.stringify({ email, password }),
		})
		if (response.status === 200) {
			setLoading(false)
			resetClient()
			clearAuthErrors()
			setLoggedIn(true)
			return
		}
		const userError: string = await getUserError(response)
		setUserErrors([{ message: userError, field: [] }])
		setLoading(false)
		return
	}

	const logout = async () => {
		await setLoading(true)
		const response = await fetch("/api/auth/logout", {
			method: "POST",
		})
		if (response.status === 200) {
			setLoading(false)
			resetClient()
			setLoggedIn(false)
			clearAuthErrors()
			setLogoutRedirect(true)
			return
		}
		setUserErrors([{ message: "There was a problem logging you out", field: [] }])
		setLoading(false)
		return
	}

	const verify = async (token: string) => {
		await setLoading(true)
		const response = await fetch("/api/auth/verify_account", {
			method: "POST",
			body: JSON.stringify({ token }),
		})
		if (response.status === 200) {
			setLoading(false)
			resetClient()
			clearAuthErrors()
			setLoggedIn(true)
			setShowVerifyComplete(true)
			return
		}
		const userError: string = await getUserError(response)
		setUserErrors([{ message: userError, field: [] }])
		setLoading(false)
		return
	}

	const clearAuthErrors = () => {
		setUserErrors([])
		setApolloError(undefined)
	}

	return {
		login,
		logout,
		verify,
		loading,
		userErrors,
		apolloError,
		clearAuthErrors,
		loggedIn,
		logoutRedirect,
		setLogoutRedirect,
		showVerifyComplete,
		setShowVerifyComplete,
	}
}

export const AuthContainer = createContainer(useAuth)
