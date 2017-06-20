package oauth

import "github.com/impinj/go-inject/sample/token"

type OAuthTokenService struct {}

type oauthToken struct {
	username string
}

func (ts OAuthTokenService) GetTokenByUsername(username string) token.Token {
	return oauthToken{
		username: username,
	}
}