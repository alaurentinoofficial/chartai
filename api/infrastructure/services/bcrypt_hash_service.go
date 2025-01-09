package infrastructure_services

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHashService struct{}

func NewBcryptHashService() *BcryptHashService {
	return &BcryptHashService{}
}

func (h *BcryptHashService) Hash(value string) string {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(hashedValue)
}

func (h *BcryptHashService) Validate(hashedValue string, value string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	if err != nil {
		return false
	}

	return true
}
