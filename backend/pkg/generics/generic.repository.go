package generics

import (
	"backend/database"
	"backend/pkg/common"

	"fmt"
	"log"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository[Entity common.Entity, DTO common.DTO] interface {
	Create(payload Entity) (Entity, error)
	Update(payload Entity) (Entity, error)
	Delete(payload Entity) error
	FindOne(id uuid.UUID, relations []string) (Entity, error)
	FindAll(pageable common.Pageable, conditions common.SQLConditions, relations []string, orderBys common.OrderBys) (*common.Page[Entity], error)
	FindOneRandom() (Entity, error)
	Exists(id uuid.UUID) (bool, error)
	Count(conditions common.SQLConditions) (int64, error)
	HardDelete(payload Entity) error
	GetOneDeleted(id uuid.UUID) (Entity, error)
	GetDeleted(pageable common.Pageable, conditions common.SQLConditions, relations []string, orderBys common.OrderBys) (*common.Page[Entity], error)
}

type GenericRepository[Entity common.Entity, DTO common.DTO] struct{}

func NewGenericRepositoryGORM[Entity common.Entity, DTO common.DTO]() GenericRepository[Entity, DTO] {
	var service GenericRepository[Entity, DTO]
	return service
}

func (imp GenericRepository[Entity, DTO]) Create(payload Entity) (Entity, error) {
	err := database.DB.Create(&payload).Error
	return payload, err
}

func (imp GenericRepository[Entity, DTO]) Update(payload Entity) (Entity, error) {
	if payload.GetID() == uuid.Nil {
		return payload, fmt.Errorf("ID cannot be nil")
	}
	err := database.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Omit("created_at").
		Save(&payload).Error
	return payload, err
}

func (imp GenericRepository[Entity, DTO]) Delete(payload Entity) (err error) {
	err = database.DB.Delete(&payload, payload.GetID()).Error
	if err != nil {
		log.Printf("Error deleting %s: %s", payload.GetID(), err)
	}
	return
}

func (imp GenericRepository[Entity, DTO]) FindOne(id uuid.UUID, relations []string) (Entity, error) {
	var entity Entity
	err := database.DB.
		Scopes(
			Preload(relations),
		).
		First(&entity, "id = ?", id).Error
	return entity, err
}

func (imp GenericRepository[Entity, DTO]) FindOneRandom() (Entity, error) {
	var entity Entity
	err := database.DB.Order("RANDOM()").First(&entity).Error
	return entity, err
}

func (imp GenericRepository[Entity, DTO]) FindAll(pageable common.Pageable, conditions common.SQLConditions, relations []string, orderBys common.OrderBys) (*common.Page[Entity], error) {
	var entities []Entity
	limit := pageable.Size
	offset := (pageable.Page - 1) * pageable.Size

	var count int64
	database.DB.Model(&entities).Count(&count)

	result := database.DB.
		Limit(limit).
		Offset(offset).
		Scopes(
			Preload(relations),
			Filters(conditions),
			Order(orderBys),
		).
		Find(&entities)

		// Debug query
	database.DB.Debug().
		Limit(limit).
		Offset(offset).
		Scopes(
			Preload(relations),
			Filters(conditions),
			Order(orderBys),
		).
		Find(&entities)

	var filtered int64
	database.DB.Model(&entities).
		Scopes(
			Preload(relations),
			Filters(conditions),
		).
		Count(&filtered)

	return &common.Page[Entity]{
		Items:    entities,
		Page:     pageable.Page,
		Size:     pageable.Size,
		Total:    count,
		Filtered: filtered,
	}, result.Error
}

func (imp GenericRepository[Entity, DTO]) Exists(id uuid.UUID) (bool, error) {
	var entity Entity
	err := database.DB.First(&entity, "id = ?", id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (imp GenericRepository[Entity, DTO]) Count(conditions common.SQLConditions) (int64, error) {
	var count int64
	var entity Entity
	err := database.DB.
		Model(entity).
		Scopes(
			Filters(conditions),
		).
		Count(&count).Error
	return count, err
}

func (imp GenericRepository[Entity, DTO]) HardDelete(payload Entity) error {
	return database.DB.Unscoped().Delete(&payload).Error
}

func (imp GenericRepository[Entity, DTO]) GetOneDeleted(id uuid.UUID) (Entity, error) {
	var entity Entity
	err := database.DB.Unscoped().First(&entity, "id = ?", id).Error
	return entity, err
}

func (imp GenericRepository[Entity, DTO]) GetDeleted(pageable common.Pageable, conditions common.SQLConditions, relations []string, orderBys common.OrderBys) (*common.Page[Entity], error) {
	var entities []Entity
	limit := pageable.Size
	offset := (pageable.Page - 1) * pageable.Size

	var count int64
	database.DB.Unscoped().Model(&entities).Count(&count)

	result := database.DB.
		Unscoped().
		Limit(limit).
		Offset(offset).
		Scopes(
			Preload(relations),
			Filters(conditions),
			Order(orderBys),
		).
		Find(&entities)

	var filtered int64
	database.DB.Model(&entities).
		Scopes(
			Preload(relations),
			Filters(conditions),
		).
		Count(&filtered)

	return &common.Page[Entity]{
		Items:    entities,
		Page:     pageable.Page,
		Size:     pageable.Size,
		Total:    count,
		Filtered: filtered,
	}, result.Error
}

// Filters returns a function that applies the given conditions to a gorm.DB
func Filters(conditions common.SQLConditions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, condition := range conditions {
			db = filter(condition, common.And)(db)
		}
		return db
	}
}

func filter(condition common.SQLCondition, compositor common.SQLCompositor) func(db *gorm.DB) *gorm.DB {
	if condition.IsComposite() {
		return composition(condition.(common.SQLCompositeCondition))
	}
	return LeaftFilter(condition.(common.SQLLeafCondition), compositor)
}

func composition(condition common.SQLCompositeCondition) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		composite := database.DB
		for _, c := range condition.Conditions {
			composite = composite.Scopes(filter(c, condition.Type))
		}
		switch condition.Type {
		case common.And:
			return db.Where(composite)
		case common.Or:
			return db.Or(composite)
		case common.Not:
			return db.Not(composite)
		default:
			return db.Where(composite)
		}
	}

}

// AndFilterComposition returns a function that applies the given condition to a gorm.DB
func AndFilterComposition(condition common.SQLCondition) func(db *gorm.DB) *gorm.DB {
	if condition.IsComposite() {
		composite := condition.(common.SQLCompositeCondition)

		return func(db *gorm.DB) *gorm.DB {

			compositeQuery := db

			for _, c := range composite.Conditions {
				compositeQuery = compositeQuery.Scopes(AndFilterComposition(c))
			}

			return db.Where(compositeQuery)
		}
	}

	return LeaftFilter(condition.(common.SQLLeafCondition), common.And)
}

func OrFilterComposition(condition common.SQLCondition) func(db *gorm.DB) *gorm.DB {
	if condition.IsComposite() {
		composite := condition.(common.SQLCompositeCondition)
		return func(db *gorm.DB) *gorm.DB {
			compositeQuery := db
			for _, c := range composite.Conditions {
				compositeQuery = compositeQuery.Scopes(OrFilterComposition(c))
			}

			return db.Or(compositeQuery)
		}
	}
	return LeaftFilter(condition.(common.SQLLeafCondition), common.Or)
}

func NotFilterComposition(condition common.SQLCondition) func(db *gorm.DB) *gorm.DB {
	if condition.IsComposite() {
		composite := condition.(common.SQLCompositeCondition)
		return func(db *gorm.DB) *gorm.DB {
			compositeQuery := db
			for _, c := range composite.Conditions {
				compositeQuery = compositeQuery.Scopes(NotFilterComposition(c))
			}

			return db.Not(compositeQuery)
		}
	}
	return LeaftFilter(condition.(common.SQLLeafCondition), common.Not)
}

func LeaftFilter(condition common.SQLLeafCondition, compositor common.SQLCompositor) func(db *gorm.DB) *gorm.DB {
	fieldParts := strings.Split(condition.Field, ".")

	if len(fieldParts) > 2 {
		log.Fatalf("Invalid field: %s", condition.Field)
		log.Fatal("Field can only be nested at most one level deep")
		log.Fatal("For example field or relation.field are valid")
		return nil
	}

	if len(fieldParts) == 1 {
		switch compositor {
		case common.And:
			return AndLeaftFilter(condition)
		case common.Or:
			return OrLeaftFilter(condition)
		case common.Not:
			return NotLeaftFilter(condition)
		default:
			return AndLeaftFilter(condition)
		}
	}

	relation := fieldParts[0]
	field := fieldParts[1]

	p := pluralize.NewClient()
	tableName := p.Plural(relation)
	foreignKey := p.Singular(relation) + "_id"

	ids := []string{}
	sq := database.DB.Select("id").Table(tableName).
		Scopes(LeaftFilter(
			common.SQLLeafCondition{
				Field:      field,
				Value:      condition.Value,
				Comparator: condition.Comparator,
			},
			compositor)).
		Find(&ids)

	switch compositor {
	case common.And:
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(foreignKey+" IN (?)", sq)
		}
	case common.Or:
		return func(db *gorm.DB) *gorm.DB {
			return db.Or(foreignKey+" IN (?)", sq)
		}
	case common.Not:
		return func(db *gorm.DB) *gorm.DB {
			return db.Not(foreignKey+" IN (?)", sq)
		}
	default:
		return nil
	}
}

func AndLeaftFilter(condition common.SQLLeafCondition) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch condition.Comparator {
		case common.Equal:
			return db.Where(condition.Field+" = ?", condition.Value)
		case common.Like:
			return db.Where(condition.Field+" like ?", condition.Value)
		case common.ILike:
			return db.Where("LOWER("+condition.Field+") like LOWER(?)", condition.Value)
		case common.GreaterThan:
			return db.Where(condition.Field+" > ?", condition.Value)
		case common.LessEqualThan:
			return db.Where(condition.Field+" <= ?", condition.Value)
		case common.GreaterEqualThan:
			return db.Where(condition.Field+" >= ?", condition.Value)
		case common.LessThan:
			return db.Where(condition.Field+" < ?", condition.Value)
		case common.In:
			return db.Where(condition.Field+" in (?)", strings.Split(condition.Value, "|"))
		case common.IsNull:
			return db.Where(condition.Field + " IS NULL")
		case common.IsNotNull:
			return db.Where(condition.Field + " IS NOT NULL")
		default:
			log.Panic("Invalid comparator")
			return db.Where(condition.Field+" = ?", condition.Value)

		}
	}
}

func OrLeaftFilter(condition common.SQLLeafCondition) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch condition.Comparator {
		case common.Equal:
			return db.Or(condition.Field+" = ?", condition.Value)
		case common.Like:
			return db.Or(condition.Field+" like ?", condition.Value)
		case common.ILike:
			return db.Or("LOWER("+condition.Field+") like LOWER(?)", condition.Value)
		case common.GreaterThan:
			return db.Or(condition.Field+" > ?", condition.Value)
		case common.LessEqualThan:
			return db.Or(condition.Field+" <= ?", condition.Value)
		case common.GreaterEqualThan:
			return db.Or(condition.Field+" >= ?", condition.Value)
		case common.LessThan:
			return db.Or(condition.Field+" < ?", condition.Value)
		case common.In:
			return db.Or(condition.Field+" in (?)", strings.Split(condition.Value, "|"))
		case common.IsNull:
			return db.Or(condition.Field + " IS NULL")
		case common.IsNotNull:
			return db.Or(condition.Field + " IS NOT NULL")
		default:
			log.Panic("Invalid comparator")
			return db.Or(condition.Field+" = ?", condition.Value)
		}
	}
}

func NotLeaftFilter(condition common.SQLLeafCondition) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch condition.Comparator {
		case common.Equal:
			return db.Not(condition.Field+" = ?", condition.Value)
		case common.Like:
			return db.Not(condition.Field+" like ?", condition.Value)
		case common.ILike:
			return db.Not("LOWER("+condition.Field+") like LOWER(?)", condition.Value)
		case common.GreaterThan:
			return db.Not(condition.Field+" > ?", condition.Value)
		case common.LessEqualThan:
			return db.Not(condition.Field+" <= ?", condition.Value)
		case common.GreaterEqualThan:
			return db.Not(condition.Field+" >= ?", condition.Value)
		case common.LessThan:
			return db.Not(condition.Field+" < ?", condition.Value)
		case common.In:
			return db.Not(condition.Field+" in (?)", strings.Split(condition.Value, "|"))
		case common.IsNull:
			return db.Not(condition.Field + " IS NULL")
		case common.IsNotNull:
			return db.Not(condition.Field + " IS NOT NULL")
		default:
			log.Panic("Invalid comparator")
			return db.Not(condition.Field+" = ?", condition.Value)
		}
	}
}

func Order(orderBys common.OrderBys) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, orderBy := range orderBys {
			if orderBy.Direction == common.Asc {
				db = db.Order(orderBy.Field + " asc")
			} else {
				db = db.Order(orderBy.Field + " desc")
			}
		}
		return db
	}
}

func Preload(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}
		return db
	}
}
