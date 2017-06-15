package inject

import "reflect"

type Provider interface {
	GetContext() reflect.Type
	GetName() string
	GetType() reflect.Type
	IsComplete() bool

	Resolve() interface{}
}