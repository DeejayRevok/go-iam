package refreshToken

const RefreshTokenDefaultExpirationHours = 8

type RefreshToken struct {
	Id  string
	Sub string
	Exp int64
}
