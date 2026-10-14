import { extractQuantity } from '../extract'
import { parseDecimal } from '@/utils/number'
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

	const value = parseDecimal(match[1])
	return value === null || value <= 0 ? null : value
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
	expectedColumns = 0,
): ParsedItem[] {
	const result: ParsedItem[] = []

	for (let i = startIndex + 1; i < rows.length; i++) {
		const cells = Array.from(rows[i].children)
		// защита от съезда колонок: строка с другим числом ячеек читается по позиции некорректно
		if (expectedColumns && cells.length !== expectedColumns) continue

		const getCellText = (idx: number) => normalizeText(cells[idx]?.textContent) || null

		const nameRaw = getCellText(nameIndex)
		const qtyRaw = getCellText(qtyIndex)

		const name = nameRaw?.replace(/\s+/g, ' ')
		const quantity = parseQuantity(qtyRaw)
		console.log(nameRaw, name, '|', qtyRaw, quantity)

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
// ================= Нумерация списка =================
// Явная нумерация снимается всегда: "1. ", "1) ", "№2 ", "п/п 3 ", "п. 3 "
const DEFINITE_PREFIX = /^\s*(?:(?:№\s*\d+|п\/п\s*\d+|п\.\s*\d+)\s*|\d+\.\s+|\d+\)\s*)/i
// Голое "N " снимается только внутри последовательности номеров
const BARE_PREFIX = /^\s*(\d+)\s+/

/**
 * Снимает нумерацию списка, не задевая названия, начинающиеся с числа ("3.5 мм", "40x40").
 * Явные маркеры снимаются сразу, голое "N " — только если строки образуют последовательность 1, 2, 3…
 */
export function stripListNumbering(lines: string[]): string[] {
	const definite = lines.map(line => {
		const match = line.match(DEFINITE_PREFIX)
		return match ? line.slice(match[0].length) : null
	})

	const bare = lines.map(line => {
		const match = line.match(BARE_PREFIX)
		return match ? { number: Number(match[1]), rest: line.slice(match[0].length) } : null
	})

	const stripBare = new Array<boolean>(lines.length).fill(false)

	let i = 0
	while (i < lines.length) {
		const start = bare[i]
		if (!start) {
			i++
			continue
		}

		let j = i
		while (j + 1 < lines.length && bare[j + 1] && start.number + (j + 1 - i) === bare[j + 1]!.number) {
			j++
		}

		// снимаем, только если нумерация образует последовательность минимум из двух строк
		if (j > i) {
			for (let k = i; k <= j; k++) stripBare[k] = true
		}
		i = j + 1
	}

	return lines.map((line, idx) => {
		if (definite[idx] !== null) return definite[idx]!
		if (stripBare[idx]) return bare[idx]!.rest
		return line
	})
}

/** Парсит одну строку в { name, quantity } */
export function parseLine(line: string, minNameLen: number): ParsedItem | null {
	const text = line.trim()
	if (text.length < minNameLen) return null

	const { name, quantity } = extractQuantity(text)
	return { name, quantity }
}
