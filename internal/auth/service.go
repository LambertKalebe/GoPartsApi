package auth

import (
	"database/sql"
	"errors"
	"g0/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var errInvalidCredentials = errors.New("usuário ou senha inválidos")
var errUserAlreadyExists = errors.New("usuário já existe")

func serviceLogin(req loginRequest) (loginResponse, error) {
	user, err := checkUser(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return loginResponse{}, errInvalidCredentials
		}
		return loginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PassHash),
		[]byte(req.Password),
	)

	if err != nil {
		return loginResponse{}, errInvalidCredentials
	}

	token, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
	)

	if err != nil {
		return loginResponse{}, err
	}

	return toLoginResponse(token), nil
}

func serviceRegister(req registerRequest) (registerResponse, error) {
	_, err := checkUser(req.Username)

	if err == nil {
		return registerResponse{}, errUserAlreadyExists
	}

	EncryptedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	err = insertUser(req.Username, string(EncryptedPass), "user")
	if err != nil {
		return registerResponse{}, err
	}
	return registerResponse{}, nil
}

func serviceLogout() logoutResponse {
	return toLogoutResponse()
}
func serviceAboutMe(token *jwt.Token) (aboutMeResponse, error) {
	claims := token.Claims.(jwt.MapClaims)

	u := user{
		ID:       int(claims["user_id"].(float64)),
		Username: claims["username"].(string),
		Role:     claims["role"].(string),
	}

	return toAboutMeResponse(u), nil
}
