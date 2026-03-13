// const reWithUnit = /(\d+)\s*(?:шт|кг)|(?:шт|кг)\.?\s*(\d+)/i
const reWithUnit = /([\d,.]+)\s*(?:шт|кг)$|(?:шт|кг)\.?\s*([\d,.]+)$/i
const reEndDigits = /\s+(\d+)$/

/**
 * Функция парсит строку, извлекает количество и возвращает очищенное наименование.
 * @param {string} input - Строка из Excel
 * @returns {{ name: string, quantity: number }}
 */
export const extractQuantity = (input: string) => {
	const s = input.trim()

	// 1. Ищем паттерны с "шт" (регистронезависимо)
	// Ищем число до или после "шт", "шт.", "шт "
	const unitMatch = s.match(reWithUnit)

	if (unitMatch) {
		const qtyStr = unitMatch[1] || unitMatch[2]
		const quantity = parseFloat(qtyStr.replace(',', '.'))

		// Удаляем найденное количество и "шт" из названия
		const name = s.replace(reWithUnit, '').replace(/\s+/g, ' ').trim()
		return { name, quantity }
	}

	// 2. Если "шт" нет, проверяем, не заканчивается ли строка на стандарт (ГОСТ, ОСТ, ТУ и т.д.)
	// Если в конце стандарт, то числа там — это часть стандарта, их трогать нельзя.
	const isStandard = /(ОСТ|ГОСТ|ТУ|ASME|B|Series)\s*[\d.\-\s]+$/i.test(s)

	if (!isStandard) {
		// Ищем число в самом конце строки, которому предшествует пробел
		const endMatch = s.match(reEndDigits)

		if (endMatch) {
			const quantity = parseFloat(endMatch[1].replace(',', '.'))
			const name = s.replace(reEndDigits, '').trim()
			return { name, quantity }
		}
	}

	// 3. Если ничего не нашли, возвращаем как есть
	return { name: s, quantity: 0 }
}
