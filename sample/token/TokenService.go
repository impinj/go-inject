package token

type Token interface {
	// Dunno, token-y things.
}

type TokenService interface {
	GetTokenByUsername(username string) Token
}
