package inject

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const injectTag string = "inject"

type Graph interface {
	Complete(v interface{}) error
	Find(typeInfo, context reflect.Type, name string) (interface{}, error)
	Provide(providers ...Provider)
	Resolve() error
}

func NewGraph() Graph {
	return &graph{}
}

type graph struct {
	providers []Provider
}

func (g *graph) Provide(providers ...Provider) {
	g.providers = append(g.providers, providers...)
}

func (g graph) Complete(v interface{}) error {
	switch k := reflect.TypeOf(v).Kind(); k {
	case reflect.Ptr:

	default:
		return fmt.Errorf("Tried to complete a non-pointer object (%s).", k)
	}

	q := []interface{}{v}
	for len(q) > 0 {
		c := q[0]
		q = q[1:]

		el := deferenceValue(reflect.ValueOf(c))
		providersByType := selectProvidersByType(g.providers, el.Type())
		for _, provider := range providersByType {
			if provider.IsComplete() {
				reflect.ValueOf(v).Elem().Set(reflect.ValueOf(provider.Resolve()))
				return nil
			}
		}

		if err := completeHelper(g.providers, el); err != nil {
			return err
		}

		for _, field := range selectInjectableFields(el) {
			if !isComplete(deferenceValue(field)) {
				q = append(q, field.Interface())
			}
		}
	}

	return nil
}

func (g graph) Find(typeInfo, context reflect.Type, name string) (interface{}, error) {
	if provider, err := findHelper(g.providers, typeInfo, context, name); err != nil {
		return nil, err
	} else {
		return provider.Resolve(), nil
	}
}

func (g graph) Resolve() error {
	var errs []error

	for _, provider := range g.providers {
		if valProv, ok := provider.(*ValueProvider); ok &&
			valProv.Value != nil &&
			!provider.IsComplete() {
			switch reflect.ValueOf(provider.Resolve()).Kind() {
			case reflect.Ptr:
				if err := g.Complete(valProv.Value); err != nil {
					err = fmt.Errorf("Encountered error while trying to complete (%s): %s",
						provider.GetType().String(),
						err)
					errs = append(errs, err)
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func completeHelper(providers []Provider, el reflect.Value) error {
	for i, field := range selectInjectableFields(el) {
		fieldInfo := el.Type().Field(i)
		name := fieldInfo.Tag.Get(injectTag)

		if !field.CanSet() {
			return fmt.Errorf("Cannot set a field (%s) on a non-struct object (%s).", fieldInfo.Name, el.Type())
		}

		if provider, err := findHelper(providers, field.Type(), el.Type(), name); err != nil {
			return fmt.Errorf("Encountered error attempting to set a field (%s) of type %s: %s",
				fieldInfo.Name,
				fieldInfo.Type,
				err)
		} else {
			field.Set(reflect.ValueOf(provider.Resolve()))
		}
	}

	return nil
}

func deferenceValue(el reflect.Value) reflect.Value {
	for {
		switch el.Kind() {
		case reflect.Interface:
			fallthrough

		case reflect.Ptr:
			el = el.Elem()

		default:
			return el
		}
	}
}

func findHelper(providers []Provider, typeInfo, context reflect.Type, name string) (Provider, error) {
	providersForType := selectProvidersByType(providers, typeInfo)
	providersForType = selectProvidersByContext(providersForType, context)
	providersForType = selectProvidersByName(providersForType, name)

	switch len(providersForType) {
	case 0:
		return nil, fmt.Errorf("Could not find provider for %s.", typeInfo)

	case 1:
		return providersForType[0], nil

	default:
		return nil, fmt.Errorf("Found multiple providers for type: %s, context: %s, name: %s.",
			typeInfo,
			context,
			name)
	}
}

func isComplete(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Struct:

	default:
		return true
	}

	for _, field := range selectInjectableFields(v) {
		switch field.Kind() {
		case reflect.Interface:
			fallthrough

		case reflect.Ptr:
			if field.IsNil() {
				return false
			}
		}
	}

	return true
}

func selectInjectableFields(v reflect.Value) map[int]reflect.Value {
	injectableFields := map[int]reflect.Value{}
	typeInfo := v.Type()
	for i := 0; i < typeInfo.NumField(); i++ {
		fieldInfo := typeInfo.Field(i)
		if _, ok := structTagLookup(fieldInfo.Tag, injectTag); ok {
			injectableFields[i] = v.Field(i)
		}
	}

	return injectableFields
}

func selectProvidersByContext(providers []Provider, context reflect.Type) []Provider {
	if context == nil {
		return providers
	}

	var contextMatches, noContexts []Provider
	for _, provider := range providers {
		ctx := provider.GetContext()
		if ctx == nil {
			noContexts = append(noContexts, provider)
		} else if context.ConvertibleTo(ctx) ||
			reflect.PtrTo(context).ConvertibleTo(ctx) {
			contextMatches = append(contextMatches, provider)
		}
	}

	if len(contextMatches) == 0 {
		return noContexts
	}

	return contextMatches
}

func selectProvidersByName(providers []Provider, name string) []Provider {
	var retVal []Provider
	for _, provider := range providers {
		if strings.EqualFold(provider.GetName(), name) {
			retVal = append(retVal, provider)
		}
	}

	return retVal
}

func selectProvidersByType(providers []Provider, typeInfo reflect.Type) []Provider {
	var providersForType []Provider
	for _, candidate := range providers {
		if candidate.GetType().AssignableTo(typeInfo) {
			providersForType = append(providersForType, candidate)
		}
	}

	return providersForType
}

func structTagLookup(structTag reflect.StructTag, key string) (string, bool) {
	tags := strings.Split((string)(structTag), " ")

	pattern := `(.+?):"(.*?)"`
	r := regexp.MustCompile(pattern)

	for _, tag := range tags {
		matches := r.FindStringSubmatch(tag)
		if len(matches) == 3 {
			return matches[2], matches[1] == key
		}
	}

	return "", false
}
