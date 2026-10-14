import { createTextColumn } from 'react-datasheet-grid'

import { parseDecimal } from '@/utils/number'

/**
 * Колонка грида для дробных чисел, принимающая запятую в качестве разделителя.
 * Замена floatColumn из react-datasheet-grid: тот использует parseFloat и молча обрезает "1,5" до 1.
 */
export const decimalColumn = createTextColumn<number | null>({
	alignRight: true,
	parseUserInput: parseDecimal,
	parsePastedValue: parseDecimal,
	formatBlurredInput: value =>
		typeof value === 'number' ? new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 20 }).format(value) : '',
	formatInputOnFocus: value => (typeof value === 'number' ? String(value) : ''),
})
