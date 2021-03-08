/** Returns a promise race between the provided promise and a timeout (30 seconds by default) */
export const promiseTimeout = (promise: Promise<any>, seconds?: number) => {
	const timeout = new Promise<void>((_, reject) => {
		const id = setTimeout(() => {
			clearTimeout(id)
			reject(TIMED_OUT)
		}, (seconds || 30) * 1000)
	})
	return Promise.race([promise, timeout])
}

export const TIMED_OUT = "Timed out"
