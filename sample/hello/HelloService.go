package hello

import (
	"github.com/impinj/go-inject/sample/token"
)

type HelloService struct {
	TokenProvider token.TokenService `inject:""`
	Url string `inject:"helloservice.url"`
}
