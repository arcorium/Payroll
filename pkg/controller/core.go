package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"log"
)

func (a *API) Login(c *fiber.Ctx) error {
	SetDefaultContext(c)

	user := model.User{}
	var err error
	if err = c.BodyParser(&user); err != nil {
		return err
	}

	if util.IsEmpty(user.Username) || util.IsEmpty(user.Password) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	// Validate user with the username and password
	user, err = a.userRepo.ValidateUser(user.Username, user.Password)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// Prevent multi-login user
	if user.IsLoggedIn {
		return SendErrorResponse(c, fasthttp.StatusConflict, "user already logged in")
	}

	// JWT
	claims := jwt.MapClaims{}

	refreshToken, err := util.GenerateRefreshToken(claims)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't create rtoken")
	}
	// claims
	claims = jwt.MapClaims{
		//"id":         id,
		"id":         user.Id,
		"authorized": true,
		"admin":      user.Type == model.Admin,
	}

	accessToken, err := util.GenerateAccessToken(claims)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't create atoken")
	}

	// Load token to database
	token := model.Token{Token: refreshToken, UserId: user.Id}
	_, err = a.tokenRepo.AddToken(&token)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't update token into database")
	}

	// Set into logged in
	err = a.userRepo.UpdateLoggedIn(user.Id, true)

	// Set refresh token and access token in cookie
	refreshCookie, accessCookie := GenerateTokenCookie(refreshToken, accessToken)
	SetCookies(c, refreshCookie, accessCookie)

	response := model.ResponseToken{RefreshToken: refreshToken, AccessToken: accessToken}
	return SendSuccessResponse(c, fasthttp.StatusOK, response)
}

func (a *API) Logout(c *fiber.Ctx) error {
	SetDefaultContext(c)

	token := model.Token{}
	var err error
	if err = c.BodyParser(&token); err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	if util.IsEmpty(token.Token) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	// Remove JWT from database
	token, err = a.tokenRepo.RemoveTokenByToken(token.Token)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusNotFound, "token not found")
	}

	// Get user from database by the id
	user, err := a.userRepo.GetUserById(token.UserId)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// Don't logout user that not logged in
	if !user.IsLoggedIn {
		return SendErrorResponse(c, fasthttp.StatusConflict, "user is not logged in yet")
	}

	err = a.userRepo.UpdateLoggedIn(user.Id, false)

	// Clear cookies
	DeleteCookies(c, util.JWT_COOKIE_REFRESH_NAME, util.JWT_COOKIE_ACCESS_NAME)

	return SendSuccessResponse(c, fasthttp.StatusOK, model.ResponseID{Id: user.Id})
}

func (a *API) RequestToken(c *fiber.Ctx) error {
	SetDefaultContext(c)

	// old refresh token
	token := model.Token{}

	// Parse body
	if err := c.BodyParser(&token); err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	// Get token
	token, err := a.tokenRepo.ValidateToken(token.Token)
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

	refreshToken, err := util.GenerateRefreshToken(claims)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't create rtoken")
	}

	// claims
	claims = jwt.MapClaims{
		//"id":         id,
		"id":         user.Id,
		"authorized": user.IsLoggedIn,
		"admin":      user.Type == model.Admin,
	}
	accessToken, err := util.GenerateAccessToken(claims)

	// Invalidate last token by updating the token
	err = a.tokenRepo.UpdateToken(token.Token, refreshToken)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't update token")
	}

	// Set refresh token and access token in cookie
	refreshCookie, accessCookie := GenerateTokenCookie(refreshToken, accessToken)
	SetCookies(c, refreshCookie, accessCookie)

	// Response
	response := model.ResponseToken{RefreshToken: refreshToken, AccessToken: accessToken}
	return SendSuccessResponse(c, fasthttp.StatusOK, response)
}

// validateAuthorization middleware to check claims in jwt on key "authorized" and
// set context locals type for superOnly middleware
func (a *API) validateAuthorization() fiber.Handler {

	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		if !claims["authorized"].(bool) {
			log.Println("not authorized")
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
