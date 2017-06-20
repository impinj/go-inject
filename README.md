#go-inject

`import "github.com/impinj/go-inject"`

Go-inject is a dependency injection framework for Golang. It uses Go's reflection package to discover annotated fields and populate them.

## Set Me Up!

You'll need Golang v1.5+


```bash
go get github.com/impinj/go-inject
```

**Example**

To use go-inject, just annotate your classes which require injection. Because go-inject uses Go's reflection package, only public fields can be completed. If you provide a name for your injection, only objects with matching names will be supplied.

```go
package token

type Token interface {
        // Dunno, token-y things.
}

type TokenService interface {
        GetTokenByUsername(username string) Token
}
```

```go
package hello

import (
        "github.com/impinj/go-inject/sample/token"
)

type HelloService struct {
    TokenProvider token.TokenService `inject:""`
    Url string `inject:"helloservice.url"`
}
```

Elsewhere in your project, provide your object graph with objects for completion. After providing all required values, resolve the object graph to complete any partial objects.

```go
package main

import (
        "github.com/impinj/go-inject/inject"
        "github.com/impinj/go-inject/sample/hello"
        "github.com/impinj/go-inject/sample/oauth"
        "github.com/impinj/go-inject/sample/token"
)

func main() {
    var g inject.Graph
    
    g.Provide(
            &inject.ValueProvider{
                    Name: "helloservice.url",
                    Value: "https://www.impinj.com/hello",
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
    }
    
    // Run actual business logic.
}
```