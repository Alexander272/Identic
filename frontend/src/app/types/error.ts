export interface IFetchError {
	data: {
		message: string
		code: string
	}
}

export interface IBaseFetchError {
	error: {
		data: {
			message: string
			code: string
		}
		status: number
	}
	status: number
}
