package core_services

type HashService interface {
	Hash(value string) string
	Validate(hashedValue string, value string) bool
}
