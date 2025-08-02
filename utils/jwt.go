package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var JwtKey = []byte(os.Getenv("JWT_SECRET"))



type Claims struct{
	Email string `json:"email"`
	Role string`json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(email,role string)(string,error){
	claims :=&Claims{
		Email:email,
		Role:role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
		},

	}
token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
return token.SignedString(JwtKey)
}

func ValidateJWT(tokenString string )(*Claims,error){
	token,err:=jwt.ParseWithClaims(tokenString,&Claims{},func(token*jwt.Token)(interface{},error){
		return JwtKey,nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}