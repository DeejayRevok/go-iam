package session

const SessionExpirationHours = 4

type Session struct {
	UserID    string
	LastUsage int64
	Exp       int64
}
