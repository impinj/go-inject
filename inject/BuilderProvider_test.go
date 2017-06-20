package inject_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/impinj/go-inject/inject"
	. "reflect"
)

var _ = Describe("BuilderProvider", func() {
	var (
		p Provider
	)

	type Struct struct {
		State int
	}

	Describe("Resolve", func() {
		Context("Given a builder function", func() {
			DescribeTable("The provider resolves to the correct type",
				func(builder interface{}) {
					p = BuilderProvider{
						Builder: builder,
					}

					Expect(TypeOf(p.Resolve())).To(Equal(TypeOf(builder).Out(0)))
				},
				Entry("Value", func() Struct {
					return Struct{}
				}),
				Entry("Pointer", func() *Struct {
					return &Struct{}
				}),
			)

			DescribeTable("The provider resolves to different objects",
				func(builder interface{}) {
					p = BuilderProvider{
						Builder: builder,
					}

					r1 := p.Resolve()
					r2 := p.Resolve()

					Expect(ValueOf(r1)).ToNot(Equal(ValueOf(r2)))
				},
				Entry("Value", func() Struct {
					return Struct{}
				}),
				Entry("Pointer", func() *Struct {
					return &Struct{}
				}),
			)
		})
	})
})
