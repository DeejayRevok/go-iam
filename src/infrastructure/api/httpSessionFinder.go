package api

import (
	"go-uaa/src/domain/session"
	"go-uaa/src/infrastructure/jwt"
	"net/http"
)

type HTTPSessionFinder struct {
	sessionDeserializer *jwt.JWTSessionDeserializer
}

func (finder *HTTPSessionFinder) Find(request *http.Request) (*session.Session, error) {
	cookies := request.Cookies()
	var session string
	for _, cookie := range cookies {
		if cookie.Name == "session" {
			session = cookie.Value
			break
		}
	}
	if session == "" {
		return nil, nil
	}

	return finder.sessionDeserializer.Deserialize(session)
}

func NewHTTPSessionFinder(sessionDeserializer *jwt.JWTSessionDeserializer) *HTTPSessionFinder {
	return &HTTPSessionFinder{
		sessionDeserializer: sessionDeserializer,
	}
}
