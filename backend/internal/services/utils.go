package services

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// Оставляем только буквы, цифры и базовые разделители
	reCleanSymbols = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9\s\.,/№]`)
	reMultiSpace   = regexp.MustCompile(`\s+`)
)

func NormalizeString(name string) string {
	if name == "" {
		return ""
	}

	// 1. Убираем невидимые спецсимволы и Unicode-мусор
	// Это эффективнее, чем TrimSpace, так как убирает \u00A0 (неразрывный пробел)
	name = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, name)

	// 2. Нижний регистр и базовая очистка
	name = strings.ToLower(name)

	// 3. Удаляем запрещенные спецсимволы (оставляем только разрешенные: №, /, -, ., ,)
	name = reCleanSymbols.ReplaceAllString(name, " ")

	// 4. Замена латиницы на кириллицу (самые коварные символы)
	replacer := strings.NewReplacer(
		"a", "а", "e", "е", "o", "о", "p", "р",
		"c", "с", "x", "х", "y", "у", "h", "н", "k", "к",
	)
	name = replacer.Replace(name)

	// 5. Схлопываем лишние пробелы внутри и по краям
	name = reMultiSpace.ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)

	return name
}
