import { extractQuantity } from '../extract'
import type { ParsedItem } from './types'

// ================= Утилиты =================
export function normalizeText(str: string | null | undefined): string {
	return (str || '')
		.replace(/\u00A0/g, ' ') // &nbsp; → пробел
		.replace(/[\u200B-\u200D\uFEFF]/g, '') // невидимые символы
		.replace(/\s+/g, ' ')
		.trim()
}

export function parseQuantity(text: string | null | undefined): number | null {
	if (!text) return null
	// Ищем первое число, игнорируя разделители и пробелы
	const match = text.match(/(\d[\d\s.,]*)/)
	if (!match) return null

	const cleaned = match[1].replace(/[.,\s]/g, '.')
	const value = parseFloat(cleaned)
	return isNaN(value) || value <= 0 ? null : value
}

// ================= Таблицы / Grid =================
export function getRowCells(row: HTMLElement): string[] {
	return Array.from(row.children).map(el => normalizeText(el.textContent))
}

export function findHeaderRow<T extends HTMLElement>(rows: T[], namePattern: RegExp, qtyPattern: RegExp): T | null {
	for (const row of rows) {
		const cells = getRowCells(row)
		if (cells.some(c => namePattern.test(c)) && cells.some(c => qtyPattern.test(c))) {
			return row
		}
	}
	return null
}

export function findColumnIndexes(
	cells: string[],
	namePattern: RegExp,
	qtyPattern: RegExp,
): { nameIndex: number; qtyIndex: number } | null {
	const nameIndex = cells.findIndex(c => namePattern.test(c))
	const qtyIndex = cells.findIndex(c => qtyPattern.test(c))
	return nameIndex === -1 || qtyIndex === -1 ? null : { nameIndex, qtyIndex }
}

export function extractFromGrid<T extends HTMLElement>(
	rows: T[],
	startIndex: number,
	nameIndex: number,
	qtyIndex: number,
	minValid: number,
): ParsedItem[] {
	const result: ParsedItem[] = []

	for (let i = startIndex + 1; i < rows.length; i++) {
		const cells = Array.from(rows[i].children)
		const getCellText = (idx: number) => cells[idx]?.textContent?.replace(/\u00A0/g, ' ').trim() || null

		const nameRaw = getCellText(nameIndex)
		const qtyRaw = getCellText(qtyIndex)

		const name = nameRaw?.replace(/\s+/g, ' ')
		const quantity = parseQuantity(qtyRaw)

		if (name && name.length >= 2 && quantity) {
			result.push({ name, quantity })
		}
	}

	return result.length >= minValid ? result : []
}

// ================= Плоские текстовые списки =================
const INLINE_REGEX = /^(.+?)\s*[—\-–:]\s*(\d[\d\s.,]*)/i

export function extractFromInlineText(elements: Element[], minValid: number): ParsedItem[] {
	const result: ParsedItem[] = []

	for (const el of elements) {
		const text = el.textContent?.replace(/\u00A0/g, ' ').trim()
		if (!text) continue

		const match = text.match(INLINE_REGEX)
		if (!match) continue

		const name = match[1].replace(/\s+/g, ' ').trim()
		const quantity = parseQuantity(match[2])

		if (name.length >= 3 && quantity) {
			result.push({ name, quantity })
		}
	}

	return result.length >= minValid ? result : []
}

// ================= Паттерны =================
export const DEFAULT_NAME_PATTERN = /наим|товар|назван|номен|item|product|description|позиция/i
export const DEFAULT_QTY_PATTERN = /кол|qty|кол-во|quantity|шт|pcs|amount|count/i

/** Собирает текстовые блоки из DOM */
export function extractRawLines(doc: Document): string[] {
	const elements = doc.querySelectorAll('div, tr, li, p')
	const lines: string[] = []

	for (const el of elements) {
		const hasNested = el.querySelector('div, tr, li, p')
		if (hasNested) continue // Пропускаем элементы с вложенными

		const text = normalizeText(getSmartText(el))
		console.log('normalizeText', text)

		if (text.length > 5) {
			lines.push(text)
		}
	}
	return lines
}

function getSmartText(el: Element): string {
	if (el.tagName === 'TR') {
		// Склеиваем ячейки таблицы через пробел, чтобы <td> не слипались
		return Array.from(el.children)
			.map(child => child.textContent?.trim() || '')
			.join(' ')
	}
	// Для остальных (p, li, div) берем текст как есть
	return el.textContent || ''
}
/** Парсит одну строку в { name, quantity } */
export function parseLine(line: string, minNameLen: number): ParsedItem | null {
	// 1. Убираем ведущую нумерацию: "1 ", "№2 ", "3. "
	const text = line.replace(/^\s*(?:№|п\/п|п\.|n)?\s*\d+[.)\-\s]*/i, '').trim()
	if (text.length < minNameLen) return null

	// let name: string | null = null
	// let qtyStr: string | null = null

	// 2. Формат с разделителем: "Название — 3 шт"
	// const sepMatch = text.match(/^(.+?)\s+[—\-–:]\s+(\d[\d\s.,]*)/)
	// if (sepMatch) {
	// 	name = sepMatch[1].trim()
	// 	qtyStr = sepMatch[2]
	// } else {
	// 	// 3. Формат с числом в конце: "Название 16 ШТ"
	// 	const endMatch = text.match(/^(.+)\s+(\d[\d\s.,]*)\s*(?:шт|pcs|ед\.?|unit|ШТ)?\s*$/i)
	// 	if (endMatch) {
	// 		name = endMatch[1].trim()
	// 		qtyStr = endMatch[2]
	// 	}
	// }

	const { name, quantity } = extractQuantity(text)

	// if (!name || !qtyStr) return null

	// const quantity = parseQuantity(qtyStr)
	// if (!quantity || name.length < minNameLen) return null

	return { name, quantity }
}
