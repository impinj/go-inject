package inject

import "reflect"

type ValueProvider struct {
	Name     string
	Complete bool
	Context  reflect.Type
	Value    interface{}
}

func (p ValueProvider) GetContext() reflect.Type {
	return p.Context
}

func (p ValueProvider) GetName() string {
	return p.Name
}

func (p ValueProvider) GetType() reflect.Type {
	return reflect.TypeOf(p.Value)
}

func (p ValueProvider) IsComplete() bool {
	return p.Complete
}

func (p ValueProvider) Resolve() interface{} {
	return p.Value
}
