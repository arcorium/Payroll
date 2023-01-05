package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"github.com/gofiber/fiber/v2"
	"time"
)

func SendSuccessResponse(c *fiber.Ctx, httpStatus_ int, data_ any) error {
	if err := c.SendStatus(httpStatus_); err != nil {
		return err
	}

	return c.JSON(model.NewSuccessResponse(CONDITION_SUCCESS, data_))
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

func GenerateTokenCookie(refreshToken_ string, accessToken_ string) (*fiber.Cookie, *fiber.Cookie) {
	refreshCookie := GenerateCookie(util.JWT_COOKIE_REFRESH_NAME, refreshToken_, time.Now().Add(util.JWT_COOKIE_REFRESH_TIMEOUT))
	accessCookie := GenerateCookie(util.JWT_COOKIE_ACCESS_NAME, accessToken_, time.Now().Add(util.JWT_COOKIE_ACCESS_TIMEOUT))

	return refreshCookie, accessCookie
}

func GenerateCookie(key_ string, value_ string, expiration_ time.Time) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Expires = expiration_
	cookie.Name = key_
	cookie.Value = value_
	cookie.HTTPOnly = true

	return cookie
}

func DeleteCookies(c *fiber.Ctx, keys_ ...string) {
	for _, key := range keys_ {
		SetCookies(c, GenerateCookie(key, "", time.Now().Add(time.Second*-3)))
	}
}

func SetCookies(c *fiber.Ctx, cookies_ ...*fiber.Cookie) {
	for _, cookie := range cookies_ {
		c.Cookie(cookie)
	}
}

const (
	CONDITION_ERROR   = "error"
	CONDITION_SUCCESS = "success"
)
