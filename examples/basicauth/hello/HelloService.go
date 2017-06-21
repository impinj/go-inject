package hello

import (
	"github.com/impinj/go-inject/examples/basicauth/token"
)

type HelloService struct {
	TokenProvider token.TokenService `inject:""`
	Url string `inject:"helloservice.url"`
}
