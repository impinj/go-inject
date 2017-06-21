package inject_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"errors"
	. "github.com/impinj/go-inject/inject"
	"github.com/impinj/go-inject/inject/mock"
	. "reflect"
)

var _ = Describe("BuilderProvider", func() {
	var (
		p BuilderProvider
	)

	type Struct struct {
		State int
	}

	DescribeTable("GetType",
		func(builderFunc interface{}, expectedError string) {
			if "" != expectedError {
				defer func() {
					if r := recover(); r != nil {
						if r != expectedError {
							panic(r)
						}
					}
				}()
			}
			p.Builder = builderFunc
			Expect(p.GetType()).To(Equal(TypeOf(builderFunc).Out(0)))
		},
		Entry("Value", 5, "reflect: Out of non-func type"),
		Entry("Func", func() int { return 5 }, nil),
	)

	Describe("IsComplete", func() {
		BeforeEach(func() {
			p = BuilderProvider{}
		})

		It("Defaults to false", func() {
			Expect(p.IsComplete()).To(BeFalse())
		})
	})

	Describe("Resolve", func() {
		var (
			mockCtrl  *gomock.Controller
			mockGraph *mock_inject.MockGraph
			calls     []*gomock.Call
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockGraph = mock_inject.NewMockGraph(mockCtrl)
			p.ResolveContext = mockGraph
			calls = []*gomock.Call{}
		})

		JustBeforeEach(func() {
			gomock.InOrder(calls...)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

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

			DescribeTable("With arguments",
				func(succeeds bool, builder interface{}, mockResponseForArgs ...[]interface{}) {
					p.Builder = builder

					typeInfo := TypeOf(builder)
					for i := 0; i < typeInfo.NumIn(); i++ {
						argTypeInfo := typeInfo.In(i)
						switch argTypeInfo.Kind() {
						case Interface:
							calls = append(calls,
								mockGraph.EXPECT().Find(argTypeInfo, nil, "").Return(mockResponseForArgs[i]...),
							)

						default:
							calls = append(calls,
								mockGraph.EXPECT().Complete(gomock.Any()).Do(func(v interface{}) {
									Expect(argTypeInfo).To(Equal(TypeOf(v).Elem()))
								}).Return(mockResponseForArgs[i][1]),
							)
						}
					}

					v := p.Resolve()
					if succeeds {
						Expect(v).ToNot(BeNil())
					} else {
						Expect(v).To(BeNil())
					}
				},
				Entry("Value args", true, func(_ string, _ int) *Struct {
					return &Struct{}
				}, []interface{}{"hello", nil}, []interface{}{5, nil}),
				Entry("Value arg with error", false, func(_ string) *Struct {
					return &Struct{}
				}, []interface{}{nil, errors.New("Encountered error")}),
				Entry("Interface args", true, func(_ interface{}) *Struct {
					return &Struct{}
				}, []interface{}{5, nil}),
				Entry("Interface arg with error", false, func(_ interface{}) *Struct {
					return &Struct{}
				}, []interface{}{nil, errors.New("Encountered error")}),
			)
		})

		DescribeTable("Given a non-func as a builder",
			func(v interface{}) {
				p.Builder = v
				Expect(p.Resolve()).To(BeNil())
			},
			Entry("Nil", nil),
			Entry("Value", 5),
		)
	})
})
