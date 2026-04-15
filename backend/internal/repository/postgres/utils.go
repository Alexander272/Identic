package postgres

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type QueryBuilder struct {
	baseQuery string
	where     []string
	args      []interface{}
	groupBy   string
	orderBy   string
	limit     int
	cursor    interface{}
}

func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery: baseQuery,
		where:     []string{},
		args:      []interface{}{},
	}
}

func (qb *QueryBuilder) AddSearch(field, search string) {
	if search == "" {
		return
	}
	// ILIKE для PostgreSQL, безопасный placeholder
	qb.args = append(qb.args, "%"+search+"%")
	qb.where = append(qb.where, fmt.Sprintf("%s ILIKE $%d", field, len(qb.args)))
}
func (qb *QueryBuilder) AddMultiSearch(fields []string, search string) {
	if search == "" || len(fields) == 0 {
		return
	}

	parts := []string{}
	for _, f := range fields {
		qb.args = append(qb.args, "%"+search+"%")
		parts = append(parts, fmt.Sprintf("%s ILIKE $%d", f, len(qb.args)))
	}

	qb.where = append(qb.where, "("+strings.Join(parts, " OR ")+")")
}

func (qb *QueryBuilder) SetGroupBy(field string) {
	qb.groupBy = field
}

func (qb *QueryBuilder) SetSort(field string, desc bool) {
	order := field
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}
	qb.orderBy = order
}

type SortField struct {
	Field string
	Desc  bool
}

func (qb *QueryBuilder) SetMultiSort(fields []SortField) {
	if len(fields) == 0 {
		return
	}

	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		dir := "ASC"
		if f.Desc {
			dir = "DESC"
		}
		parts = append(parts, fmt.Sprintf("%s %s", f.Field, dir))
	}

	qb.orderBy = strings.Join(parts, ", ")
}
func (qb *QueryBuilder) SetMultiCompositeCursor(fields []string, values []interface{}, descSlice []bool) {
	if len(fields) != len(values) || len(fields) != len(descSlice) {
		return // или panic в дев-режиме
	}

	// Построение условия: (f1 < v1) OR (f1 = v1 AND f2 < v2) OR ...
	var orConditions []string

	for i := range fields {
		var andParts []string

		// Равенство по всем предыдущим полям
		for j := 0; j < i; j++ {
			andParts = append(andParts, fmt.Sprintf("%s = $%d", fields[j], j+1))
		}

		// Сравнение по текущему полю
		op := "<"
		if !descSlice[i] {
			op = ">"
		}
		andParts = append(andParts, fmt.Sprintf("%s %s $%d", fields[i], op, i+1))

		orConditions = append(orConditions, "("+strings.Join(andParts, " AND ")+")")
	}

	qb.where = append(qb.where, "("+strings.Join(orConditions, ") OR (")+")")

	// Добавляем аргументы
	qb.args = append(qb.args, values...)
}

func (qb *QueryBuilder) SetCursor(field string, cursor interface{}, desc bool) {
	if cursor == nil {
		return
	}

	qb.cursor = cursor
	op := ">"
	if desc {
		op = "<"
	}

	qb.args = append(qb.args, cursor)
	qb.where = append(qb.where, fmt.Sprintf("%s %s $%d", field, op, len(qb.args)))
}
func (qb *QueryBuilder) SetCompositeCursor(field string, dateVal string, idVal string, desc bool) {
	op := ">"
	if desc {
		op = "<"
	}
	qb.args = append(qb.args, dateVal, idVal)
	qb.where = append(qb.where,
		fmt.Sprintf("(%s, o.id) %s ($%d::timestamp, $%d::uuid)", field, op, len(qb.args)-1, len(qb.args)),
	)
}

func (qb *QueryBuilder) SetLimit(limit int) {
	qb.limit = limit
}

func (qb *QueryBuilder) AddCompositeFilter(fieldNames []string, compare []string, values []string) {
	if len(fieldNames) != len(compare) || len(fieldNames) != len(values) {
		return
	}

	parts := []string{}

	for i := range fieldNames {
		filter, args := qb.prepareFilter(fieldNames[i], compare[i], values[i])
		parts = append(parts, filter)
		if args != nil {
			qb.args = append(qb.args, args)
		}
	}
	qb.where = append(qb.where, "("+strings.Join(parts, " OR ")+")")
}
func (qb *QueryBuilder) AddFilter(fieldName string, compare string, value string) {
	if compare == "" {
		return
	}

	filter, args := qb.prepareFilter(fieldName, compare, value)
	qb.where = append(qb.where, filter)
	if args != nil {
		qb.args = append(qb.args, args)
	}

	// switch compare {
	// case "con":
	// 	qb.args = append(qb.args, "%"+value+"%")
	// 	qb.where = append(qb.where, fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)))
	// case "start":
	// 	qb.args = append(qb.args, value+"%")
	// 	qb.where = append(qb.where, fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)))
	// case "end":
	// 	qb.args = append(qb.args, "%"+value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)))
	// case "like":
	// 	qb.args = append(qb.args, value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)))
	// case "nlike":
	// 	qb.args = append(qb.args, value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s NOT ILIKE $%d", fieldName, len(qb.args)))
	// case "in":
	// 	// Разделяем строку по разделителю
	// 	vals := strings.Split(value, "|")
	// 	qb.args = append(qb.args, vals)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s::text ILIKE ANY ($%d)", fieldName, len(qb.args)))
	// case "nin":
	// 	vals := strings.Split(value, "|")
	// 	qb.args = append(qb.args, vals)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s::text NOT ILIKE ANY ($%d)", fieldName, len(qb.args)))
	// case "eq":
	// 	qb.args = append(qb.args, value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s = $%d", fieldName, len(qb.args)))
	// case "neq":
	// 	qb.args = append(qb.args, value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s != $%d", fieldName, len(qb.args)))
	// case "gte":
	// 	qb.args = append(qb.args, value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s >= $%d", fieldName, len(qb.args)))
	// case "lte":
	// 	qb.args = append(qb.args, value)
	// 	qb.where = append(qb.where, fmt.Sprintf("%s <= $%d", fieldName, len(qb.args)))
	// case "null":
	// 	qb.where = append(qb.where, fmt.Sprintf("(%s IS NULL OR %s::text = '')", fieldName, fieldName))
	// }
}
func (qb *QueryBuilder) prepareFilter(fieldName string, compare string, value string) (string, interface{}) {
	switch compare {
	case "con":
		return fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)+1), "%" + value + "%"
	case "start":
		return fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)+1), value + "%"
	case "end":
		return fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)+1), "%" + value
	case "like":
		return fmt.Sprintf("%s ILIKE $%d", fieldName, len(qb.args)+1), value
	case "nlike":
		return fmt.Sprintf("%s NOT ILIKE $%d", fieldName, len(qb.args)+1), value
	case "in":
		// Разделяем строку по разделителю
		return fmt.Sprintf("%s::text ILIKE ANY ($%d)", fieldName, len(qb.args)+1), strings.Split(value, "|")
	case "nin":
		return fmt.Sprintf("%s::text NOT ILIKE ANY ($%d)", fieldName, len(qb.args)+1), strings.Split(value, "|")
	case "eq":
		return fmt.Sprintf("%s = $%d", fieldName, len(qb.args)+1), value
	case "neq":
		return fmt.Sprintf("%s != $%d", fieldName, len(qb.args)+1), value
	case "gte":
		return fmt.Sprintf("%s >= $%d", fieldName, len(qb.args)+1), value
	case "lte":
		return fmt.Sprintf("%s <= $%d", fieldName, len(qb.args)+1), value
	case "null":
		return fmt.Sprintf("(%s IS NULL OR %s::text = '')", fieldName, fieldName), nil
	}

	return "", nil
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := qb.baseQuery

	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}

	if qb.groupBy != "" {
		query += " GROUP BY " + qb.groupBy
	}

	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	return query, qb.args
}

// CursorState хранит значения всех полей сортировки для пагинации
type CursorState struct {
	Values []interface{} `json:"v"` // значения полей в порядке сортировки
	Types  []string      `json:"t"` // типы полей: "time", "uuid", "int", "string"
	Desc   []bool        `json:"d"` // направление сортировки для каждого поля
}

// Encode кодирует курсор в строку (безопасно для URL)
func (c CursorState) Encode() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to encode cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// BuildCursorFromRow создаёт курсор из последней строки результата
// rowValues — значения полей в том же порядке, что и в сортировке
// fieldTypes — типы полей для корректного форматирования
// desc — направление сортировки
func BuildCursorFromRow(rowValues []interface{}, fieldTypes []string, desc []bool) (string, error) {
	cursor := CursorState{
		Values: make([]interface{}, 0, len(rowValues)),
		Types:  fieldTypes,
		Desc:   desc,
	}

	// Копируем значения, приводя time.Time к RFC3339 для консистентности
	for i, val := range rowValues {
		switch v := val.(type) {
		case time.Time:
			cursor.Values = append(cursor.Values, v.Format(time.RFC3339Nano))
		case *time.Time:
			if v != nil {
				cursor.Values = append(cursor.Values, v.Format(time.RFC3339Nano))
			} else {
				cursor.Values = append(cursor.Values, nil)
			}
		default:
			cursor.Values = append(cursor.Values, v)
		}
		_ = i // suppress unused
	}

	return cursor.Encode()
}

// DecodeCursor парсит строку курсора обратно в структуру
func DecodeCursor(cursor string) (*CursorState, error) {
	if cursor == "" {
		return nil, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}
	var c CursorState
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}
	return &c, nil
}

// ParseCursorValues парсит значения курсора с учётом типов
// Возвращает значения, готовые к подстановке в запрос
func (c *CursorState) ParseCursorValues() ([]interface{}, error) {
	if len(c.Values) != len(c.Types) {
		return nil, fmt.Errorf("cursor values/types mismatch")
	}

	result := make([]interface{}, 0, len(c.Values))
	for i, val := range c.Values {
		switch c.Types[i] {
		case "time", "timestamp":
			if str, ok := val.(string); ok {
				t, err := time.Parse(time.RFC3339Nano, str)
				if err != nil {
					t, err = time.Parse(time.RFC3339, str) // fallback
					if err != nil {
						return nil, fmt.Errorf("invalid time cursor value: %w", err)
					}
				}
				result = append(result, t)
			} else {
				result = append(result, val)
			}
		case "uuid":
			// pgx сам сконвертирует string → uuid, если тип в БД uuid
			result = append(result, val)
		case "int", "int64":
			switch v := val.(type) {
			case float64: // JSON decodes numbers as float64
				result = append(result, int64(v))
			default:
				result = append(result, v)
			}
		case "float32", "float64":
			result = append(result, val)
		default:
			result = append(result, val)
		}
	}
	return result, nil
}
