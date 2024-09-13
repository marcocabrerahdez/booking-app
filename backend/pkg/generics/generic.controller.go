package generics

import (
	"backend/pkg/common"
	"backend/pkg/helpers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GenericController interface {
	GetResourceNames() ResourceNames

	Get() fiber.Handler
	GetAll() fiber.Handler
	Update() fiber.Handler
	Create() fiber.Handler
	Delete() fiber.Handler
	Count() fiber.Handler
	HardDelete() fiber.Handler
	GetDeleted() fiber.Handler
	GetAllDeleted() fiber.Handler
}

type GenericControllerImpl[E common.Entity, DTO common.DTO] struct {
	names      ResourceNames
	repository GenericRepository[E, DTO]
}

type ResourceNames struct {
	Singular string
	Plural   string
}

func NewController[E common.Entity, DTO common.DTO](names ResourceNames) GenericControllerImpl[E, DTO] {
	var controller GenericControllerImpl[E, DTO]
	controller.repository = NewGenericRepositoryGORM[E, DTO]()
	controller.names = names
	return controller
}

func (imp GenericControllerImpl[E, DTO]) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s id", imp.names.Singular))
		}

		relations := common.RelationsFromQuery(c)

		entity, err := imp.repository.FindOne(id, relations)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("%s not found", imp.names.Singular))
		}

		return Found(c, entity.ToDTO(), fmt.Sprintf("Found %s", imp.names.Singular))
	}
}

func (imp GenericControllerImpl[E, DTO]) GetAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageable, err := common.PageableFromQuery(c)
		if err != nil {
			return BadRequest(c, err, "Invalid pagination parameters")
		}

		relations := common.RelationsFromQuery(c)
		conditions := common.ConditionsFromQuery(c)
		orders := common.OrderBysFromQuery(c)

		result, err := imp.repository.FindAll(pageable, conditions, relations, orders)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("%s not found", imp.names.Plural))
		}

		dtos := make([]DTO, len(result.Items))
		for i, entity := range result.Items {
			dtos[i] = entity.ToDTO().(DTO)
		}

		return Found(c,
			common.NewPage[DTO](
				dtos,
				result.Page,
				result.Size,
				result.Total,
				result.Filtered),
			fmt.Sprintf("Found %s", imp.names.Plural))

	}
}

func (imp GenericControllerImpl[E, DTO]) Update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validate id
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s id", imp.names.Singular))
		}
		// Validate payload
		var payload DTO
		if err := c.BodyParser(&payload); err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload", imp.names.Singular))
		}

		// Validate entity existence
		exists, err := imp.repository.Exists(id)
		if err != nil || !exists {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload", imp.names.Singular))
		}

		// Update entity
		entity := payload.ToEntity().(E)

		entity.SetID(id)
		entity, err = imp.repository.Update(entity)
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload", imp.names.Singular))
		}

		return Updated(c, entity.ToDTO(), fmt.Sprintf("%s updated", imp.names.Singular))
	}
}

func (imp GenericControllerImpl[E, DTO]) Create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dto DTO
		if err := c.BodyParser(&dto); err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload", imp.names.Singular))
		}

		entity := dto.ToEntity().(E)
		entity, err := imp.repository.Create(entity)
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload", imp.names.Singular))
		}

		return Created(c, entity.ToDTO(), fmt.Sprintf("%s created", imp.names.Singular))
	}
}

func (imp GenericControllerImpl[E, DTO]) Delete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s id", imp.names.Singular))
		}

		relations := common.RelationsFromQuery(c)

		entity, err := imp.repository.FindOne(id, relations)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("%s not found", imp.names.Singular))
		}

		err = imp.repository.Delete(entity)
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload: %s", imp.names.Singular, helpers.PrettyStruct(entity)))
		}

		return Deleted(c, fmt.Sprintf("%s deleted", imp.names.Singular))
	}
}

func (imp GenericControllerImpl[E, DTO]) GetResourceNames() ResourceNames {
	return imp.names
}

func (imp GenericControllerImpl[E, DTO]) Count() fiber.Handler {
	return func(c *fiber.Ctx) error {
		conditions := common.ConditionsFromQuery(c)

		count, err := imp.repository.Count(conditions)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("Error counting %s", imp.names.Plural))
		}

		return Found(c, count, fmt.Sprintf("Counted %s", imp.names.Plural))
	}
}

func (imp GenericControllerImpl[E, DTO]) HardDelete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s id", imp.names.Singular))
		}

		entity, err := imp.repository.GetOneDeleted(id)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("%s not found", imp.names.Singular))
		}

		err = imp.repository.HardDelete(entity)
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s payload", imp.names.Singular))
		}

		return Deleted(c, fmt.Sprintf("%s deleted", imp.names.Singular))
	}
}

func (imp GenericControllerImpl[E, DTO]) GetAllDeleted() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageable, err := common.PageableFromQuery(c)
		if err != nil {
			return BadRequest(c, err, "Invalid pagination parameters")
		}

		filters := common.ConditionsFromQuery(c)
		relations := common.RelationsFromQuery(c)
		orderBys := common.OrderBysFromQuery(c)

		result, err := imp.repository.GetDeleted(pageable, filters, relations, orderBys)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("%s not found", imp.names.Plural))
		}

		dtos := make([]DTO, len(result.Items))
		for i, entity := range result.Items {
			dtos[i] = entity.ToDTO().(DTO)
		}

		return Found(c,
			common.NewPage[DTO](
				dtos,
				result.Page,
				result.Size,
				result.Total,
				result.Filtered),
			fmt.Sprintf("Found %s", imp.names.Plural))

	}
}

func (imp GenericControllerImpl[E, DTO]) GetDeleted() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return BadRequest(c, err, fmt.Sprintf("Invalid %s id", imp.names.Singular))
		}

		entity, err := imp.repository.GetOneDeleted(id)
		if err != nil {
			return NotFound(c, err, fmt.Sprintf("%s not found", imp.names.Singular))
		}

		return Found(c, entity.ToDTO(), fmt.Sprintf("Found %s", imp.names.Singular))
	}
}
