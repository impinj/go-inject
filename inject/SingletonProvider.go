package inject

type SingletonProvider struct {
	Provider

	value interface{}
}

func (p *SingletonProvider) Resolve() interface{} {
	if p.Provider == nil {
		return nil
	}

	if p.value == nil {
		p.value = p.Provider.Resolve()
	}

	return p.value
}
