package auth

import "github.com/dgrijalva/jwt-go"

type JSONWebToken struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

type RefreshJWTReq struct {
	GrantType    string `json:"grantType" sql:"-"`
	RefreshToken string `json:"refreshToken" sql:"-"`
}