package controller

import (
	"Penggajian/pkg/model"
	"github.com/gofiber/fiber/v2"
)

func SendSuccessResponse(c *fiber.Ctx, httpStatus_ int, message_ string, data_ any) error {
	if err := c.SendStatus(httpStatus_); err != nil {
		return err
	}

	return c.JSON(model.NewSuccessResponse(CONDITION_SUCCESS, message_, data_))
}

func SendErrorResponse(c *fiber.Ctx, httpStatus_ int, message_ string) error {
	if err := c.SendStatus(httpStatus_); err != nil {
		return err
	}

	return c.JSON(model.NewErrorResponse(CONDITION_ERROR, message_))
}

func SetDefaultContext(c *fiber.Ctx, accepts_ ...string) {
	// Accept request format
	c.Accepts("json", "application/json")
	// Custom accept
	if len(accepts_) > 0 {
		c.Accepts(accepts_...)
	}

	// Language
	c.AcceptsLanguages("id", "en")

	// Encoding
	c.AcceptsEncodings("compress", "br")
}

const (
	CONDITION_ERROR   = "error"
	CONDITION_SUCCESS = "success"
)
