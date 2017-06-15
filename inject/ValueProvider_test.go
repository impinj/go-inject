package inject_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "go-inject/inject"
	. "reflect"
)

var _ = Describe("ValueProvider", func() {
	var (
		p Provider
	)

	type Struct struct {
		State int
	}

	Describe("Resolve", func() {
		Context("Given a concrete object", func() {
			DescribeTable("The provider resolves to the correct type",
				func(obj interface{}, typeInfo Type) {
					p = ValueProvider{
						Value: obj,
					}

					Expect(TypeOf(p.Resolve())).To(Equal(typeInfo))
				},
				Entry("Value", Struct{}, TypeOf(Struct{})),
				Entry("Pointer", &Struct{}, PtrTo(TypeOf(Struct{}))),
			)

			DescribeTable("The provider resolves to the same object",
				func(obj interface{}) {
					p = ValueProvider{
						Value: obj,
					}

					r1 := p.Resolve()
					r2 := p.Resolve()

					Expect(ValueOf(r1)).To(Equal(ValueOf(r2)))
				},
				Entry("Value", Struct{}),
				Entry("Pointer", &Struct{}),
			)
		})
	})
})
