package inject_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"go-inject/inject/mock"
	"go-inject/inject"
)

var _ = Describe("SingletonProvider", func() {
	var (
		mockCtrl *gomock.Controller
		mockProvider *mock_inject.MockProvider
		provider *inject.SingletonProvider
		resolvedValue interface{}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockProvider = mock_inject.NewMockProvider(mockCtrl)
		provider = &inject.SingletonProvider{}

	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Resolve", func() {
		Context("With no wrapped provider", func() {
			BeforeEach(func() {
				provider.Provider = nil
			})

			It("Returns nil", func() {
				Expect(provider.Resolve()).To(BeNil())
			})
		})

		Context("With a wrapped provider", func() {
			BeforeEach(func() {
				provider.Provider = mockProvider
				resolvedValue = struct{}{}
				mockProvider.EXPECT().Resolve().Return(resolvedValue)
			})

			It("Returns the expected value", func() {
				Expect(provider.Resolve()).To(BeIdenticalTo(resolvedValue))
			})

			It("Does not invoke the wrapped Resolve() twice", func() {
				provider.Resolve()
				provider.Resolve()
			})
		})
	})
})