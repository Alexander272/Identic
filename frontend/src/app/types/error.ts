export interface IFetchError {
	data: {
		message: string
		code: string
	}
	status: number
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
