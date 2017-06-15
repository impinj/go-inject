package inject_test

type StructA struct {
	B *StructB `inject:""`
}

type StructB struct {
	A *StructA `inject:""`
}

type InterfaceA interface {
	MethodA()
}

type InterfaceB interface {
	MethodB()
}

type ImplA struct {
	B InterfaceB `inject:""`
}

func (a ImplA) MethodA() {}

type PtrImplA struct {
	B InterfaceB `inject:""`
}

func (a *PtrImplA) MethodA() {}

type ImplB struct {
	A InterfaceA `inject:""`
}

func (b ImplB) MethodB() {}

type PtrImplB struct {
	A InterfaceA `inject:""`
}

func (a *PtrImplB) MethodB() {}

type Decorator struct {
	Decorated InterfaceA `inject:""`
}

func (d Decorator) MethodA() {
	d.Decorated.MethodA()
}

type Decorated struct{}

func (d Decorated) MethodA() {}

type PtrDecorated struct{}

func (d *PtrDecorated) MethodA() {}

type ServiceInterface interface {
	ServiceMethod()
}

type ServiceValueImpl struct {
	X InterfaceA `inject:""`
}

func (s ServiceValueImpl) ServiceMethod() {}

type ServicePtrImpl struct {
	X InterfaceA `inject:""`
}

func (s *ServicePtrImpl) ServiceMethod() {}

type CustomA struct {
	Val int
}

func (a CustomA) MethodA() {}
