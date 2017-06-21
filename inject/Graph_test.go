package inject_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"fmt"
	. "github.com/impinj/go-inject/inject"
	"reflect"
)

var _ = Describe("Graph", func() {
	var (
		g Graph
	)

	BeforeEach(func() {
		g = NewGraph()
	})

	Describe("Complete", func() {
		Context("Attempt to complete a non-pointer object", func() {
			var v StructB

			It("Fails gloriously", func() {
				Expect(g.Complete(v)).ToNot(Succeed())
			})
		})

		Context("No found providers", func() {
			var v StructB

			It("Returns an error", func() {
				Expect(g.Complete(&v)).ToNot(Succeed())
			})
		})

		Context("Too many providers", func() {
			var (
				a, b StructA
			)

			BeforeEach(func() {
				g.Provide(
					&ValueProvider{
						Value: &a,
					},
					&ValueProvider{
						Value: &b,
					},
				)
			})

			It("Returns an error", func() {
				var v StructB
				Expect(g.Complete(&v)).ToNot(Succeed())
			})
		})

		Context("Complete a primitive", func() {
			var x, val int

			BeforeEach(func() {
				x = 0
				val = 5

				g.Provide(
					&ValueProvider{
						Complete: true,
						Value:    val,
					},
				)
			})

			Context("As a value", func() {
				It("Does not resolve correctly", func() {
					Expect(g.Complete(x)).ToNot(Succeed())
					Expect(x).To(Equal(0))
				})
			})

			Context("As a pointer", func() {
				It("Resolves correctly", func() {
					Expect(g.Complete(&x)).To(Succeed())
					Expect(x).To(Equal(val))
				})
			})
		})

		Context("Simple graph", func() {
			type Struct struct {
				X int `inject:""`
			}

			var (
				val int
				v   Struct
			)

			BeforeEach(func() {
				val = 5

				g.Provide(
					&ValueProvider{
						Value: val,
					},
				)
			})

			It("Resolves correctly", func() {
				Expect(v.X).To(Equal(0))
				Expect(g.Complete(&v)).To(Succeed())
				Expect(v.X).To(Equal(val))
			})
		})

		Context("Graph with a cycle", func() {
			var (
				a StructA
				b StructB
			)

			BeforeEach(func() {
				g.Provide(
					&ValueProvider{
						Value: &a,
					},
					&ValueProvider{
						Value: &b,
					},
				)
			})

			It("Resolves correctly", func() {
				Expect(g.Complete(&a)).To(Succeed())
				Expect(g.Complete(&b)).To(Succeed())
				Expect(a.B).To(Equal(&b))
				Expect(b.A).To(Equal(&a))
			})
		})

		Context("Graph with multiple layers", func() {
			type LevelOne struct {
			}

			type LevelTwo struct {
				LOne *LevelOne `inject:""`
			}

			type LevelThree struct {
				LTwo *LevelTwo `inject:""`
			}

			var (
				one   LevelOne
				two   LevelTwo
				three LevelThree
			)

			BeforeEach(func() {
				g.Provide(
					&ValueProvider{
						Value: &one,
					},
					&ValueProvider{
						Value: &two,
					},
					&ValueProvider{
						Value: &three,
					},
				)
			})

			It("Resolves all correctly", func() {
				Expect(g.Complete(&three)).To(Succeed())
				Expect(three.LTwo).To(Equal(&two))
				Expect(two.LOne).To(Equal(&one))
			})
		})

		Context("Graph with named injections", func() {
			type Struct struct {
				X int `inject:"ValA"`
				Y int `inject:"ValB"`
				Z int `inject:""`
			}

			var (
				v                Struct
				valA, valB, valC int
			)

			BeforeEach(func() {
				g.Provide(
					&ValueProvider{
						Name:  "ValA",
						Value: valA,
					},
					&ValueProvider{
						Name:  "ValB",
						Value: valB,
					},
					&ValueProvider{
						Value: valC,
					},
				)
			})

			It("Resolves correctly", func() {
				Expect(g.Complete(&v)).To(Succeed())
			})
		})

		Context("Graph with context selected injections", func() {
			Context("Struct type context", func() {
				var (
					a, b Decorated
					d    Decorator
				)

				BeforeEach(func() {
					d = Decorator{}
					g.Provide(
						&ValueProvider{
							Value: &a,
						},
						&ValueProvider{
							Context: reflect.TypeOf((*Decorator)(nil)).Elem(),
							Value:   &b,
						},
					)
				})

				Context("Simple discovery of decorated service", func() {
					It("Resolves correctly", func() {
						Expect(g.Complete(&d)).To(Succeed())
						Expect(d.Decorated).To(Equal(&b))
					})
				})

				Context("External service discovering decorated service", func() {
					var service ServiceValueImpl

					BeforeEach(func() {
						g.Provide(
							&ValueProvider{
								Context: reflect.TypeOf((*ServiceValueImpl)(nil)).Elem(),
								Value:   &d,
							},
						)
					})

					It("Resolves correctly", func() {
						Expect(g.Complete(&service)).To(Succeed())
						Expect(service.X).To(Equal(&d))
						Expect(service.X.(*Decorator).Decorated).To(Equal(&a))
					})
				})
			})

			Context("Interface type context", func() {
				var b CustomA

				BeforeEach(func() {
					b = CustomA{Val: 10}

					g.Provide(
						&ValueProvider{
							Context: reflect.TypeOf((*ServiceInterface)(nil)).Elem(),
							Value:   &b,
						},
					)
				})

				Context("Value receiver", func() {
					var s ServiceValueImpl

					BeforeEach(func() {
						s = ServiceValueImpl{}
					})

					It("Resolves correctly", func() {
						Expect(g.Complete(&s)).To(Succeed())
						Expect(s.X).To(Equal(&b))
					})
				})

				Context("Pointer receiver", func() {
					var s ServicePtrImpl

					BeforeEach(func() {
						s = ServicePtrImpl{}
					})

					It("Resolves correctly", func() {
						Expect(g.Complete(&s)).To(Succeed())
						Expect(s.X).To(Equal(&b))
					})
				})
			})
		})

		Context("Graph with value-receiver interfaces", func() {
			var (
				a InterfaceA
				b InterfaceB
			)

			BeforeEach(func() {
				a = ImplA{}
				b = ImplB{}

				g.Provide(
					&ValueProvider{
						Value: a,
					},
					&ValueProvider{
						Value: b,
					},
				)
			})

			It("Does not resolve", func() {
				Expect(g.Complete(&a)).ToNot(Succeed())
				Expect(g.Complete(&b)).ToNot(Succeed())
			})
		})

		Context("Graph with pointer-receiver interfaces", func() {
			var (
				a InterfaceA
				b InterfaceB
			)

			BeforeEach(func() {
				a = &PtrImplA{}
				b = &PtrImplB{}

				g.Provide(
					&ValueProvider{
						Value: a,
					},
					&ValueProvider{
						Value: b,
					},
				)
			})

			It("Resolves correctly", func() {
				Expect(g.Complete(&a)).To(Succeed())
				Expect(g.Complete(&b)).To(Succeed())
				Expect(a.(*PtrImplA).B).To(Equal(b))
				Expect(b.(*PtrImplB).A).To(Equal(a))
			})
		})

		Context("Graph with builder", func() {
			var decorated PtrDecorated

			Context("Without dependencies", func() {
				BeforeEach(func() {
					g.Provide(
						&BuilderProvider{
							Builder: func() InterfaceA {
								return &decorated
							},
							ResolveContext: g,
						},
					)
				})

				It("Resolves correctly", func() {
					var decorator Decorator
					Expect(g.Complete(&decorator)).To(Succeed())
					Expect(decorator.Decorated).To(Equal(&decorated))
				})
			})

			Context("With dependencies", func() {
				Context("Primitive type", func() {
					BeforeEach(func() {
						g.Provide(
							&ValueProvider{
								Complete: true,
								Value:    5,
							},
							&BuilderProvider{
								Builder: func(v int) InterfaceA {
									return &decorated
								},
								ResolveContext: g,
							},
						)
					})

					It("Resolves correctly", func() {
						var decorator Decorator
						Expect(g.Complete(&decorator)).To(Succeed())
						Expect(decorator.Decorated).To(Equal(&decorated))
					})
				})

				Context("Pointer type", func() {
					BeforeEach(func() {
						g.Provide(
							&ValueProvider{
								Value: &decorated,
							},
							&BuilderProvider{
								Builder: func(v *Decorated) *ServiceValueImpl {
									return &ServiceValueImpl{
										X: v,
									}
								},
								ResolveContext: g,
							},
						)
					})

					It("Resolves correctly", func() {
						var decorator Decorator
						Expect(g.Complete(&decorator)).To(Succeed())
						Expect(decorator.Decorated).To(Equal(&decorated))
					})
				})

				Context("Interface type", func() {
					BeforeEach(func() {
						g.Provide(
							&ValueProvider{
								Value: &decorated,
							},
							&BuilderProvider{
								Builder: func(v InterfaceA) *ServiceValueImpl {
									return &ServiceValueImpl{
										X: v,
									}
								},
								ResolveContext: g,
							},
						)
					})

					It("Resolves correctly", func() {
						type ServiceB struct {
							Y *ServiceValueImpl `inject:""`
						}

						var b ServiceB
						Expect(b.Y).To(BeNil())
						Expect(g.Complete(&b)).To(Succeed())
						Expect(b.Y.X).To(Equal(&decorated))
					})
				})
			})
		})
	})

	Describe("Find", func() {
		var (
			typeInfo, context reflect.Type
			name              string
			expected          error
		)

		BeforeEach(func() {
			typeInfo = reflect.TypeOf((*StructA)(nil))
		})

		Context("With no matching providers", func() {
			BeforeEach(func() {
				expected = fmt.Errorf("Could not find provider for %s.", typeInfo)
			})

			It("Returns nil and error", func() {
				v, err := g.Find(typeInfo, context, name)
				Expect(v).To(BeNil())
				Expect(err).To(Equal(expected))
			})
		})

		Context("With just the right number of providers (1)", func() {
			BeforeEach(func() {
				g.Provide(
					&ValueProvider{Value: &StructA{}},
				)
			})

			It("Returns the value and no error", func() {
				v, err := g.Find(typeInfo, context, name)
				Expect(reflect.TypeOf(v)).To(Equal(typeInfo))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("With too many matching providers", func() {
			BeforeEach(func() {
				expected = fmt.Errorf("Found multiple providers for type: %s, context: %s, name: %s.",
					typeInfo,
					context,
					name)

				g.Provide(
					&ValueProvider{Value: &StructA{}},
					&ValueProvider{Value: &StructA{}},
				)
			})

			It("Returns nil and error", func() {
				v, err := g.Find(typeInfo, context, name)
				Expect(v).To(BeNil())
				Expect(err).To(Equal(expected))
			})
		})
	})

	Describe("Resolve", func() {
		DescribeTable("Resolve with different value types",
			func(providers ...Provider) {
				g.Provide(providers...)
				Expect(g.Resolve()).To(Succeed())
			},
			Entry("An int", &ValueProvider{Value: 5}),
			Entry("A map", &ValueProvider{Value: map[string]string{}}),
		)

		Context("There exists an unmet dependency", func() {
			BeforeEach(func() {
				g.Provide(
					&ValueProvider{Value: &StructA{}},
				)
			})

			It("Returns the expected error correctly", func() {
				Expect(g.Resolve()).To(Equal(fmt.Errorf("Encountered error while trying to complete (%s): %s",
					reflect.TypeOf((*StructA)(nil)),
					fmt.Errorf("Encountered error attempting to set a field (%s) of type %s: %s",
						"B",
						reflect.TypeOf((*StructB)(nil)),
						fmt.Errorf("Could not find provider for %s.",
							reflect.TypeOf((*StructB)(nil)))))))
			})
		})

		Context("Builder function", func() {
			var (
				val int
				a   CustomA
			)

			BeforeEach(func() {
				val = 10

				a = CustomA{
					Val: val,
				}

				g.Provide(
					&ValueProvider{
						Value: &a,
					},
					&BuilderProvider{
						Builder: func(v InterfaceA) *ServiceValueImpl {
							return &ServiceValueImpl{
								X: v,
							}
						},
						ResolveContext: g,
					},
				)
			})

			It("It resolves correctly", func() {
				type ServiceB struct {
					Y *ServiceValueImpl `inject:""`
				}

				var b ServiceB

				g.Provide(
					&ValueProvider{
						Value: &b,
					},
				)

				Expect(b.Y).To(BeNil())
				Expect(g.Resolve()).To(Succeed())
				Expect(b.Y.X).To(Equal(&a))
				Expect(b.Y.X.(*CustomA).Val).To(Equal(val))
			})
		})

		Context("Graph with context selected injections", func() {
			var (
				b Decorated
				d Decorator
			)

			BeforeEach(func() {
				g.Provide(
					&ValueProvider{
						Context: reflect.TypeOf((*InterfaceA)(nil)).Elem(),
						Value:   &b,
					},
				)
			})

			It("Resolves correctly", func() {
				g.Provide(
					&ValueProvider{Value: &d},
				)
				Expect(g.Resolve()).To(Succeed())
				Expect(d.Decorated).To(Equal(&b))
			})
		})
	})
})
