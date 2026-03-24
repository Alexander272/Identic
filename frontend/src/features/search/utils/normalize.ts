export const normalize = (str: string) => {
	if (!str) return ''
	return str
		.replace(/\s+/g, ' ') // Заменяет любые пробелы/переносы/табы на ОДИН пробел
		.trim() // Убирает пробелы по краям
}
