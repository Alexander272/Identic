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

	let structureFound = false
	// ===== 1. Классические <table> =====
	const tables = Array.from(doc.querySelectorAll('table'))
	if (tables.length > 0) structureFound = true
	for (const table of tables) {
		console.log('parse Table')

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

	// ===== 2. CSS-Grid / Div-таблицы =====
	const gridCandidates = Array.from(
		doc.querySelectorAll<HTMLElement>(
			'[class*="row"], [class*="item"], [class*="list"], [class*="grid-cell"], [class*="k-grid"], [data-qa="cell"]',
		),
	)
		.map(el => {
			if (el.getAttribute('data-qa') === 'cell' || el.classList.contains('controls-Grid__header-cell')) {
				return el.parentElement
			}
			return el
		})
		.filter(
			(el, index, self) => el && el.children.length >= 2 && self.indexOf(el) === index, // Убираем дубликаты родителей
		) as HTMLElement[]

	if (gridCandidates.length >= 2) {
		console.log('parser grid')

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

	// Выходим если нашли таблицу, но её структура невалидна
	if (structureFound) {
		console.log('Structure detected but headers were invalid. Aborting before text parse.')
		return []
	}

	// ===== 3. Плоские текстовые строки =====
	// Исправлено: ищем элементы с прямым текстовым содержимым, а не только с 1-2 детьми
	const lines = extractRawLines(doc)
	const parsed: ParsedItem[] = []
	console.log('parse text')

	// 2. Парсим каждую строку
	for (const line of lines) {
		const item = parseLine(line, minNameLen)
		if (item) parsed.push(item)
	}

	return parsed.length >= minValid ? parsed : []
}
