package controller

import (
	"Penggajian/model"
	"Penggajian/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
)

func (a *API) Login(c *fiber.Ctx) error {
	SetDefaultContext(c)

	user := model.User{}
	var err error
	if err = c.BodyParser(&user); err != nil {
		return err
	}

	if util.IsEmpty(user.Username) || util.IsEmpty(user.Password) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// Validate user with the username and password
	user, err = a.userRepo.ValidateUser(user.Username, user.Password)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// JWT
	claims := jwt.MapClaims{}

	refreshToken, err := util.GenerateRefreshToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_RTOKEN)
	}
	// claims
	claims = jwt.MapClaims{
		//"id":         id,
		"id":         user.Id,
		"authorized": true,
		"admin":      user.Type == model.Admin,
	}

	accessToken, err := util.GenerateAccessToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_ATOKEN)
	}

	// Load token to database
	token := model.Token{Token: refreshToken, UserId: user.Id}
	// either update or inserting data
	token, err = a.tokenRepo.UpsertTokenByUserId(user.Id, &token)

	// Set into logged in
	err = a.userRepo.UpdateLoggedIn(user.Id, true)

	// Set refresh token in cookie and response access token
	refreshCookie := GenerateTokenCookie(refreshToken)
	SetCookies(c, refreshCookie)

	response := model.ResponseToken{UserId: user.Id.Hex(), AccessToken: accessToken}

	return SendSuccessResponse(c, fasthttp.StatusOK, response)
}

func (a *API) Logout(c *fiber.Ctx) error {
	SetDefaultContext(c)

	refreshToken := c.Cookies(util.JWT_COOKIE_REFRESH_NAME, "")
	if util.IsEmpty(refreshToken) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "rtoken doesn't found")
	}

	// Remove token
	token, err := a.tokenRepo.RemoveTokenByToken(refreshToken)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "token not found")
	}

	err = a.userRepo.UpdateLoggedIn(token.UserId, false)

	// Clear cookies
	DeleteCookies(c, util.JWT_COOKIE_REFRESH_NAME)

	return SendSuccessResponse(c, fasthttp.StatusOK, model.NewResponseID(token.UserId))
}

func (a *API) RequestToken(c *fiber.Ctx) error {
	SetDefaultContext(c)

	// old refresh token
	cookie := c.Cookies(util.JWT_COOKIE_REFRESH_NAME, "")
	if util.IsEmpty(cookie) {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "data is not satisfied")
	}

	// Get token
	token, err := a.tokenRepo.GetTokenByToken(cookie)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "token not found")
	}

	// Get user
	user, err := a.userRepo.GetUserById(token.UserId)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// Generate new jwt
	claims := jwt.MapClaims{}

	refreshToken, err := util.GenerateRefreshToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_RTOKEN)
	}

	// claims
	claims = jwt.MapClaims{
		//"id":         id,
		"id":         user.Id,
		"authorized": user.IsLoggedIn,
		"admin":      user.Type == model.Admin,
	}
	accessToken, err := util.GenerateAccessToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_ATOKEN)
	}

	// Invalidate last token by updating the token
	err = a.tokenRepo.UpdateToken(token.Token, refreshToken)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to update token")
	}

	// Set refresh token and access token in cookie
	refreshCookie := GenerateTokenCookie(refreshToken)
	SetCookies(c, refreshCookie)

	// Response
	response := model.ResponseToken{AccessToken: accessToken}
	return SendSuccessResponse(c, fasthttp.StatusOK, response)
}

// validateAuthorization middleware to check claims in jwt on key "authorized" and
// set context locals type for superOnly middleware
func (a *API) validateAuthorization() fiber.Handler {

	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		if !claims["authorized"].(bool) {
			return SendErrorResponse(c, fasthttp.StatusUnauthorized, "account not eligible")
		}

		// Set local in context
		c.Locals("type", claims["admin"])

		return c.Next()
	}
}

// superOnly middleware for only allowing super admin user type to access,
// data's got from context locals named type
func (a *API) superOnly() fiber.Handler {

	return func(c *fiber.Ctx) error {
		if c.Locals("type").(bool) {
			return SendErrorResponse(c, fasthttp.StatusForbidden, "not allowed")
		}

		return c.Next()
	}
}
