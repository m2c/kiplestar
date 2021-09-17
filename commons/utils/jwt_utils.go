package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 产生json web token
func CreateJWT(data map[string]interface{}, expire time.Duration, secret string) (string, error) {
	t := make(jwt.MapClaims)
	for s := range data {
		t[s] = data[s]
	}
	t["exp"] = time.Now().Add(expire).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string, secret string) (map[string]interface{}, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claim.Claims.(jwt.MapClaims), nil
}

func GetJWTStringValue(token, secret, key string) (string, error) {
	mapData, err := ParseToken(token, secret)
	if err != nil {
		return "", err
	}
	return mapData[key].(string), nil
}
