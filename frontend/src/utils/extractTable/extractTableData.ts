import type { ExtractedData, ExtractConfig, ParsedItem } from './types'
import {
	findHeaderRow,
	findColumnIndexes,
	getRowCells,
	extractFromGrid,
	DEFAULT_NAME_PATTERN,
	DEFAULT_QTY_PATTERN,
	extractRawLines,
	parseLine,
} from './helpers'

export function extractTableData(html: string, config: ExtractConfig = {}): ExtractedData {
	if (!html || typeof html !== 'string') return []

	const minValid = config.minValidRows ?? 1
	const minNameLen = config.minNameLength ?? 3

	let doc: Document
	try {
		const parser = new DOMParser()
		doc = parser.parseFromString(html, 'text/html')
		if (doc.querySelector('parsererror')) return []
	} catch {
		return []
	}

	// ===== 1. Классические <table> =====
	const tables = Array.from(doc.querySelectorAll('table'))
	for (const table of tables) {
		const rows = Array.from(table.querySelectorAll('tr'))
		if (rows.length < 2) continue

		const header = findHeaderRow(rows, DEFAULT_NAME_PATTERN, DEFAULT_QTY_PATTERN)
		if (!header) continue

		const cells = getRowCells(header)
		const indexes = findColumnIndexes(cells, DEFAULT_NAME_PATTERN, DEFAULT_QTY_PATTERN)
		if (!indexes) continue

		const result = extractFromGrid(rows, rows.indexOf(header), indexes.nameIndex, indexes.qtyIndex, minValid)
		if (result.length) return result
	}

	// ===== 2. CSS-Grid / Div-таблицы (ваш пример 2) =====
	const gridCandidates = Array.from(
		doc.querySelectorAll<HTMLElement>('[class*="row"], [class*="item"], [class*="list"]'),
	).filter(el => el.children.length >= 2)

	if (gridCandidates.length >= 2) {
		const header = findHeaderRow(gridCandidates, DEFAULT_NAME_PATTERN, DEFAULT_QTY_PATTERN)
		if (header) {
			const cells = getRowCells(header)
			const indexes = findColumnIndexes(cells, DEFAULT_NAME_PATTERN, DEFAULT_QTY_PATTERN)
			if (indexes) {
				const result = extractFromGrid(
					gridCandidates,
					gridCandidates.indexOf(header),
					indexes.nameIndex,
					indexes.qtyIndex,
					minValid,
				)
				if (result.length) return result
			}
		}
	}

	// ===== 3. Плоские текстовые строки (ваш пример 1) =====
	// Исправлено: ищем элементы с прямым текстовым содержимым, а не только с 1-2 детьми
	const lines = extractRawLines(doc)
	const parsed: ParsedItem[] = []

	// 2. Парсим каждую строку
	for (const line of lines) {
		const item = parseLine(line, minNameLen)
		if (item) parsed.push(item)
	}

	// 3. Убираем возможные дубли (из-за вложенных div/tr)
	const unique = new Map<string, ParsedItem>()
	for (const item of parsed) {
		const key = `${item.name}|${item.quantity}`
		if (!unique.has(key)) unique.set(key, item)
	}

	const result = Array.from(unique.values())
	return result.length >= minValid ? result : []
}
