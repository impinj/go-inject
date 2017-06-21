package inject

import "reflect"

type BuilderProvider struct {
	Name           string
	Context        reflect.Type
	Builder        interface{}
	ResolveContext Graph
}

func (p BuilderProvider) GetContext() reflect.Type {
	return p.Context
}

func (p BuilderProvider) GetName() string {
	return p.Name
}

func (p BuilderProvider) GetType() reflect.Type {
	return reflect.TypeOf(p.Builder).Out(0)
}

func (p BuilderProvider) IsComplete() bool {
	return false
}

func (p BuilderProvider) Resolve() interface{} {
	typeInfo := reflect.TypeOf(p.Builder)
	if typeInfo == nil || typeInfo.Kind() != reflect.Func {
		return nil
	}

	args := make([]reflect.Value, typeInfo.NumIn())
	for i := 0; i < typeInfo.NumIn(); i++ {
		argTypeInfo := typeInfo.In(i)
		switch argTypeInfo.Kind() {
		case reflect.Interface:
			if found, err := p.ResolveContext.Find(argTypeInfo, nil, ""); err == nil {
				args[i] = reflect.ValueOf(found)
			} else {
				return nil
			}

		default:
			args[i] = reflect.New(argTypeInfo)
			if err := p.ResolveContext.Complete(args[i].Interface()); err != nil {
				return nil
			}

			args[i] = args[i].Elem()
		}
	}

	v := reflect.ValueOf(p.Builder).Call(args)
	return v[0].Interface()
}
