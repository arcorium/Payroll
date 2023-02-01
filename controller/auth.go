package controller

import (
	"Penggajian/model"
	"Penggajian/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// JWT
	claims := jwt.MapClaims{}
	claims["id"] = user.Id.Hex()

	refreshToken, err := util.GenerateRefreshToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_RTOKEN)
	}
	// claims
	claims = jwt.MapClaims{
		"username":   user.Username,
		"id":         user.Id.Hex(),
		"authorized": true,
		"admin":      user.Type == model.Admin,
	}

	accessToken, err := util.GenerateAccessToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_ATOKEN)
	}

	// Load token to database
	//token := model.Token{Token: refreshToken, UserId: user.Id}
	// either update or inserting data
	//token, err = a.tokenRepo.UpsertTokenByUserId(user.Id, &token)

	// Set into logged in
	err = a.userRepo.UpdateLoggedIn(user.Id, true)

	// Set refresh token in cookie and response access token
	refreshCookie := GenerateTokenCookie(refreshToken)
	SetCookies(c, refreshCookie)

	response := model.ResponseToken{UserId: user.Id.Hex(), AccessToken: accessToken, RefreshToken: refreshToken}

	return SendSuccessResponse(c, fasthttp.StatusOK, response)
}

func (a *API) Logout(c *fiber.Ctx) error {
	SetDefaultContext(c)

	token := model.ResponseToken{}
	if util.IsError(c.BodyParser(&token)) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	cookie := c.Cookies(util.JWT_COOKIE_REFRESH_NAME, "")
	if util.IsEmpty(cookie) {
		if util.IsEmpty(token.RefreshToken) {
			return SendErrorResponse(c, fasthttp.StatusBadRequest, "rtoken doesn't found")
		}
		cookie = token.RefreshToken
	}

	// Remove token from database
	//token, err := a.tokenRepo.RemoveTokenByToken(refreshToken)
	//if err != nil {
	//	return SendErrorResponse(c, fasthttp.StatusNotFound, "token not found")
	//}

	// validate and get claims from refresh token
	refreshToken, err := jwt.Parse(cookie, a.jwtValidateToken)
	if util.IsError(err) {
		DeleteCookies(c, util.JWT_COOKIE_REFRESH_NAME)
		return SendErrorResponse(c, fasthttp.StatusBadRequest, err.Error())
	}
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)

	// Get id from access token
	user := c.Locals("user")
	accessToken := user.(*jwt.Token)
	accessClaims := accessToken.Claims.(jwt.MapClaims)

	// Check equivalent id
	if refreshClaims["id"].(string) != accessClaims["id"].(string) {
		DeleteCookies(c, util.JWT_COOKIE_REFRESH_NAME)
		return SendErrorResponse(c, fasthttp.StatusConflict, "token doesn't match")
	}

	userId, err := primitive.ObjectIDFromHex(accessClaims["id"].(string))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed id")
	}
	err = a.userRepo.UpdateLoggedIn(userId, false)

	// Clear cookies
	DeleteCookies(c, util.JWT_COOKIE_REFRESH_NAME)

	return SendSuccessResponse(c, fasthttp.StatusOK, model.NewResponseID(userId))
}

func (a *API) RequestToken(c *fiber.Ctx) error {
	SetDefaultContext(c)

	tokenReq := model.ResponseToken{}
	if util.IsError(c.BodyParser(&tokenReq)) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, RESPONSE_MSG_BODY)
	}

	// old refresh token
	cookie := c.Cookies(util.JWT_COOKIE_REFRESH_NAME, "")
	if util.IsEmpty(cookie) {
		if util.IsEmpty(tokenReq.RefreshToken) {
			return SendErrorResponse(c, fasthttp.StatusBadRequest, "rtoken doesn't found")
		}
		cookie = tokenReq.RefreshToken
	}

	// Get token from database
	//token, err := a.tokenRepo.GetTokenByToken(cookie)
	//if err != nil {
	//	return SendErrorResponse(c, fasthttp.StatusNotFound, "token not found")
	//}

	// validate and get claims from refresh token
	token, err := jwt.Parse(cookie, a.jwtValidateToken)
	if util.IsError(err) {
		DeleteCookies(c, util.JWT_COOKIE_REFRESH_NAME)
		return SendErrorResponse(c, fasthttp.StatusBadRequest, err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)

	// Get user
	userId, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if util.IsError(err) {
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "malformed id")
	}

	user, err := a.userRepo.GetUserById(userId)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// Generate new jwt
	claims = jwt.MapClaims{}
	claims["id"] = user.Id.Hex()

	refreshToken, err := util.GenerateRefreshToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_RTOKEN)
	}

	// claims
	claims = jwt.MapClaims{
		"username":   user.Username,
		"id":         user.Id,
		"authorized": user.IsLoggedIn,
		"admin":      user.Type == model.Admin,
	}
	accessToken, err := util.GenerateAccessToken(claims, []byte(a.config.SecretKey))
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, RESPONSE_MSG_ATOKEN)
	}

	// Invalidate last token by updating the token
	//err = a.tokenRepo.UpdateToken(token.Token, refreshToken)
	//if err != nil {
	//	return SendErrorResponse(c, fasthttp.StatusInternalServerError, "failed to update token")
	//}

	// Set refresh token in cookie
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
