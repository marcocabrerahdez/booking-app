package common

import "github.com/spf13/viper"

type SQLOperator string

const (
	Like             SQLOperator = "like"
	Equal            SQLOperator = "="
	ILike            SQLOperator = "ilike"
	GreaterThan      SQLOperator = ">"
	LessThan         SQLOperator = "<"
	GreaterEqualThan SQLOperator = ">="
	LessEqualThan    SQLOperator = "<="
	In               SQLOperator = "in"
	IsNull           SQLOperator = "is null"
	IsNotNull        SQLOperator = "is not null"
)

type SQLCompositor string

const (
	And SQLCompositor = "and"
	Or  SQLCompositor = "or"
	Not SQLCompositor = "not"
)

type SQLLeafCondition struct {
	Field      string
	Value      string
	Comparator SQLOperator
}

type SQLCompositeCondition struct {
	Type       SQLCompositor
	Conditions []SQLCondition
}

type SQLCondition interface {
	IsComposite() bool
}

func (c SQLLeafCondition) IsComposite() bool {
	return false
}

func (c SQLCompositeCondition) IsComposite() bool {
	return true
}

type SQLConditions []SQLCondition

var NoConditions = SQLConditions{}

type OrderDirection string

const (
	Asc  OrderDirection = "asc"
	Desc OrderDirection = "desc"
)

type OrderBy struct {
	Field     string
	Direction OrderDirection
}

type OrderBys []OrderBy

var NoOrder = OrderBys{}

var NoFields = []string{}

type Pageable struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func PageableFrom(page int, size int) Pageable {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = viper.GetInt("services.pagination.size.default")
	}
	if size > viper.GetInt("services.pagination.size.max") {
		size = viper.GetInt("services.pagination.size.default")
	}

	return Pageable{
		Page: page,
		Size: size,
	}
}
