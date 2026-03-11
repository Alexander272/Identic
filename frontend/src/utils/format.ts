export const numberFormat = new Intl.NumberFormat('ru-Ru').format

export const removeSpace = (values: string[]) => {
	return values.map(v => v.replace(/\s+/g, ''))
}
