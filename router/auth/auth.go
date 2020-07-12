package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	e "github.com/donohutcheon/gowebserver/controllers/errors"

)

const AccessTokenLifeSpan = 36000
const RefreshTokenLifeSpan = 864000
const APITokenLifeSpan = 31536000

type JSONWebToken struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

type RefreshJWTReq struct {
	GrantType    string `json:"grantType" sql:"-"`
	RefreshToken string `json:"refreshToken" sql:"-"`
}

type TokenResponse struct {
	ExpiresIn int64 `json:"expiresIn"`
	AccessToken string  `json:"accessToken" sql:"-"`
	RefreshToken string  `json:"refreshToken" sql:"-"`
}

type APITokenResponse struct {
	ExpiresIn int64 `json:"expiresIn"`
	APIToken string `json:"apiToken" sql:"-"`
}

func CreateToken(userID int64) (*TokenResponse, error){
	token := new(TokenResponse)
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + AccessTokenLifeSpan
	token.ExpiresIn = expireDateTime
	accessToken := &JSONWebToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}

	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	accessTokenString, _ := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	token.AccessToken = accessTokenString

	refreshToken := &JSONWebToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: epochSecs + RefreshTokenLifeSpan,
			IssuedAt:  epochSecs,
		},
	}
	signedRefreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	refreshTokenString, _ := signedRefreshToken.SignedString([]byte(os.Getenv("token_password")))
	token.RefreshToken = refreshTokenString

	return token, nil
}

func RefreshToken(rawToken string) (*TokenResponse, error) {
	tk := new(JSONWebToken)

	token, err := jwt.ParseWithClaims(rawToken, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil { //Malformed token, returns with http code 403 as usual
		return nil, e.Wrap("Token rejected", http.StatusForbidden, err)
	}

	if !token.Valid { //Token is invalid, maybe not signed on this server
		return nil, e.NewError("token is not valid", nil, http.StatusForbidden)
	}

	fmt.Printf("UserID %d", tk.UserID)

	//Create JWT token
	tokenResp, err := CreateToken(tk.UserID)
	if err != nil {
		return nil, e.Wrap("token creation failed", http.StatusInternalServerError, err)
	}

	return tokenResp, nil
}

func CreateAPIToken(userID int64) (*APITokenResponse, error){
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + APITokenLifeSpan

	accessToken := &JSONWebToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}

	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	apiTokenString, err := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return nil, err
	}

	token := &APITokenResponse{
		ExpiresIn: expireDateTime,
		APIToken: apiTokenString,
	}

	return token, nil
}
