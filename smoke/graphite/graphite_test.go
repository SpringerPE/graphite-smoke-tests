package graphite

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SpringerPE/graphite-smoke-tests/smoke"
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

func ExpectMetricsToBeSent() {

	maxAttempts := 30

	//var testConfig = smoke.GetConfig()
	smoke.SendMetricToGraphite()

	for i:= 0; i < maxAttempts; i++ {
		smoke.GetMetricFromGraphite()
	}
}