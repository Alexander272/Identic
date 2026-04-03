export interface ExtractConfig {
	/** Минимальная длина имени для отсечения мусора */
	minNameLength?: number
	/** Минимальное количество совпадений для успешного возврата */
	minValidRows?: number
}

/** Формат выходных данных */
export type OutputFormat = 'object' | 'string'

/** Элемент результата в объектном формате */
export interface ParsedItem {
	name: string
	quantity: number
}

export type ExtractedData = ParsedItem[]

/** Настройки парсинга */
export interface ExtractTableOptions {
	/** Формат возврата: 'object' | 'string' */
	outputFormat?: OutputFormat
	/** Единица измерения для строкового формата */
	defaultUnit?: string
	/** Минимальная длина названия товара */
	minNameLength?: number
	/** Разрешить дробные количества */
	allowDecimals?: boolean
	/** Дополнительные паттерны для поиска колонки названия */
	namePatterns?: RegExp
	/** Дополнительные паттерны для поиска колонки количества */
	quantityPatterns?: RegExp
}

/** Внутренние индексы колонок */
export interface ColumnIndexes {
	nameIndex: number
	qtyIndex: number
}

/** Контекст парсера — передаётся между функциями */
export interface ParserContext {
	namePatterns: RegExp
	qtyPatterns: RegExp
	options: Required<ExtractTableOptions>
}
