package main

import (
	"github.com/impinj/go-inject/inject"
	"github.com/impinj/go-inject/examples/basicauth/hello"
	"github.com/impinj/go-inject/examples/basicauth/oauth"
	"github.com/impinj/go-inject/examples/basicauth/token"
)

func main() {
	g := inject.NewGraph()

	g.Provide(
		&inject.ValueProvider{
			Name: "helloservice.url",
			Value: "https://www.example.org/hello",
		},
		&inject.SingletonProvider{
			Provider: &inject.BuilderProvider{
				Builder: func() token.TokenService {
					return &oauth.OAuthTokenService{}
				},
			},
		},
		&inject.ValueProvider{
			Value: &hello.HelloService{},
		},
	)

	if err := g.Resolve(); err != nil {
		// Handle error appropriately.
		panic(err)
	}

	// Run actual business logic.
}
