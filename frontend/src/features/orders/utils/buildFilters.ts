import type { IFilter } from '../types/params'

export const buildFilterUrlParams = (req: IFilter[]): URLSearchParams => {
	const params: string[][] = []

	for (let i = 0; i < req.length; i++) {
		const f = req[i]

		if (i == 0 || f.field != req[i - 1].field) {
			params.push([`filters[${f.field}]`, ''])
		}
		params.push([`${f.field}[${f.compareType}]`, f.value])
	}

	return new URLSearchParams(params)
}
