package generics

import (
	"backend/pkg/common"
	"backend/pkg/helpers"

	"github.com/gofiber/fiber/v2"
)

func BadRequest(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusBadRequest).
		JSON(common.NewErrorResponse(err, message))
}

func Unauthorized(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusUnauthorized).
		JSON(common.NewErrorResponse(err, message))
}

func NotFound(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusNotFound).
		JSON(common.NewErrorResponse(err, message))
}

func PayloadValidationFailed(c *fiber.Ctx, errors []*helpers.ValidationErrors, message string) error {
	return c.Status(fiber.StatusBadRequest).
		JSON(common.NewValidationErrorResponse(errors, message))
}

func InternalServerError(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusInternalServerError).
		JSON(common.NewErrorResponse(err, message))
}

func Found(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).
		JSON(common.NewSuccessResponse(data, message))
}

func Updated(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).
		JSON(common.NewSuccessResponse(data, message))
}

func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).
		JSON(common.NewSuccessResponse(data, message))
}

func Deleted(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).
		JSON(common.NewSuccessResponse("n/a", message))
}

func Ok(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).
		JSON(common.NewSuccessResponse(data, message))
}

func Unimplemented(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotImplemented).
		JSON(common.NewErrorResponse(nil, message))
}
