package controller

import (
	"Penggajian/model"
	"Penggajian/util"
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

func GenerateTokenCookie(refreshToken_ string) *fiber.Cookie {
	refreshCookie := GenerateCookie(util.JWT_COOKIE_REFRESH_NAME, refreshToken_, time.Now().Add(util.JWT_COOKIE_REFRESH_TIMEOUT), "/api/v1/")

	return refreshCookie
}

func GenerateCookie(key_ string, value_ string, expiration_ time.Time, path_ ...string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Expires = expiration_
	cookie.Name = key_
	cookie.Value = value_
	cookie.HTTPOnly = true

	if path_ == nil {
		cookie.Path = "/"
	} else {
		cookie.Path = path_[0]
	}
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

const (
	RESPONSE_MSG_PARAMETER = "parameter is not satisfied"
	RESPONSE_MSG_BODY      = "body is not satisfied"
	RESPONSE_MSG_MALFORM   = "malformed id"
	RESPONSE_MSG_RTOKEN    = "failed to generate refresh token"
	RESPONSE_MSG_ATOKEN    = "failed to generate access token"
)
