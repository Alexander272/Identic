package services

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
)

var (
	// Оставляем только буквы, цифры и базовые разделители
	// reCleanSymbols = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9\s\.,/№]`)
	// reCleanSymbols = regexp.MustCompile(`[^a-zа-я0-9]+`)
	reCleanSymbols = regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]+`)
	reMultiSpace   = regexp.MustCompile(`\s+`)
	// reStd          = regexp.MustCompile(`(?:гост|ост|ту)[\s\-]*[\d\.\-]+`)
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

	// Удаляем стандарты (гост, ост, ту)
	// name = reStd.ReplaceAllString(name, "")

	//// 3. Удаляем запрещенные спецсимволы (оставляем только разрешенные: №, /, -, ., ,)
	//3. Удаляем все символы кроме букв, цифр
	name = reCleanSymbols.ReplaceAllString(name, " ")

	// // 4. Замена латиницы на кириллицу (самые коварные символы)
	// replacer := strings.NewReplacer(
	// 	"a", "а", "e", "е", "o", "о", "p", "р",
	// 	"c", "с", "x", "х", "y", "у", "h", "н", "k", "к",
	// )
	// name = replacer.Replace(name)

	// 5. Схлопываем лишние пробелы внутри и по краям
	name = reMultiSpace.ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)

	return name
}

func CalculateHash(positions []*models.PositionDTO) string {
	if len(positions) == 0 {
		return ""
	}

	// 1. Агрегируем (на случай дублей в строках)
	totals := make(map[string]float64)
	for _, p := range positions {
		name := strings.ToLower(strings.TrimSpace(p.Name))
		totals[name] += float64(p.Quantity)
	}

	// 2. Выгружаем в слайс для сортировки (важно для консистентности хеша)
	keys := make([]string, 0, len(totals))
	for k := range totals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. Собираем строку: "товар1:10.000;товар2:5.500"
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte(':')
		sb.WriteString(strconv.FormatFloat(totals[k], 'f', 3, 64))
		sb.WriteByte(';')
	}

	// 4. Считаем SHA256
	hash := sha256.Sum256([]byte(sb.String()))
	return fmt.Sprintf("%x", hash)
}

func splitPositions(orderId string, dto []*models.PositionDTO) (created, updated, deleted, notChanged []*models.PositionDTO) {
	for i := range dto {
		switch dto[i].Status {
		case models.PositionCreated:
			dto[i].Id = uuid.NewString()
			dto[i].OrderId = orderId
			created = append(created, dto[i])
		case models.PositionUpdated:
			updated = append(updated, dto[i])
		case models.PositionDeleted:
			deleted = append(deleted, dto[i])
		default:
			notChanged = append(notChanged, dto[i])
		}
	}

	return
}

func ParseUserAgent(ua string) models.LoginMetadataDTO {
	meta := models.LoginMetadataDTO{
		Success: true,
	}

	if ua == "" {
		meta.IsDesktop = true
		meta.Device = "unknown"
		meta.Browser = "Unknown"
		meta.OS = "Unknown"
		return meta
	}

	uaParsed := useragent.Parse(ua)

	meta.Browser = uaParsed.Name
	meta.BrowserVersion = uaParsed.Version
	meta.OS = uaParsed.OS
	meta.OSVersion = uaParsed.OSVersion

	meta.IsMobile = uaParsed.Mobile
	meta.IsTablet = uaParsed.Tablet
	meta.IsBot = uaParsed.Bot
	meta.IsDesktop = uaParsed.Desktop

	if meta.IsMobile {
		meta.Device = "mobile"
	} else if meta.IsTablet {
		meta.Device = "tablet"
	} else if meta.IsBot {
		meta.Device = "bot"
	} else if meta.IsDesktop {
		meta.Device = "desktop"
	} else {
		meta.Device = "unknown"
	}

	if meta.Browser == "" {
		meta.Browser = "Unknown"
	}
	if meta.OS == "" {
		meta.OS = "Unknown"
	}

	return meta
}
