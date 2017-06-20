package inject

type SingletonProvider struct {
	Provider
	Value interface{}
}

func (p SingletonProvider) IsComplete() bool {
	return p.Value != nil
}

func (p *SingletonProvider) Resolve() interface{} {
	if p.Provider == nil {
		return nil
	}

	if p.Value == nil {
		p.Value = p.Provider.Resolve()
	}

	return p.Value
}
