package accessToken

const AccessTokenDefaultExpirationHours = 1

type AccessToken struct {
	Iss   string
	Sub   string
	Exp   int64
	Iat   int64
	Scope string
}
