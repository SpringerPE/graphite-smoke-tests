package graphite

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("Graphite:", func() {



	BeforeEach(func() {


		})

	AfterEach(func() {


		})

	Context("send single metric value", func() {
		It("can be sent and retrieved", func() {
			Expect(sendMetric()).To(Equal(0))
			})
		})
})

func sendMetric() int{
	return 0
}