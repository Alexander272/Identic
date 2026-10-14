/**
 * Приводит пользовательский ввод к числу.
 * Поддерживает запятую и точку как десятичный разделитель,
 * а также пробелы и неразрывные пробелы как разделители разрядов: "1 000,5".
 * Возвращает null, если строку не удалось распознать как число.
 */
export const parseDecimal = (raw: string): number | null => {
	const value = raw.replace(/\u00A0/g, ' ').trim()
	if (!value) return null

	const lastComma = value.lastIndexOf(',')
	const lastDot = value.lastIndexOf('.')

	let normalized = value
	if (lastComma !== -1 && lastDot !== -1) {
		// присутствуют оба разделителя — правый считаем десятичным, левый разрядами
		const decimalSeparator = lastComma > lastDot ? ',' : '.'
		const groupSeparator = decimalSeparator === ',' ? '.' : ','
		normalized = value.split(groupSeparator).join('').replace(decimalSeparator, '.')
	} else if (lastComma !== -1) {
		normalized = value.replace(/,/g, '.')
	}

	normalized = normalized.replace(/\s+/g, '')

	const number = Number(normalized)
	return Number.isFinite(number) ? number : null
}
