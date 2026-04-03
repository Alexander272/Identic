import { extractTableData } from './extractTable/extractTableData'

export const handleGlobalPaste = (e: ClipboardEvent) => {
	// Проверяем наличие clipboardData
	if (!e.clipboardData) return

	const html = e.clipboardData.getData('text/html')
	console.log(html)

	// Проверяем, есть ли HTML и не было ли событие уже отменено
	if (html && !e.defaultPrevented) {
		// const parser = new DOMParser()
		// const doc = parser.parseFromString(html, 'text/html')
		// const divs = doc.querySelectorAll('div')

		// const lines = Array.from(divs)
		// 	.map(div => (div.textContent || '').replace(/\s+/g, ' ').trim())
		// 	.filter(Boolean)

		const table = extractTableData(html)

		console.log('table', table)

		if (table.length === 0) return

		// Останавливаем оригинальное событие на фазе capture
		e.stopImmediatePropagation()
		e.preventDefault()

		const lines = table.map(item => `${item.name}\t${item.quantity}`)

		const cleanedText = lines.join('\n')

		// Создаем "чистый" контейнер данных
		const dataTransfer = new DataTransfer()
		dataTransfer.setData('text/plain', cleanedText)

		// Генерируем новое событие paste
		const newEvent = new ClipboardEvent('paste', {
			clipboardData: dataTransfer,
			bubbles: true,
			cancelable: true,
		})

		// e.target может быть EventTarget | null, приводим к Node/HTMLElement для dispatch
		if (e.target instanceof Node) {
			e.target.dispatchEvent(newEvent)
		}
	}
}
