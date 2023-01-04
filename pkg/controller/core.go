package controller

import (
	"Penggajian/pkg/model"
	"Penggajian/pkg/util"
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
		return SendErrorResponse(c, fasthttp.StatusBadRequest, "body is not satisfied")
	}

	// Validate user with the username and password
	user, err = a.userRepo.ValidateUser(user.Username, user.Password)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// JWT
	//id := uuid.NewString() // unique id each jwt
	claims := jwt.MapClaims{}

	refreshToken, err := util.GenerateRefreshToken(claims)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't create rtoken")
	}
	// claims
	claims = jwt.MapClaims{
		//"id":         id,
		"name":       user.Username,
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

	// Prevent multi-login user
	if err = a.userRepo.IsLoggedIn(user.Id); err == nil {
		return SendErrorResponse(c, fasthttp.StatusConflict, "user already logged in")
	}

	// Set into logged in
	err = a.userRepo.UpdateLoggedIn(user.Id, true)

	response := model.ResponseToken{RefreshToken: refreshToken, AccessToken: accessToken}
	return SendSuccessResponse(c, fasthttp.StatusOK, "user "+user.Username+" logged in!", response)
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
		return SendErrorResponse(c, fasthttp.StatusNotFound, "token can't be deleted")
		//return c.JSON(bson.M{"data": "kontol"})
	}

	// Get user from database
	user, err := a.userRepo.GetUserById(token.UserId)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusUnauthorized, "user not found")
	}

	// Don't logout user that not logged in
	if err = a.userRepo.IsLoggedIn(user.Id); err != nil {
		return SendErrorResponse(c, fasthttp.StatusConflict, "user is not logged in yet")
	}

	err = a.userRepo.UpdateLoggedIn(user.Id, false)

	return SendSuccessResponse(c, fasthttp.StatusOK, "user "+user.Username+" logged out!", nil)
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
		"name":       user.Username,
		"authorized": true,
		"admin":      user.Type == model.Admin,
	}
	accessToken, err := util.GenerateAccessToken(claims)

	// Invalidate last token by updating the token
	err = a.tokenRepo.UpdateToken(token.Token, refreshToken)
	if err != nil {
		return SendErrorResponse(c, fasthttp.StatusInternalServerError, "can't update token")
	}

	// Response
	response := model.ResponseToken{RefreshToken: refreshToken, AccessToken: accessToken}
	return SendSuccessResponse(c, fasthttp.StatusOK, "created pair of token", response)
}
