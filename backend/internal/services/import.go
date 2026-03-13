package services

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/Identic/backend/internal/constants"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/repository/postgres"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
)

type ImportService struct {
	txManager TransactionManager
	orders    Orders
	positions Positions
}

func NewImportService(txManager TransactionManager, orders Orders, positions Positions) *ImportService {
	return &ImportService{
		txManager: txManager,
		orders:    orders,
		positions: positions,
	}
}

type Import interface {
	Load(ctx context.Context, dto *models.ImportDTO) error
}

func (s *ImportService) Load(ctx context.Context, dto *models.ImportDTO) error {
	return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
		file, err := dto.File.Open()
		if err != nil {
			return fmt.Errorf("failed to open file. error: %w", err)
		}
		defer file.Close()

		excel, err := excelize.OpenReader(file)
		if err != nil {
			return fmt.Errorf("failed to open excel file. error: %w", err)
		}
		defer excel.Close()

		sheets := excel.GetSheetList()

		// Канал для передачи готовых пачек заказов от парсеров к БД
		// Буфер равен количеству листов, чтобы парсеры не ждали БД
		results := make(chan []*models.OrderDTO, len(sheets))

		// Используем errgroup для контроля горутин
		g, gCtx := errgroup.WithContext(ctx)

		for _, sheet := range sheets {
			if !strings.Contains(sheet, "заяв") {
				continue
			}

			sName := sheet
			g.Go(func() error {
				// parseSheet только читает данные и возвращает слайс
				orders, err := s.parseSheet(gCtx, excel, sName)
				if err != nil {
					return fmt.Errorf("sheet %s: %w", sName, err)
				}

				// Отправляем результат в канал
				select {
				case results <- orders:
					return nil
				case <-gCtx.Done():
					return gCtx.Err()
				}
			})

			// if err := s.loadSheet(ctx, tx, sheet, excel); err != nil {
			// 	return err
			// }
		}

		// Закрываем канал, когда ВСЕ парсеры закончат работу
		go func() {
			g.Wait()
			close(results)
		}()

		// ЗАПИСЬ В БД (Consumer)
		// Транзакция открывается здесь, чтобы последовательно влить данные
		return s.txManager.WithinTransaction(ctx, func(tx postgres.Tx) error {
			for batch := range results {
				if len(batch) == 0 {
					continue
				}

				// Здесь важно, чтобы CreateSeveral делал Bulk Insert
				if err := s.orders.CreateSeveral(ctx, tx, batch); err != nil {
					return fmt.Errorf("db save failed: %w", err)
				}
			}

			// Если хоть один парсер вернул ошибку, errgroup её вернет,
			// и WithinTransaction откатит всё назад.
			return g.Wait()
		})
	})
}

func (s *ImportService) parseSheet(ctx context.Context, excel *excelize.File, sheet string) ([]*models.OrderDTO, error) {
	rows, err := excel.Rows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows. error: %w", err)
	}
	defer rows.Close()

	template := constants.ImportTemplate
	ordersBuffer := make([]*models.OrderDTO, 0, constants.MaxOrdersInBatch)
	var currentOrder *models.OrderDTO
	row := make([]string, template.Count)
	totalPositionsInBuffer := 0 // Количество позиций в буфере
	rowNum := 0
	rowNumInOrder := 0

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		rowNum++
		rowNumInOrder++

		origRow, err := rows.Columns(excelize.Options{RawCellValue: true})
		if err != nil {
			return nil, fmt.Errorf("failed to get columns. error: %w", err)
		}
		for i := range row {
			row[i] = ""
		}
		copy(row, origRow)

		if row[template.NameColumn] == "" || row[template.NameColumn] == "Наименование продукции" {
			continue
		}

		tmp, err := s.parseRow(row, template, rowNum, rowNumInOrder)
		if err != nil {
			return nil, err
		}

		isNewOrder := false
		if currentOrder == nil {
			isNewOrder = true
			// } else if tmp.order.Bill != "" && tmp.order.Bill != currentOrder.Bill {
			// 	isNewOrder = true
			// } else if tmp.order.Bill == "" && (tmp.order.Customer != currentOrder.Customer || tmp.order.Consumer != currentOrder.Consumer) {
		} else {
			sameCustomer := false
			if tmp.order.Customer != "" && currentOrder.Customer != tmp.order.Customer {
				sameCustomer = true
			}
			sameConsumer := false
			if tmp.order.Consumer != "" && currentOrder.Consumer != tmp.order.Consumer {
				sameConsumer = true
			}

			if sameCustomer || sameConsumer {
				isNewOrder = true
			}
		}

		if isNewOrder {
			rowNumInOrder = 1
			tmp.position.RowNumber = rowNumInOrder

			if currentOrder != nil {
				totalPositionsInBuffer += len(currentOrder.Positions)
				ordersBuffer = append(ordersBuffer, currentOrder)
			}
			currentOrder = &models.OrderDTO{
				Customer:  tmp.order.Customer,
				Consumer:  tmp.order.Consumer,
				Manager:   strings.TrimSpace(row[template.ManagerColumn]),
				Bill:      tmp.order.Bill,
				Date:      tmp.order.Date,
				Positions: []*models.PositionDTO{tmp.position},
			}
		} else {
			currentOrder.Positions = append(currentOrder.Positions, tmp.position)
		}
	}

	return ordersBuffer, nil
}

// ! Deprecated
func (s *ImportService) loadSheet(ctx context.Context, tx postgres.Tx, sheet string, excel *excelize.File) error {
	rows, err := excel.Rows(sheet)
	if err != nil {
		return fmt.Errorf("failed to get rows. error: %w", err)
	}
	defer rows.Close()

	template := constants.ImportTemplate
	ordersBuffer := make([]*models.OrderDTO, 0, constants.MaxOrdersInBatch)
	var currentOrder *models.OrderDTO
	row := make([]string, template.Count)
	totalPositionsInBuffer := 0 // Количество позиций в буфере
	rowNum := 0
	skippedRows := 0
	ordersCreated := 0
	positionsCreated := 0

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return err
		}
		rowNum++

		origRow, err := rows.Columns(excelize.Options{RawCellValue: true})
		if err != nil {
			return fmt.Errorf("failed to get columns. error: %w", err)
		}
		for i := range row {
			row[i] = ""
		}
		copy(row, origRow)

		if row[template.NameColumn] == "" || row[template.NameColumn] == "Наименование продукции" {
			skippedRows++
			continue
		}

		// dateStr := strings.TrimSpace(row[template.DateColumn])
		// date, err := time.Parse("02/01/2006", dateStr)
		// if err != nil {
		// 	logger.Debug("import", logger.StringAttr("date", dateStr))
		// 	return fmt.Errorf("row %d: invalid date format: %w", rowNum, err)
		// }
		// qtyStr := strings.TrimSpace(row[template.QuantityColumn])
		// quantityFloat, err := strconv.ParseFloat(qtyStr, 64)
		// if err != nil {
		// 	return fmt.Errorf("row %d: invalid quantity: %w", rowNum, err)
		// }
		// quantity := int(quantityFloat)

		// customer := strings.TrimSpace(row[template.CustomerColumn])
		// consumer := strings.TrimSpace(row[template.ConsumerColumn])
		// bill := strings.TrimSpace(row[template.BillColumn])

		// position := &models.PositionDTO{
		// 	Name:     strings.TrimSpace(row[template.NameColumn]),
		// 	Quantity: quantity,
		// 	Notes:    strings.TrimSpace(row[template.NotesColumn]),
		// }

		tmp, err := s.parseRow(row, template, rowNum, 0)
		if err != nil {
			return err
		}

		isNewOrder := false
		if currentOrder == nil {
			isNewOrder = true
			// } else if tmp.order.Bill != "" && tmp.order.Bill != currentOrder.Bill {
			// 	isNewOrder = true
			// } else if tmp.order.Bill == "" && (tmp.order.Customer != currentOrder.Customer || tmp.order.Consumer != currentOrder.Consumer) {
		} else if tmp.order.Customer != currentOrder.Customer || tmp.order.Consumer != currentOrder.Consumer {
			isNewOrder = true
		}

		if isNewOrder {
			ordersCreated++
			positionsCreated++
			if currentOrder != nil {
				totalPositionsInBuffer += len(currentOrder.Positions)
				ordersBuffer = append(ordersBuffer, currentOrder)
			}
			currentOrder = &models.OrderDTO{
				Customer:  tmp.order.Customer,
				Consumer:  tmp.order.Consumer,
				Manager:   strings.TrimSpace(row[template.ManagerColumn]),
				Bill:      tmp.order.Bill,
				Date:      tmp.order.Date,
				Positions: []*models.PositionDTO{tmp.position},
			}
		} else {
			positionsCreated++
			currentOrder.Positions = append(currentOrder.Positions, tmp.position)
		}

		if len(ordersBuffer) >= constants.MaxOrdersInBatch || totalPositionsInBuffer >= constants.MaxPositionsInBatch {
			if err := s.orders.CreateSeveral(ctx, tx, ordersBuffer); err != nil {
				return fmt.Errorf("create batch: %w", err)
			}
			for i := range ordersBuffer {
				ordersBuffer[i] = nil
			}

			ordersBuffer = ordersBuffer[:0] // Очистка без аллокации
		}
	}

	logger.Debug("import",
		logger.IntAttr("rows", rowNum),
		logger.IntAttr("skipped", skippedRows),
		logger.IntAttr("orders", ordersCreated),
		logger.IntAttr("positions", positionsCreated))
	if currentOrder != nil {
		ordersBuffer = append(ordersBuffer, currentOrder)
	}
	if len(ordersBuffer) > 0 {
		if err := s.orders.CreateSeveral(ctx, tx, ordersBuffer); err != nil {
			return fmt.Errorf("create final batch: %w", err)
		}
	}
	return nil
}

type rowData struct {
	order    *models.OrderDTO
	position *models.PositionDTO
}

var reNumbers = regexp.MustCompile(`[0-9,.]+`)
var reNum = regexp.MustCompile(`\d+(?:[.,]\d+)?`)
var reCount = regexp.MustCompile(`^\d+[.)]\s*`)
var reWithUnit = regexp.MustCompile(`(?i)(\d+)\s*(?:шт|кг)|(?:шт|кг)\.?\s*(\d+)`)
var reStandards = regexp.MustCompile(`(?i)(ОСТ|ГОСТ|ТУ|ASME|B)\s*[\d\.\-]+$`)
var reEndDigits = regexp.MustCompile(`\s+(\d+)$`)
var reSpace = regexp.MustCompile(`\s+`)

func (s *ImportService) parseRow(row []string, t *models.ImportTemplate, rowNum, rowInOrder int) (*rowData, error) {
	date := time.Time{}
	dateStr := strings.TrimSpace(row[t.DateColumn])
	if dateStr != "" {
		// var err error
		// dateRe := regexp.MustCompile(`\d{2}\/\d{2}\/\d{4}`)
		// dateTemplate := "02/01/2006"

		// if !dateRe.MatchString(dateStr) {
		// 	dateTemplate = "02/01/06"
		// }

		// date, err = time.Parse(dateTemplate, dateStr)
		dateNum, err := strconv.ParseFloat(dateStr, 64)
		if err != nil {
			logger.Debug("import", logger.StringAttr("date", dateStr))
			return nil, fmt.Errorf("row %d: invalid date: %w", rowNum, err)
		}

		date, err = excelize.ExcelDateToTime(dateNum, false)
		if err != nil {
			logger.Debug("import", logger.StringAttr("date", dateStr))
			return nil, fmt.Errorf("row %d: invalid date: %w", rowNum, err)
		}
	}
	name := strings.TrimSpace(row[t.NameColumn])
	var quantity float64 = 0

	match := reNumbers.FindString(row[t.QuantityColumn])
	match = strings.ReplaceAll(match, ",", ".")

	if match == "" {
		name, quantity = s.extractQuantity(row[t.NameColumn])
	} else {
		var err error
		quantity, err = strconv.ParseFloat(match, 64)
		if err != nil {
			logger.Debug("import", logger.StringAttr("quantity", match))
			return nil, fmt.Errorf("row %d: invalid quantity: %w", rowNum, err)
		}
	}
	search := NormalizeString(name)

	return &rowData{
		order: &models.OrderDTO{
			Customer: strings.TrimSpace(row[t.CustomerColumn]),
			Consumer: strings.TrimSpace(row[t.ConsumerColumn]),
			Manager:  strings.TrimSpace(row[t.ManagerColumn]),
			Bill:     strings.TrimSpace(row[t.BillColumn]),
			Date:     date,
		},
		position: &models.PositionDTO{
			RowNumber: rowInOrder,
			Name:      name,
			Search:    search,
			Quantity:  quantity,
			Notes:     strings.TrimSpace(row[t.NotesColumn]),
		},
	}, nil
}

func (s *ImportService) parseLine(line string) (name string, quantity float64) {
	line = strings.TrimSpace(line)

	// 1. Находим все вхождения чисел в строке
	matches := reNum.FindAllString(line, -1)
	indices := reNum.FindAllStringIndex(line, -1)

	if len(matches) > 0 {
		// Берем ПОСЛЕДНЕЕ число из найденных
		lastIdx := len(matches) - 1
		rawVal := matches[lastIdx]
		pos := indices[lastIdx]

		// Конвертируем в число
		cleanVal := strings.ReplaceAll(rawVal, ",", ".")
		quantity, _ = strconv.ParseFloat(cleanVal, 64)

		// Название — это всё, что ДО начала последнего числа
		name = line[:pos[0]]

		// Чистим хвост названия от мусора: "шт", "ШТ", "шт.", "-", пробелы
		name = strings.TrimRight(name, " -–штШТ. ")

		// (Опционально) Убираем порядковый номер в начале: "2. "
		name = reCount.ReplaceAllString(name, "")
	} else {
		name = line
	}

	return name, quantity
}

func (s *ImportService) extractQuantity(input string) (string, float64) {
	// 1. Очищаем строку от лишних пробелов по краям
	line := strings.TrimSpace(input)

	// Убираем порядковый номер в начале: "2. "
	line = reCount.ReplaceAllString(line, "")

	// Регулярное выражение для поиска количества.
	// Оно ищет число (\d+), которое:
	// - Либо стоит перед/после "шт" (с точкой или без)
	// - Либо стоит в конце строки или перед стандартом

	// Вариант А: Ищем паттерны с "шт"
	matches := reWithUnit.FindStringSubmatch(line)

	if len(matches) > 0 {
		qtyStr := ""
		if matches[1] != "" {
			qtyStr = matches[1]
		} else {
			qtyStr = matches[2]
		}
		qty, _ := strconv.ParseFloat(qtyStr, 64)

		// Удаляем найденное количество из названия
		name := reWithUnit.ReplaceAllString(line, "")
		name = reSpace.ReplaceAllString(name, " ")
		return strings.TrimSpace(name), qty
	}

	// Вариант Б: Если "шт" нет, ищем число в конце, которое не похоже на ГОСТ/ОСТ
	// Проверяем, не заканчивается ли строка на стандарт
	isStandard := reStandards.MatchString(line)

	if !isStandard {
		if m := reEndDigits.FindStringSubmatch(line); m != nil {
			qty, _ := strconv.ParseFloat(m[1], 64)
			name := reEndDigits.ReplaceAllString(line, "")
			name = reSpace.ReplaceAllString(name, " ")
			return strings.TrimSpace(name), qty
		}
	}
	return line, 0
}
