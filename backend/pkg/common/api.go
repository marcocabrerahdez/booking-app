package common

import (
	"backend/pkg/helpers"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
)

type ResponseStatus string

const (
	Success ResponseStatus = "success"
	Error   ResponseStatus = "error"
)

type ApiResponse[T any] struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    T              `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

type Page[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
	Total    int64 `json:"total"`
	Filtered int64 `json:"filtered"`
}

func NewPage[T any](items []T, page int, size int, total int64, filtered int64) Page[T] {
	return Page[T]{
		Items:    items,
		Page:     page,
		Size:     size,
		Total:    total,
		Filtered: filtered,
	}
}

func NewSuccessResponse[T any](data T, message string) ApiResponse[T] {
	return ApiResponse[T]{
		Status:  Success,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(err error, message string) ApiResponse[any] {
	error := ""
	if err != nil {
		error = err.Error()
	}
	return ApiResponse[any]{
		Status:  Error,
		Message: message,
		Error:   error,
	}
}

func NewValidationErrorResponse(errors []*helpers.ValidationErrors, message string) ApiResponse[any] {
	return ApiResponse[any]{
		Status:  Error,
		Message: message,
		Error:   "Validation failed",
		Data:    errors,
	}
}

func OrderBysFromQuery(c *fiber.Ctx) OrderBys {
	orders := c.Query("orders", "")
	if orders == "" {
		return OrderBys{}
	}
	orderList := strings.Split(orders, ",")

	var orderBys OrderBys
	for _, order := range orderList {
		orderSplit := strings.Split(order, ":")
		orderBy := OrderBy{}
		var field, direction string
		if len(orderSplit) == 1 {
			field = orderSplit[0]
			direction = "asc"
		} else if len(orderSplit) == 2 {
			field = orderSplit[0]
			direction = orderSplit[1]
		} else {
			continue
		}
		orderBy.Field = toSnakePreserveDot(field)
		switch direction {
		case "asc":
			orderBy.Direction = Asc
		case "desc":
			orderBy.Direction = Desc
		default:
			orderBy.Direction = Asc
		}
		orderBys = append(orderBys, orderBy)
	}
	return orderBys
}

func PageableFromQuery(c *fiber.Ctx) (Pageable, error) {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		return Pageable{}, err
	}
	size, err := strconv.Atoi(c.Query("size", "10"))
	if err != nil {
		return Pageable{}, err
	}
	return Pageable{
		Page: page,
		Size: size,
	}, nil
}

func RelationsFromQuery(c *fiber.Ctx) []string {
	relations := c.Query("relations", "")
	if relations == "" {
		return []string{}
	}
	relationList := strings.Split(relations, ",")
	for i, relation := range relationList {
		path := strings.Split(relation, ".")
		for j, part := range path {
			path[j] = toCamelPreserveDot(part)
		}
		relation = strings.Join(path, ".")

		relationList[i] = relation
	}
	return relationList
}

func ConditionsFromQuery(c *fiber.Ctx) SQLConditions {
	filters := c.Query("filters", "")

	if filters == "" {
		return SQLConditions{}
	}

	return parseFilters(filters)

}

// Parse filters will recibe a string with a filter list:
// filter1,filter2,filter3
// Each filter will be either a leaf condition or a composite condition
// A leaf condition will be in the format:
// field:comparator:value
// A composite condition will be in the format:
// type:(filterList)

func parseFilters(filters string) SQLConditions {
	// Split the filters by comma, but ignore commas inside parenthesis
	filterList := splitIgnoreParenthesis(filters, ",")

	conditions := SQLConditions{}

	for _, filter := range filterList {
		condition := parseFilter(filter)
		conditions = append(conditions, condition)
	}

	return conditions
}

func parseFilter(filter string) SQLCondition {
	if strings.Contains(filter, "(") {
		return parseCompositeCondition(filter)
	}
	return parseLeafCondition(filter)
}

func parseCompositeCondition(filter string) SQLCondition {
	// type;(filterList)
	parts := strings.Split(filter, ";")

	typeStr := parts[0]
	rest := strings.Join(parts[1:], ";")
	// Remove the parenthesis
	filterListStr := rest[1 : len(rest)-1]

	filterList := parseFilters(filterListStr)

	return SQLCompositeCondition{
		Type:       SQLCompositor(typeStr),
		Conditions: filterList,
	}
}

func parseLeafCondition(filter string) SQLCondition {
	parts := strings.Split(filter, ";")
	if len(parts) != 3 {
		return SQLLeafCondition{}
	}

	field := parts[0]
	comparatorStr := parts[1]
	value := parts[2]

	var comparator SQLOperator
	switch comparatorStr {
	case "eq":
		comparator = Equal
	case "like":
		comparator = Like
	case "ilike":
		comparator = ILike
	case "gt":
		comparator = GreaterThan
	case "lt":
		comparator = LessThan
	case "gte":
		comparator = GreaterEqualThan
	case "lte":
		comparator = LessEqualThan
	case "in":
		comparator = In
	case "isnull":
		comparator = IsNull
	case "isnotnull":
		comparator = IsNotNull
	default:
		return SQLLeafCondition{}
	}

	return SQLLeafCondition{
		Field:      field,
		Comparator: comparator,
		Value:      value,
	}
}

func toSnakePreserveDot(str string) string {
	parts := strings.Split(str, ".")
	for i, part := range parts {
		parts[i] = strcase.ToSnake(part)
	}
	return strings.Join(parts, ".")
}

func toCamelPreserveDot(str string) string {
	parts := strings.Split(str, ".")
	for i, part := range parts {
		parts[i] = strcase.ToCamel(part)
	}
	return strings.Join(parts, ".")
}

func splitIgnoreParenthesis(str string, sep string) []string {
	parts := strings.Split(str, sep)
	var result []string
	var current string
	parenthesis := 0
	for _, part := range parts {
		current += part
		parenthesis += strings.Count(part, "(") - strings.Count(part, ")")
		if parenthesis == 0 {
			result = append(result, current)
			current = ""
		} else {
			current += sep
		}
	}
	return result
}
