package infrastructure_services

import (
	"fmt"
	"log"

	infrastructure_config "github.com/alaurentinoofficial/chartai/infrastructure/configs"
	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenService struct {
	signature string
}

func NewJwtTokenService(config *infrastructure_config.Config) *JwtTokenService {
	return &JwtTokenService{signature: config.TokenSignature}
}

func (h *JwtTokenService) Generate(accountId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"accountId": accountId,
	})

	signedToken, err := token.SignedString([]byte(h.signature))
	if err != nil {
		log.Fatal(err.Error())
	}

	return signedToken
}

func (h *JwtTokenService) GetClaims(tokenString string) (map[string]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}

		return h.signature, nil
	})
	if token == nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		result := make(map[string]string)
		for k, v := range claims {
			result[k] = v.(string)
		}

		return result, nil
	}

	return map[string]string{}, nil
}

func (h *JwtTokenService) Validate(token string) bool {
	return true
}
