package core_services

type TokenService interface {
	Generate(accountId string) string
	Validate(token string) bool
	GetClaims(tokenString string) (map[string]string, error)
}
